// Package handlers handles web requests to preview or create workflows.
package handlers

import (
	"github.com/go-chi/chi"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/execute"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// form describes request form fields required to generate workflow.
type form struct {
	// Seed is a seed used in both failure and target selection.
	Seed int64 `json:"seed"`
	// Stages is a number of stages in generated workflow.
	Stages int `json:"stages"`
}

// Router returns a handler with configured routes.
func Router(cfg *config.Config, finder targets.TargetFinder, executor execute.Executor, logger *zap.SugaredLogger) http.Handler {
	r := chi.NewRouter()

	r.Get("/preview", func(writer http.ResponseWriter, request *http.Request) {
		preview(writer, request, cfg, finder, logger.Named("preview"))
	})
	r.Post("/create", func(writer http.ResponseWriter, request *http.Request) {
		create(writer, request, cfg, finder, executor, logger.Named("create"))
	})

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
