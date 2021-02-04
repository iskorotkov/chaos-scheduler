package workflows

import (
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/k8s"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generate"
	"go.uber.org/zap"
	"net/http"
)

type previewResponse struct {
	Scenario generate.Scenario `json:"scenario"`
}

func preview(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) {
	entry := r.Context().Value("config")
	cfg, ok := entry.(*config.Config)
	if !ok {
		msg := "couldn't get config from request context"
		logger.Error(msg,
			"config", entry)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	scenario, err := generateScenario(r, cfg, logger.Named("scenario-generation"))
	if err != nil {
		if err == formParseError || err == paramsError {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.Header().Add("Content-Type", "application/json")

	data := previewResponse{Scenario: scenario}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Errorw(err.Error(),
			"data", data)
		http.Error(w, "couldn't encode response as JSON", http.StatusInternalServerError)
		return
	}
}

func generateScenario(r *http.Request, cfg *config.Config, logger *zap.SugaredLogger) (generate.Scenario, error) {
	form, ok := parseForm(r, logger.Named("params"))
	if !ok {
		return generate.Scenario{}, formParseError
	}

	finder, err := k8s.NewFinder(logger.Named("targets"))
	if err != nil {
		logger.Error(err)
		return generate.Scenario{}, internalError
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
		logger.Errorw(err.Error(),
			"params", params)

		if err == workflows.NotEnoughTargetsError ||
			err == workflows.NotEnoughFailuresError ||
			err == workflows.AssembleError ||
			err == workflows.TargetsFetchError {
			return generate.Scenario{}, internalError
		} else {
			return generate.Scenario{}, paramsError
		}
	}

	return scenario, nil
}
