package workflows

import (
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemble"
	"go.uber.org/zap"
	"net/http"
)

func create(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) {
	entry := r.Context().Value("config")
	cfg, ok := entry.(*config.Config)
	if !ok {
		msg := "couldn't get config from request context"
		logger.Error(msg,
			"config", entry)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	workflow, err := executeWorkflow(r, cfg, logger)
	if err != nil {
		if err == formParseError || err == paramsError {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	logger.Infow("workflow created",
		"name", workflow.Name,
		"namespace", workflow.Namespace)

	w.Header().Add("Content-Type", "application/json")

	data := struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	}{
		Name:      workflow.Name,
		Namespace: workflow.Namespace,
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Errorw(err.Error(),
			"data", data)
		http.Error(w, "couldn't encode response as JSON", http.StatusInternalServerError)
		return
	}
}

func executeWorkflow(r *http.Request, cfg *config.Config, logger *zap.SugaredLogger) (assemble.Workflow, error) {
	form, err := parseForm(r, logger.Named("params"))
	if err != nil {
		return assemble.Workflow{}, err
	}

	sp, err := createScenarioParams(scenarioParams{
		server:        cfg.ArgoServer,
		namespace:     cfg.AppNS,
		label:         cfg.AppLabel,
		stageDuration: cfg.StageDuration,
		seed:          form.Seed,
		stages:        form.Stages,
		failures:      enabledFailures(cfg),
	}, logger.Named("scenario-params"))
	if err != nil {
		return assemble.Workflow{}, err
	}

	wp := workflows.WorkflowParams{Extensions: enabledExtensions(cfg, logger.Named("extensions"))}

	executor, err := argo.NewExecutor(cfg.ArgoServer, logger.Named("argo"))
	if err != nil {
		logger.Error(err)
		return assemble.Workflow{}, internalError
	}

	ep := workflows.ExecutionParams{Executor: executor}

	wf, err := workflows.ExecuteWorkflow(sp, wp, ep, logger.Named("workflows"))
	if err != nil {
		logger.Errorw(err.Error(),
			"scenario params", sp,
			"workflow params", wp)

		if err == workflows.NotEnoughTargetsError ||
			err == workflows.NotEnoughFailuresError ||
			err == workflows.TargetsFetchError ||
			err == workflows.AssembleError ||
			err == workflows.ExecutionError {
			return assemble.Workflow{}, internalError
		} else {
			return assemble.Workflow{}, paramsError
		}
	}

	return wf, nil
}
