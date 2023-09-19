// Package cmd provides commands
package cmd

import (
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/processor"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/scanner"
	"github.com/shaharia-lab/teredix/pkg/storage"
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

	sources := scanner.NewSourceRegistry(scanner.GetScannerRegistries())

	st := storage.BuildStorage(appConfig)
	err = st.Prepare()
	if err != nil {
		return err
	}

	processConfig := processor.Config{BatchSize: appConfig.Storage.BatchSize}
	p := processor.NewProcessor(processConfig, st, sources.BuildFromAppConfig(*appConfig))

	resourceChan := make(chan resource.Resource)
	p.Process(resourceChan)

	return nil
}
