// Package cmd provides commands
package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	// Initialize Chi router
	r := chi.NewRouter()

	// Logger Middleware
	r.Use(middleware.Logger)

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

	// Set up your handler
	http.Handle("/metrics", promhttp.Handler())

	// Use http.Server directly to gain control over its lifecycle
	server := &http.Server{
		Addr: ":2112",
	}

	// Start server in a separate goroutine so it doesn't block
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Error("failed to start http server")
		}
	}()

	var corsMiddleware = cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Adjust this to your needs
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})
	r.Route("/api/v1", func(r chi.Router) {
		// CORS Middleware
		r.With(corsMiddleware).Get("/resources", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]interface{}{}) // Empty JSON array
		})
	})

	// Setup your API routes using r.Route, r.Get, r.Post, etc.
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Static file serving for /ui
	r.Route("/ui2", func(r chi.Router) {
		workDir, _ := os.Getwd()
		filesDir := http.Dir(workDir + "/ui")
		FileServer(r, "/ui", filesDir)
	})

	workDir, _ := os.Getwd()
	staticDirectory := http.Dir(workDir + "/ui")
	r.Handle("/ui/*", http.StripPrefix("/ui", http.FileServer(staticDirectory)))

	// Use http.Server for API routes
	apiServer := &http.Server{
		Addr:    ":8080", // Use a different port for your API
		Handler: r,
	}

	// Start API server in a separate goroutine
	go func() {
		if err := apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Error("failed to start API http server")
		}
	}()

	// Wait for context cancellation (in your case, the timeout)
	<-ctx.Done()

	// Shutdown the server gracefully with a timeout.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("failed to shutdown server gracefully")
		return err
	}

	return nil
}

// FileServer conveniently sets up a http.FileServer handler to serve static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}
