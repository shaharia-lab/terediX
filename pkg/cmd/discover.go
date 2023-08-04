// Package cmd provides commands
package cmd

import (
	"github.com/shahariaazam/teredix/pkg/config"
	"github.com/shahariaazam/teredix/pkg/processor"
	"github.com/shahariaazam/teredix/pkg/resource"
	"github.com/shahariaazam/teredix/pkg/source"
	"github.com/shahariaazam/teredix/pkg/storage"

	"github.com/spf13/cobra"
)

// NewDiscoverCommand build "discover" command
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
	err = st.Prepare()
	if err != nil {
		return err
	}

	processConfig := processor.Config{BatchSize: appConfig.Storage.BatchSize, WorkerPoolSize: appConfig.Discovery.WorkerPoolSize}
	p := processor.NewProcessor(processConfig, st, sources)

	resourceChan := make(chan resource.Resource)
	p.Process(resourceChan)

	return nil
}
