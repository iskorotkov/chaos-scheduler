package handlers

import (
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generate"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"net/http"
)

// previewResponse is a response returned after scenario was generated.
type previewResponse struct {
	Scenario generate.Scenario `json:"scenario"`
}

// preview handles requests to create and preview scenario.
func preview(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) {
	cfg, ok := r.Context().Value("config").(*config.Config)
	if !ok {
		msg := "couldn't get config from request context"
		logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	finder, ok := r.Context().Value("finder").(targets.TargetFinder)
	if !ok {
		msg := "couldn't get target finder from request context"
		logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	form, ok := parseForm(r, logger.Named("params"))
	if !ok {
		msg := "couldn't parse form data"
		logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	params := workflows.ScenarioParams{
		Seed:          form.Seed,
		Stages:        form.Stages,
		AppNS:         cfg.AppNS,
		AppLabel:      cfg.AppLabel,
		StageDuration: cfg.StageDuration,
		Failures:      enabledFailures(cfg),
		TargetFinder:  finder,
	}

	scenario, err := workflows.CreateScenario(params, logger.Named("workflows"))
	if err != nil {
		logger.Error(err)
		if err == workflows.ErrNotEnoughTargets ||
			err == workflows.ErrNotEnoughFailures ||
			err == workflows.ErrAssemble ||
			err == workflows.ErrTargetsFetch {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
