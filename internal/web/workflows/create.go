package workflows

import (
	"encoding/json"
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

	wf, params, err := createWorkflowFromForm(r, cfg, logger)
	if err != nil {
		if err == formParseError || err == scenarioParamsError {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	executor := executors.NewGRPCExecutor(cfg.ServerURL, logger.Named("execution"))
	wf, err = executor.Execute(wf)
	if err != nil {
		logger.Errorw(err.Error(),
			"config", cfg)
		http.Error(w, "couldn't execute scenario", http.StatusInternalServerError)
		return
	}

	logger.Infow("workflow created",
		"name", wf.Name,
		"namespace", wf.Namespace)

	data := workflow{Workflow: wf, Params: params}

	w.Header().Add("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Errorw(err.Error(),
			"data", data)
		http.Error(w, "couldn't encode response as JSON", http.StatusInternalServerError)
		return
	}
}
