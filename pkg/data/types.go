package data

import (
	"github.com/bacalhau-project/lilypad/pkg/data/bacalhau"
)

// used by resource providers to describe their resources
// use by job offers to describe their requirements
// when used by resource providers - these are absolute values
// when used by job offers - these are minimum requirements
type MachineSpec struct {
	// Milli-GPU
	// Whilst it's unlikely that partial GPU's make sense
	// let's not use a float and fix the precision to 1/1000
	GPU int `json:"gpu"`

	// Milli-CPU
	CPU int `json:"cpu"`

	// Megabytes
	RAM int `json:"ram"`
}

// this is what is loaded from the template file in the git repo
type Module struct {
	// the min spec that this module requires
	// e.g. does this module need a GPU?
	// the module file itself will contain this spec
	// and so the module will need to be downloaded
	// and executed for this spec to be known
	Machine MachineSpec `json:"machine"`

	// the bacalhau job spec
	Job bacalhau.Job `json:"job"`
}

// describes a workload to be run
// this pins a go-template.yaml file
// that is a bacalhau job spec
type ModuleConfig struct {

	// used for the shortcuts
	// this is in the modules package
	// where we keep a map of named modules
	// and their versions onto the
	// repo, hash and path below
	Name    string `json:"name"`
	Version string `json:"version"`

	// needs to be a http url for a git repo
	// we must be able to clone it without credentials
	Repo string `json:"repo"`
	// the git hash to pin the module
	// we will 'git checkout' this hash
	Hash string `json:"hash"`
	// once the checkout has been done
	// this is the path to the module template
	// within the repo
	Path string `json:"path"`
}

type Result struct {
	// this is the cid of the result where ID is set to empty string
	ID     string `json:"id"`
	DealID string `json:"deal_id"`
	// the CID of the actual results
	DataID           string `json:"results_id"`
	InstructionCount uint64 `json:"instruction_count"`
}

// MarketPrice means - get me the best deal
// job creators will do this by default i.e. "just buy me the cheapest"
// FixedPrice means - take it or leave it
// resource creators will do this by default i.e. "this is my price"
type PricingMode string

const (
	MarketPrice PricingMode = "MarketPrice"
	FixedPrice  PricingMode = "FixedPrice"
)

// the mediator and directory services that are trusted
// by the RP and JC - the solver will find an intersection
// of these and attach them to the deal
type TrustedParties struct {
	Mediator  []string `json:"mediator"`
	Directory []string `json:"directory"`
}

// posted to the solver by a job creator
type JobOffer struct {
	// this is the cid of the job offer where ID is set to empty string
	ID string `json:"id"`
	// this is basically a nonce so we don't have one ID pointing at multiple offers
	CreatedAt int `json:"created_at"`
	// the address of the job creator
	JobCreator string `json:"job_creator"`
	// the actual module that is being offered
	// this must hash to the ModuleID above
	Module ModuleConfig `json:"module"`
	// the spec required by the module
	// this will have been hoisted from the module itself
	Spec MachineSpec `json:"spec"`
	// the user inputs to the module
	// these values will power the go template
	Inputs map[string]string `json:"inputs"`
	// tells the solver how to match these prices
	// for JC this will normally be MarketPrice
	Mode PricingMode `json:"mode"`
	// the offered price and timeouts
	Pricing  DealPricing  `json:"pricing"`
	Timeouts DealTimeouts `json:"timeouts"`

	// which parties are trusted by the job creator
	TrustedParties TrustedParties `json:"trusted_parties"`
}

// this is what the solver keeps track of so we can know
// what the current state of the deal is
type JobOfferContainer struct {
	ID         string   `json:"id"`
	DealID     string   `json:"deal_id"`
	JobCreator string   `json:"job_creator"`
	State      uint8    `json:"state"`
	JobOffer   JobOffer `json:"job_offer"`
}

// posted to the solver by a resource provider
type ResourceOffer struct {
	// this is the cid of the resource offer where ID is set to empty string
	ID string `json:"id"`
	// this is basically a nonce so we don't have one ID pointing at multiple offers
	CreatedAt int `json:"created_at"`
	// the address of the job creator
	ResourceProvider string `json:"resource_provider"`
	// allows a resource provider to manage multiple offers
	// that are essentially the same
	Index int `json:"index"`
	// the spec being offered
	Spec MachineSpec `json:"spec"`
	// the module ID's that this resource provider can run
	// an empty list means ALL modules
	Modules []string `json:"modules"`
	// tells the solver how to match these prices
	// for RP this will normally be FixedPrice
	// we expect the default pricing to be filled in
	Mode PricingMode `json:"mode"`
	// the default pricing for this resource offer
	// i.e. this is for any module
	DefaultPricing  DealPricing  `json:"default_pricing"`
	DefaultTimeouts DealTimeouts `json:"default_timeouts"`
	// the pricing for each module
	// this allows a resource provider to charge more
	// for certain modules
	ModulePricing  map[string]DealPricing  `json:"module_pricing"`
	ModuleTimeouts map[string]DealTimeouts `json:"module_timeouts"`

	// which parties are trusted by the resource provider
	TrustedParties TrustedParties `json:"trusted_parties"`
}

// this is what the solver keeps track of so we can know
// what the current state of the deal is
type ResourceOfferContainer struct {
	ID               string        `json:"id"`
	DealID           string        `json:"deal_id"`
	ResourceProvider string        `json:"resource_provider"`
	State            uint8         `json:"state"`
	ResourceOffer    ResourceOffer `json:"job_offer"`
}

type DealMembers struct {
	JobCreator       string   `json:"job_creator"`
	ResourceProvider string   `json:"resource_provider"`
	Directory        string   `json:"directory"`
	Mediators        []string `json:"mediators"`
}

type DealTimeout struct {
	Timeout    uint64 `json:"timeout"`
	Collateral uint64 `json:"collateral"`
}

type DealTimeouts struct {
	Agree          DealTimeout `json:"agree"`
	SubmitResults  DealTimeout `json:"submit_results"`
	JudgeResults   DealTimeout `json:"judge_results"`
	MediateResults DealTimeout `json:"mediate_results"`
}

type DealPricing struct {
	InstructionPrice          uint64 `json:"instruction_price"`
	PaymentCollateral         uint64 `json:"payment_collateral"`
	ResultsCollateralMultiple uint64 `json:"results_collateral_multiple"`
	MediationFee              uint64 `json:"mediation_fee"`
}

// this is the struct that will have it's ID taken and used
// as the reference for what both parties agreed to
// the solver will publish this deal to the directory
type Deal struct {
	// this is the cid of the deal where ID is set to empty string
	ID            string        `json:"id"`
	Members       DealMembers   `json:"members"`
	Pricing       DealPricing   `json:"pricing"`
	Timeouts      DealTimeouts  `json:"timeouts"`
	JobOffer      JobOffer      `json:"job_offer"`
	ResourceOffer ResourceOffer `json:"resource_offer"`
}

type DealContainer struct {
	ID               string `json:"id"`
	JobCreator       string `json:"job_creator"`
	ResourceProvider string `json:"resource_provider"`
	JobOffer         string `json:"job_offer"`
	ResourceOffer    string `json:"resource_offer"`
	State            uint8  `json:"state"`
	Deal             Deal   `json:"deal"`
}
