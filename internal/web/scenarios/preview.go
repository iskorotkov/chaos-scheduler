package scenarios

import (
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"net/http"
)

func Preview(w http.ResponseWriter, r *http.Request) {
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

	marshaled, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		logger.Error(err)
		logger.Error(MarshalError)
		http.Error(w, MarshalError.Error(), http.StatusBadRequest)
		return
	}

	params := struct {
		GeneratedWorkflow string
		Seed              int64
		Stages            int
	}{string(marshaled), form.Seed, form.Stages}

	handler := server.Page("web/html/scenarios/preview.gohtml", params)
	handler(w, r)
}
