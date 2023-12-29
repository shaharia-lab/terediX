// Package cmd provides commands
package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCmd build root command
func NewRootCmd(version string) *cobra.Command {
	cmd := cobra.Command{
		Version: version,
		Short:   "Discover and Explore your tech resources",
		Long:    "Discover and Explore your tech resources",
	}

	cmd.AddCommand(
		NewDiscoverCommand(),
		NewRelationCommand(),
		NewValidateCommand(),
	)

	return &cmd
}
