package workflows

import (
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generate"
	"go.uber.org/zap"
	"net/http"
)

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

	data := struct {
		Scenario generate.Scenario `json:"scenario"`
	}{Scenario: scenario}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Errorw(err.Error(),
			"data", data)
		http.Error(w, "couldn't encode response as JSON", http.StatusInternalServerError)
		return
	}
}

func generateScenario(r *http.Request, cfg *config.Config, logger *zap.SugaredLogger) (generate.Scenario, error) {
	workflowParams, err := parseWorkflowParams(r, logger.Named("params"))
	if err != nil {
		return generate.Scenario{}, err
	}

	params := workflows.ScenarioParams{
		Seed:          workflowParams.Seed,
		Stages:        workflowParams.Stages,
		AppNS:         cfg.AppNS,
		AppLabel:      cfg.AppLabel,
		StageDuration: cfg.StageDuration,
		Failures:      enabledFailures(cfg),
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
