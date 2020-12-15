package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/internal/web/files"
	"github.com/iskorotkov/chaos-scheduler/internal/web/home"
	"github.com/iskorotkov/chaos-scheduler/internal/web/scenarios"
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

	r.Use(configCtx(cfg))

	mapRoutes(r, sugar)

	err = http.ListenAndServe(":8811", r)
	if err != nil {
		sugar.Fatal(err.Error())
	}
}

func mapRoutes(r chi.Router, logger *zap.SugaredLogger) {
	r.Mount("/", home.Router(logger.Named("home")))
	r.Mount("/scenarios", scenarios.Router(logger.Named("scenarios")))

	r.Get("/js/*", files.ServeFolder("./web/js"))
	r.Get("/css/*", files.ServeFolder("./web/css"))
}

func syncLogger(logger *zap.SugaredLogger) {
	err := logger.Sync()
	if err != nil {
		log.Fatal(err.Error())
	}
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
