// Package cmd provides commands
package cmd

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shaharia-lab/teredix/pkg/api"
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

// Server represent server
type Server struct {
	apiServer         *http.Server
	promMetricsServer *http.Server
	logger            *logrus.Logger
	storage           storage.Storage
}

// NewServer instantiate new server
func NewServer(logger *logrus.Logger, storage storage.Storage) *Server {
	return &Server{
		logger:  logger,
		storage: storage,
	}
}

func (s *Server) setupAPIServer(port string) {
	r := chi.NewRouter()

	// Set up CORS
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum age to cache preflight request
	})
	r.Use(crs.Handler)

	// Redirect the home page to /app/index.html
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/index.html", http.StatusMovedPermanently)
	})

	// Create a new router group
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/resources", api.GetAllResources(s.storage))
		})
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	s.apiServer = &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
}

func (s *Server) setupPromMetricsServer() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	s.promMetricsServer = &http.Server{
		Addr:    ":2112",
		Handler: mux,
	}
}

func (s *Server) startServer(server *http.Server, serverName string) {
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.WithError(err).Errorf("failed to start %s server", serverName)
		}
	}()
}

func (s *Server) shutdownServer(ctx context.Context, server *http.Server, serverName string) error {
	if err := server.Shutdown(ctx); err != nil {
		s.logger.WithError(err).Errorf("failed to shutdown %s gracefully", serverName)
		return err
	}
	s.apiServer = nil
	return nil
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

	s := NewServer(logger, st)
	s.setupAPIServer("8080")
	s.setupPromMetricsServer()

	s.startServer(s.promMetricsServer, "metrics")
	s.startServer(s.apiServer, "api")

	// Wait for context cancellation (in your case, the timeout)
	<-ctx.Done()

	// Shutdown the servers gracefully with a timeout.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.shutdownServer(shutdownCtx, s.apiServer, "API server"); err != nil {
		return err
	}

	if err := s.shutdownServer(shutdownCtx, s.promMetricsServer, "Prometheus metrics server"); err != nil {
		return err
	}

	return nil
}
