package workflows

import (
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo"
	"github.com/iskorotkov/chaos-scheduler/pkg/k8s"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemble"
	"go.uber.org/zap"
	"net/http"
)

type createResponse struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func create(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) {
	cfg, ok := configFromContext(r)
	if !ok {
		msg := "couldn't get config from request context"
		logger.Error(msg)
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

	logger.Infow("workflow created", "name", workflow.Name, "namespace", workflow.Namespace)

	w.Header().Add("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(createResponse{Name: workflow.Name, Namespace: workflow.Namespace})
	if err != nil {
		logger.Errorw(err.Error())
		http.Error(w, "couldn't encode response as JSON", http.StatusInternalServerError)
		return
	}
}

func executeWorkflow(r *http.Request, cfg *config.Config, logger *zap.SugaredLogger) (assemble.Workflow, error) {
	form, ok := parseForm(r, logger.Named("params"))
	if !ok {
		return assemble.Workflow{}, formParseError
	}

	finder, err := k8s.NewFinder(logger.Named("targets"))
	if err != nil {
		logger.Error(err)
		return assemble.Workflow{}, internalError
	}

	executor, err := argo.NewExecutor(cfg.ArgoServer, logger.Named("argo"))
	if err != nil {
		logger.Error(err)
		return assemble.Workflow{}, internalError
	}

	sp := workflows.ScenarioParams{
		Seed:          form.Seed,
		Stages:        form.Stages,
		AppNS:         cfg.AppNS,
		AppLabel:      cfg.AppLabel,
		StageDuration: cfg.StageDuration,
		Failures:      enabledFailures(cfg),
		TargetFinder:  finder,
	}
	wp := workflows.WorkflowParams{Extensions: enabledExtensions(cfg, logger.Named("extensions"))}
	ep := workflows.ExecutionParams{Executor: executor}

	wf, err := workflows.ExecuteWorkflow(sp, wp, ep, logger.Named("workflows"))
	if err != nil {
		logger.Errorw(err.Error())

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
