package workflows

import (
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
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

	workflow, err := generateWorkflow(r, cfg, logger)
	if err != nil {
		if err == formParseError || err == scenarioParamsError {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.Header().Add("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(workflow)
	if err != nil {
		logger.Errorw(err.Error(),
			"data", workflow)
		http.Error(w, "couldn't encode response as JSON", http.StatusInternalServerError)
		return
	}
}
