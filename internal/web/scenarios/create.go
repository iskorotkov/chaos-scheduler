package scenarios

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/executors"
	"net/http"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	form, err := parseScenarioParams(r)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cfg, ok := r.Context().Value("config").(*config.Config)
	if !ok {
		logger.Error(ConfigError)
		http.Error(w, ConfigError.Error(), http.StatusInternalServerError)
		return
	}

	wf, err := generateWorkflow(form, cfg)
	if err != nil {
		if err == ScenarioParamsError {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	executor := executors.NewGRPCExecutor(cfg.ServerURL)
	wf, err = executor.Execute(wf)
	if err != nil {
		logger.Error(err)
		logger.Error(ScenarioExecutionError)
		http.Error(w, ScenarioExecutionError.Error(), http.StatusInternalServerError)
		return
	}

	path := fmt.Sprintf("/scenarios/view/%s/%s", wf.Namespace, wf.Name)
	http.Redirect(w, r, path, http.StatusSeeOther)
}
