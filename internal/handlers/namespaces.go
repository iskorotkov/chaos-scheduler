package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"go.uber.org/zap"
)

func getAvailableNamespaces(w http.ResponseWriter, cfg *config.Config, log *zap.SugaredLogger) {
	type Namespace struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	type Resp []Namespace

	b, err := json.Marshal(Resp{
		Namespace{
			ID:   cfg.AppNS,
			Name: cfg.AppNS,
		},
	})
	if err != nil {
		log.Errorf("error marshaling list of namespaces: %v", err)
		http.Error(w, "error marshaling list of namespaces", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if _, err := w.Write(b); err != nil {
		log.Errorf("error writing response: %v", err)
		return
	}
}
