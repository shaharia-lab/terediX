package cmd

import (
	"teredix/pkg/config"
	"teredix/pkg/processor"
	"teredix/pkg/source"
	"teredix/pkg/storage"

	"github.com/spf13/cobra"
)

func NewDiscoverCommand() *cobra.Command {
	var cfgFile string

	cmd := cobra.Command{
		Use:   "discover",
		Short: "Start discovering resources",
		Long:  "Start discovering resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cfgFile)
		},
	}

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "a valid yaml file is required")

	return &cmd
}

func run(cfgFile string) error {
	appConfig, err := config.Load(cfgFile)
	if err != nil {
		return err
	}

	err = config.Validate(appConfig)
	if err != nil {
		return err
	}

	sources := source.BuildSources(appConfig)
	st := storage.BuildStorage(appConfig)

	processConfig := processor.Config{BatchSize: appConfig.Storage.BatchSize}
	p := processor.NewProcessor(processConfig, st, sources)
	p.Process()

	return nil
}
