package workflows

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

var (
	formParseError = errors.New("couldn't parse form data")
	paramsError    = errors.New("params are invalid")
	internalError  = errors.New("internal error occurred")
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
		logger.Errorw(err.Error(), "stages", r.FormValue("stages"))
		return form{}, false
	}

	seed, err := strconv.ParseInt(r.FormValue("seed"), 10, 64)
	if err != nil {
		logger.Errorw(err.Error(), "seed", r.FormValue("seed"))
		return form{}, false
	}

	return form{Seed: seed, Stages: int(stages)}, true
}

func configFromContext(r *http.Request) (*config.Config, bool) {
	entry := r.Context().Value("config")
	cfg, ok := entry.(*config.Config)
	return cfg, ok
}
