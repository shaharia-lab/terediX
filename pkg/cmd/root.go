package cmd

import (
	"github.com/spf13/cobra"
)

//write code in golang to create a cobra command. There will be a root command and other command will be build separately and added to the root command.

func NewRootCmd(version string) *cobra.Command {
	cmd := cobra.Command{
		Use:     "redisc",
		Version: version,
		Short:   "r",
	}

	cmd.AddCommand(
		NewDiscoverCommand(),
		NewRelationCommand(),
	)

	return &cmd
}
