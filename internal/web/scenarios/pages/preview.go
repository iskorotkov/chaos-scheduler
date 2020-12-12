package pages

import (
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"net/http"
)

func ScenarioPreviewPage(rw http.ResponseWriter, cfg config.Config, form Form) {
	wf, err := generateWorkflow(form, cfg)
	if err != nil {
		if err == ScenarioParamsError {
			server.BadRequest(rw, err)
		} else {
			server.InternalError(rw, err)
		}

		return
	}

	marshaled, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		logger.Error(err)
		server.BadRequest(rw, MarshalError)
		return
	}

	server.HTMLPage(rw, "web/html/scenarios/preview.gohtml", struct {
		GeneratedWorkflow string
		Seed              int64
		Stages            int
	}{string(marshaled), form.Seed, form.Stages})
}
