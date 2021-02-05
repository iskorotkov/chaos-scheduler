package workflows

import (
	"github.com/go-chi/chi"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type form struct {
	Seed   int64 `json:"seed"`
	Stages int   `json:"stages"`
}

func Router(logger *zap.SugaredLogger) http.Handler {
	r := chi.NewRouter()

	r.Get("/", server.WithLogger(preview, logger.Named("preview")))
	r.Post("/", server.WithLogger(create, logger.Named("create")))
	r.Get("/{namespace}/{name}", server.WithLogger(watchWS, logger.Named("watch")))

	return r
}

func parseForm(r *http.Request, logger *zap.SugaredLogger) (form, bool) {
	stages, err := strconv.ParseInt(r.FormValue("stages"), 10, 32)
	if err != nil {
		logger.Error(err.Error())
		return form{}, false
	}

	seed, err := strconv.ParseInt(r.FormValue("seed"), 10, 64)
	if err != nil {
		logger.Error(err.Error())
		return form{}, false
	}

	return form{Seed: seed, Stages: int(stages)}, true
}
