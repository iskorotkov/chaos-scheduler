package scenarios

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/executors"
	"go.uber.org/zap"
	"net/http"
)

func createAction(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) {
	entry := r.Context().Value("config")
	cfg, ok := entry.(*config.Config)
	if !ok {
		msg := "couldn't get config from request context"
		logger.Error(msg,
			"config", entry)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	wf, _, err := createWorkflowFromForm(r, cfg, logger)
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

	path := fmt.Sprintf("/scenarios/view/%s/%s", wf.Namespace, wf.Name)
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func createPage(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) {
	server.PageHandler("web/html/scenarios/create.gohtml", nil, logger.Named("page"))(w, r)
}
