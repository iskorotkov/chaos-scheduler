package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/execute"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
)

// createResponse is a response returned after workflow was created and launched.
type createResponse struct {
	// Name is a name of the generated workflow.
	Name string `json:"name"`
	// Namespace is a namespace where the workflow was launched.
	Namespace string `json:"namespace"`
}

// create handles requests to create and launch workflow.
func create(w http.ResponseWriter, r *http.Request, cfg *config.Config, finder targets.TargetFinder, executor execute.Executor, logger *zap.SugaredLogger) {
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

	sp := workflows.ScenarioParams{
		Seeds:         body.Seeds,
		Stages:        body.Stages,
		AppNS:         body.Namespace,
		AppLabel:      cfg.AppLabel,
		StageDuration: cfg.StageDuration,
		Failures:      mergeFailures(fs, body.Failures),
		Targets:       mergeTargets(ts, body.Targets),
	}
	wp := workflows.WorkflowParams{Extensions: enabledExtensions(cfg, logger.Named("extensions"))}
	ep := workflows.ExecutionParams{Executor: executor}

	workflow, err := workflows.ExecuteWorkflow(sp, wp, ep, logger.Named("workflows"))
	if err != nil {
		logger.Errorw(err.Error())
		if err == workflows.ErrNotEnoughTargets ||
			err == workflows.ErrNotEnoughFailures ||
			err == workflows.ErrScenarioParams {
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
