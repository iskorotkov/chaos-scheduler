package main

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/internal/handlers"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo"
	"github.com/iskorotkov/chaos-scheduler/pkg/k8s"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/execute"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

func main() {
	// Handle panics.
	defer func() {
		r := recover()
		if r != nil {
			log.Printf("panic occurred: %v", r)
			debug.PrintStack()
			os.Exit(1)
		}
	}()

	cfg, err := config.FromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	logger := createLogger(cfg)
	defer syncLogger(logger)

	logger.Infow("got config from environment", "config", cfg)

	finder, err := k8s.NewFinder(logger.Named("finder"))
	if err != nil {
		logger.Fatal(err)
	}

	executor, err := argo.NewExecutor(cfg.ArgoServer, logger.Named("argo"))
	if err != nil {
		logger.Fatal(err)
	}

	r := createRouter(cfg, finder, executor, logger)
	if err = http.ListenAndServe(":8811", r); err != nil {
		logger.Fatal(err.Error())
	}
}

// createRouter creates and configures chi router.
func createRouter(cfg *config.Config, finder targets.TargetFinder, executor execute.Executor, logger *zap.SugaredLogger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))

	r.Route("/api", func(r chi.Router) {
		r.Mount("/v1", handlers.Router(cfg, finder, executor, logger.Named("handlers")))
	})

	return r
}

// createLogger creates and configures zap logger.
func createLogger(cfg *config.Config) *zap.SugaredLogger {
	var (
		logger *zap.Logger
		err    error
	)
	if cfg.Development {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatal(err)
	}

	return logger.Sugar()
}

// syncLogger flushed logger buffer to stdout.
func syncLogger(logger *zap.SugaredLogger) {
	err := logger.Sync()
	if err != nil {
		log.Fatal(err.Error())
	}
}
