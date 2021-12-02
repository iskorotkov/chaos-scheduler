// Package handlers handles web requests to preview or create workflows.
package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/execute"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
)

// Router returns a handler with configured routes.
func Router(cfg *config.Config, finder targets.TargetFinder, executor execute.Executor, logger *zap.SugaredLogger) http.Handler {
	r := chi.NewRouter()

	r.Get("/targets", func(w http.ResponseWriter, r *http.Request) {
		getAvailableTargets(w, cfg, finder, logger.Named("targets"))
	})
	r.Get("/failures", func(w http.ResponseWriter, r *http.Request) {
		getAvailableFailures(w, cfg, logger.Named("failures"))
	})
	r.Get("/namespaces", func(w http.ResponseWriter, r *http.Request) {
		getAvailableNamespaces(w, cfg, logger.Named("namespaces"))
	})
	r.Get("/workflows/preview", func(w http.ResponseWriter, r *http.Request) {
		preview(w, r, cfg, finder, logger.Named("preview"))
	})
	r.Post("/workflows/create", func(w http.ResponseWriter, r *http.Request) {
		create(w, r, cfg, finder, executor, logger.Named("create"))
	})

	return r
}

type Seeds = workflows.Seeds

type Stages = workflows.Stages

type WorkflowBody struct {
	Namespace string   `json:"namespace"`
	Seeds     Seeds    `json:"seeds"`
	Stages    Stages   `json:"stages"`
	Targets   []string `json:"targets"`
	Failures  []string `json:"failures"`
}

func parseWorkflowRequest(r *http.Request, log *zap.SugaredLogger) (WorkflowBody, bool) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("error reading request body: %v", err)
		return WorkflowBody{}, false
	}

	var body WorkflowBody
	if err := json.Unmarshal(b, &body); err != nil {
		log.Errorf("error marshaling request body to json: %v", err)
		return WorkflowBody{}, false
	}

	return body, true
}
