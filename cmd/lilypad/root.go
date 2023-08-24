package lilypad

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

var Fatal = FatalErrorHandler

func init() { //nolint:gochecknoinits
	NewRootCmd()
}

func NewRootCmd() *cobra.Command {
	RootCmd := &cobra.Command{
		Use:   getCommandLineExecutable(),
		Short: "Waterlily",
		Long:  `Waterlily`,
	}
	RootCmd.AddCommand(newServeCmd())
	return RootCmd
}

func Execute() {
	RootCmd := NewRootCmd()
	RootCmd.SetContext(context.Background())
	RootCmd.SetOutput(os.Stdout)
	if err := RootCmd.Execute(); err != nil {
		Fatal(RootCmd, err.Error(), 1)
	}
}
