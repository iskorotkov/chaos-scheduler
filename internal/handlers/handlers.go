// Package handlers handles web requests to preview or create workflows.
package handlers

import (
	"github.com/go-chi/chi"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
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
func Router(logger *zap.SugaredLogger) http.Handler {
	r := chi.NewRouter()

	r.Get("/", server.WithLogger(preview, logger.Named("preview")))
	r.Post("/", server.WithLogger(create, logger.Named("create")))

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
