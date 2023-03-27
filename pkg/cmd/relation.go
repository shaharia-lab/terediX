package cmd

import (
	"github.com/spf13/cobra"
	"teredix/pkg/config"
	"teredix/pkg/storage"
)

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
