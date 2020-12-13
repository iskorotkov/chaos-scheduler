package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/internal/web/home"
	"github.com/iskorotkov/chaos-scheduler/internal/web/scenarios"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"net/http"
	"time"
)

func main() {
	cfg := config.ParseConfigFromEnv()

	r := newRouter(cfg)

	logger.Critical(http.ListenAndServe(":8811", r))
}

func newRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	useDefaultMiddleware(r)

	r.Use(configCtx(cfg))

	serveStaticFiles(r)

	r.Mount("/", home.Router())
	r.Mount("/scenarios", scenarios.Router())

	return r
}

func serveStaticFiles(r *chi.Mux) {
	r.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js"))))
	r.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))
}

func useDefaultMiddleware(r *chi.Mux) {
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
