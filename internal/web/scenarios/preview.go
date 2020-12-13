package scenarios

import (
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"go.uber.org/zap"
	"net/http"
)

func previewPage(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) {
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
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	marshaled, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		logger.Errorw(err.Error(),
			"workflow", wf)
		http.Error(w, "couldn't marshall workflow to readable format", http.StatusBadRequest)
		return
	}

	data := struct {
		GeneratedWorkflow string
		Seed              int64
		Stages            int
	}{string(marshaled), params.Seed, params.Stages}

	handler := server.PageHandler("web/html/scenarios/preview.gohtml", data, logger.Named("page"))
	handler(w, r)
}
