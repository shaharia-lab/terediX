// Package cmd provides commands
package cmd

import (
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/processor"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/scanner"
	"github.com/shaharia-lab/teredix/pkg/scheduler"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"
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
			logger := logrus.New()
			return run(cfgFile, logger)
		},
	}

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "a valid yaml file is required")

	return &cmd
}

func run(cfgFile string, logger *logrus.Logger) error {
	appConfig, err := config.Load(cfgFile)
	if err != nil {
		logger.WithError(err).Error("failed to load and parse configuration file")
		return err
	}

	st, err := storage.BuildStorage(appConfig)
	if err != nil {
		logger.WithError(err).Error("failed to build storage")
		return err
	}

	err = st.Prepare()
	if err != nil {
		logger.WithError(err).Error("failed to prepare storage system")
		return err
	}

	sch := scheduler.NewGoCron()
	scDeps := scanner.NewScannerDependencies(sch, st, logger)

	resourceChan := make(chan resource.Resource)
	p := processor.NewProcessor(processor.Config{BatchSize: appConfig.Storage.BatchSize}, st, scanner.BuildScanners(appConfig, scDeps), logger)
	err = p.Process(resourceChan, sch)
	if err != nil {
		logger.WithError(err).Error("failed to start processing scheduler jobs")
		return err
	}

	logger.Info("started processing scheduled jobs")
	select {}
}
