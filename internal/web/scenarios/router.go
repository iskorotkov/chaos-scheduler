package scenarios

import (
	"github.com/go-chi/chi"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"go.uber.org/zap"
	"net/http"
)

func Router(logger *zap.SugaredLogger) http.Handler {
	r := chi.NewRouter()

	r.Get("/preview", server.WithLogger(previewPage, logger.Named("preview")))

	r.Get("/create", server.WithLogger(createPage, logger.Named("create")))
	r.Post("/create", server.WithLogger(createAction, logger.Named("create")))

	r.Get("/view/{namespace}/{name}", server.WithLogger(statusPage, logger.Named("view")))

	return r
}
