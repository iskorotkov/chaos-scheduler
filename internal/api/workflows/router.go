package workflows

import (
	"github.com/go-chi/chi"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"go.uber.org/zap"
	"net/http"
)

func Router(logger *zap.SugaredLogger) http.Handler {
	r := chi.NewRouter()

	r.Get("/", server.WithLogger(preview, logger.Named("preview")))
	r.Post("/", server.WithLogger(create, logger.Named("create")))
	r.Get("/{namespace}/{name}", server.WithLogger(watchWS, logger.Named("watch")))

	return r
}
