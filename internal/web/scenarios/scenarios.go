package scenarios

import (
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/internal/web/scenarios/pages"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"net/http"
)

func Handle(rw http.ResponseWriter, r *http.Request, cfg config.Config) {
	if r.Method == "GET" {
		err := r.ParseForm()
		if err != nil {
			logger.Error(err)
			server.BadRequest(rw, pages.FormParseError)
		}

		if len(r.Form) == 0 {
			pages.ScenarioCreationPage(rw)
		} else {
			form, err := pages.ReadForm(r)
			if err != nil {
				server.BadRequest(rw, err)
			}

			pages.ScenarioPreviewPage(rw, cfg, form)
		}
	} else if r.Method == "POST" {
		pages.SubmissionStatusPage(rw, r, cfg)
	} else {
		server.MethodNotAvailable(rw, r)
	}
}
