package scenarios

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"go.uber.org/zap"
	"net/http"
)

func viewPage(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) {
	entry := r.Context().Value("config")
	cfg, ok := entry.(*config.Config)
	if !ok {
		msg := "couldn't get config from request context"
		logger.Errorw(msg,
			"config", entry)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	name := chi.URLParam(r, "name")
	namespace := chi.URLParam(r, "namespace")

	params := struct {
		Link      string
		Name      string
		Namespace string
	}{
		Link:      fmt.Sprintf("http://%s/workflows/%s/%s", cfg.ServerURL, namespace, name),
		Name:      name,
		Namespace: namespace,
	}

	handler := server.PageHandler("web/html/scenarios/view.gohtml", params, logger.Named("page"))
	handler(w, r)
}
