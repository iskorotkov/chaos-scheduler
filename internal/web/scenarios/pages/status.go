package pages

import (
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/executors"
	"net/http"
)

func SubmissionStatusPage(rw http.ResponseWriter, r *http.Request, cfg config.Config) {
	form, err := ReadForm(r)
	if err != nil {
		server.BadRequest(rw, err)
		return
	}

	wf, err := generateWorkflow(form, cfg)
	if err != nil {
		if err == ScenarioParamsError {
			server.BadRequest(rw, err)
		} else {
			server.InternalError(rw, err)
		}

		return
	}

	executor := executors.NewGRPCExecutor(cfg.ServerURL)
	wf, err = executor.Execute(wf)
	if err != nil {
		logger.Error(err)
		server.InternalError(rw, ScenarioExecutionError)
		return
	}

	server.HTMLPage(rw, "web/html/scenarios/submission-status.gohtml", nil)
}
