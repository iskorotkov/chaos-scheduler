package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/iskorotkov/chaos-scheduler/internal/api/workflows"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.FromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	var logger *zap.Logger
	if cfg.Development {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatal(err)
	}

	sugar := logger.Sugar()
	defer syncLogger(sugar)

	sugar.Infow("get config from environment",
		"config", cfg)

	r := chi.NewRouter()
	useDefaultMiddleware(r)
	useCors(r)
	useConfig(r, cfg)

	mapRoutes(r, sugar)

	err = http.ListenAndServe(":8811", r)
	if err != nil {
		sugar.Fatal(err.Error())
	}
}

func useConfig(r *chi.Mux, cfg *config.Config) {
	r.Use(configCtx(cfg))
}

func mapRoutes(r chi.Router, logger *zap.SugaredLogger) {
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Mount("/workflows", workflows.Router(logger.Named("workflows")))
		})
	})
}

func syncLogger(logger *zap.SugaredLogger) {
	err := logger.Sync()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func useCors(r chi.Router) {
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))
}

func useDefaultMiddleware(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
}

func configCtx(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "config", cfg))
			next.ServeHTTP(w, r)
		})
	}
}
