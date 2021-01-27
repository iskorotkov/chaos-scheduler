package workflows

import (
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/execute"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
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

	workflow, err := generateWorkflow(r, cfg, logger)
	if err != nil {
		if err == formParseError || err == paramsError {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	workflow, err = execute.Execute(cfg.ArgoServer, workflow, logger.Named("execution"))
	if err != nil {
		logger.Errorw(err.Error(),
			"config", cfg)
		http.Error(w, "couldn't execute scenario", http.StatusInternalServerError)
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

func generateWorkflow(r *http.Request, cfg *config.Config, logger *zap.SugaredLogger) (templates.Workflow, error) {
	workflowParams, err := parseWorkflowParams(r, logger.Named("params"))
	if err != nil {
		return templates.Workflow{}, err
	}

	sp := workflows.ScenarioParams{
		Seed:          workflowParams.Seed,
		Stages:        workflowParams.Stages,
		AppNS:         cfg.AppNS,
		AppLabel:      cfg.AppLabel,
		StageDuration: cfg.StageDuration,
		Failures:      enabledFailures(cfg),
	}

	wp := workflows.WorkflowParams{
		Extensions: enabledExtensions(cfg, logger),
	}

	wf, err := workflows.CreateWorkflow(sp, wp, logger.Named("workflows"))
	if err != nil {
		logger.Errorw(err.Error(),
			"scenario params", sp,
			"workflow params", wp)

		if err == workflows.NotEnoughTargetsError ||
			err == workflows.NotEnoughFailuresError ||
			err == workflows.AssembleError ||
			err == workflows.TargetsFetchError {
			return templates.Workflow{}, internalError
		} else {
			return templates.Workflow{}, paramsError
		}
	}

	return wf, nil
}
