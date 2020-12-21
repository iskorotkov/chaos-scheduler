package workflows

import (
	"encoding/json"
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/executors"
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
		if err == formParseError || err == scenarioParamsError {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	executor := executors.NewGRPCExecutor(cfg.ServerURL, logger.Named("execution"))
	workflow.Workflow, err = executor.Execute(workflow.Workflow)
	if err != nil {
		logger.Errorw(err.Error(),
			"config", cfg)
		http.Error(w, "couldn't execute scenario", http.StatusInternalServerError)
		return
	}

	logger.Infow("workflow created",
		"name", workflow.Workflow.Name,
		"namespace", workflow.Workflow.Namespace)

	w.Header().Add("Content-Type", "application/json")

	url := fmt.Sprintf("%s/%s/%s", r.URL.Path, workflow.Workflow.Namespace, workflow.Workflow.Name)
	data := struct {
		URL string `json:"url"`
	}{url}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Errorw(err.Error(),
			"data", data)
		http.Error(w, "couldn't encode response as JSON", http.StatusInternalServerError)
		return
	}
}
