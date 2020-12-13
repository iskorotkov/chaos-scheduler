package scenarios

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"net/http"
)

func Status(w http.ResponseWriter, r *http.Request) {
	cfg, ok := r.Context().Value("config").(*config.Config)
	if !ok {
		logger.Error(ConfigError)
		http.Error(w, ConfigError.Error(), http.StatusInternalServerError)
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

	handler := server.Page("web/html/scenarios/view.gohtml", params)
	handler(w, r)
}
