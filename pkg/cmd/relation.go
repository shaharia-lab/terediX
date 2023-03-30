// Package cmd provides commands
package cmd

import (
	"teredix/pkg/config"
	"teredix/pkg/storage"

	"github.com/spf13/cobra"
)

// NewRelationCommand build "relation" command
func NewRelationCommand() *cobra.Command {
	var cfgFile string

	cmd := cobra.Command{
		Use:   "relation",
		Short: "Build resource relationships",
		Long:  "Build resource relationships",
		RunE: func(cmd *cobra.Command, args []string) error {
			appConfig, err := config.Load(cfgFile)
			if err != nil {
				return err
			}

			err = config.Validate(appConfig)
			if err != nil {
				return err
			}

			st := storage.BuildStorage(appConfig)
			return st.StoreRelations(appConfig.Relation)
		},
	}

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "a valid yaml file is required")

	return &cmd
}
