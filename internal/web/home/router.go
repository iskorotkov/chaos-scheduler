package home

import (
	"github.com/go-chi/chi"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"go.uber.org/zap"
	"net/http"
)

func Router(l *zap.SugaredLogger) http.Handler {
	r := chi.NewRouter()
	r.Get("/", server.WithLogger(viewPage, l.Named("get")))
	return r
}
