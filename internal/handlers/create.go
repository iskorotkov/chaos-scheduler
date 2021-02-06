package handlers

import (
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/execute"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"net/http"
)

// createResponse is a response returned after workflow was created and launched.
type createResponse struct {
	// Name is a name of the generated workflow.
	Name string `json:"name"`
	// Namespace is a namespace where the workflow was launched.
	Namespace string `json:"namespace"`
}

// create handles requests to create and launch workflow.
func create(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) {
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

	executor, ok := r.Context().Value("executor").(execute.Executor)
	if !ok {
		msg := "couldn't get workflow executor from request context"
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

	workflow, err := workflows.ExecuteWorkflow(sp, wp, ep, logger.Named("workflows"))
	if err != nil {
		logger.Errorw(err.Error())
		if err == workflows.ErrNotEnoughTargets ||
			err == workflows.ErrNotEnoughFailures ||
			err == workflows.ErrTargetsFetch ||
			err == workflows.ErrAssemble ||
			err == workflows.ErrExecution {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
