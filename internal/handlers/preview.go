package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generate"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
)

// previewResponse is a response returned after scenario was generated.
type previewResponse struct {
	Scenario generate.Scenario `json:"scenario"`
}

// preview handles requests to create and preview scenario.
func preview(w http.ResponseWriter, r *http.Request, cfg *config.Config, finder targets.TargetFinder, logger *zap.SugaredLogger) {
	body, ok := parseWorkflowRequest(r, logger.Named("body"))
	if !ok {
		msg := "couldn't parse body"
		logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if body.Namespace != cfg.AppNS {
		msg := "only default namespace is allowed in current version"
		logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	ts, err := finder.List(body.Namespace, cfg.AppLabel)
	if err != nil {
		logger.Errorf("error getting list of targets: %v", err)
		http.Error(w, "error getting list of targets", http.StatusInternalServerError)
		return
	}

	fs := enabledFailures(cfg)

	params := workflows.ScenarioParams{
		Seeds:         body.Seeds,
		Stages:        body.Stages,
		AppNS:         body.Namespace,
		AppLabel:      cfg.AppLabel,
		StageDuration: cfg.StageDuration,
		Failures:      mergeFailures(fs, body.Failures),
		Targets:       mergeTargets(ts, body.Targets),
	}

	scenario, err := workflows.CreateScenario(params, logger.Named("workflows"))
	if err != nil {
		logger.Error(err)
		if err == workflows.ErrNotEnoughTargets ||
			err == workflows.ErrNotEnoughFailures ||
			err == workflows.ErrScenarioParams {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	w.Header().Add("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(previewResponse{Scenario: scenario})
	if err != nil {
		logger.Errorw(err.Error())
		http.Error(w, "couldn't encode response as JSON", http.StatusInternalServerError)
		return
	}
}
