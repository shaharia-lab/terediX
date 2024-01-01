// Package cmd provides commands
package cmd

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/metrics"
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
			appConfig, err := config.Load(cfgFile)
			if err != nil {
				logger.WithError(err).Error("failed to load and parse configuration file")
				return err
			}

			ctx := context.Background()

			return run(ctx, appConfig, logger)
		},
	}

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "a valid yaml file is required")

	return &cmd
}

func run(ctx context.Context, appConfig *config.AppConfig, logger *logrus.Logger) error {
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
	metricsCollector := metrics.NewCollector()
	scDeps := scanner.NewScannerDependencies(sch, st, logger, metricsCollector)

	resourceChan := make(chan resource.Resource)
	p := processor.NewProcessor(processor.Config{BatchSize: appConfig.Storage.BatchSize}, st, scanner.BuildScanners(appConfig.Sources, scDeps), logger, metricsCollector)
	err = p.Process(resourceChan, sch)
	if err != nil {
		logger.WithError(err).Error("failed to start processing scheduler jobs")
		return err
	}

	logger.Info("started processing scheduled jobs")

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Set up your handler
	mux.Handle("/metrics", promhttp.Handler())

	// Use http.Server directly to gain control over its lifecycle
	promMetricsServer := &http.Server{
		Addr:    ":2112",
		Handler: mux, // Use the new ServeMux as the handler
	}

	// Start server in a separate goroutine so it doesn't block
	go func() {
		if err := promMetricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Error("failed to start http server")
		}
	}()

	// Create a new router
	r := chi.NewRouter()

	// Define your routes
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Start another HTTP server for the API server
	apiServer := &http.Server{
		Addr:    ":8080",
		Handler: r, // Use the chi router as the handler
	}

	go func() {
		if err := apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Error("failed to start API server")
		}
	}()

	// Wait for context cancellation (in your case, the timeout)
	<-ctx.Done()

	// Shutdown the servers gracefully with a timeout.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := promMetricsServer.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("failed to shutdown Prometheus metrics server gracefully")
		return err
	}

	if err := apiServer.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("failed to shutdown API server gracefully")
		return err
	}

	return nil
}
