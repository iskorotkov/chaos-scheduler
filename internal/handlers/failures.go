package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"go.uber.org/zap"
)

func getAvailableFailures(w http.ResponseWriter, cfg *config.Config, log *zap.SugaredLogger) {
	type Failure struct {
		ID       string `json:"id,omitempty"`
		Type     string `json:"type,omitempty"`
		Name     string `json:"name,omitempty"`
		Severity string `json:"severity,omitempty"`
		Scale    string `json:"scale,omitempty"`
	}

	var failures []Failure
	for _, f := range enabledFailures(cfg) {
		failures = append(failures, Failure{
			ID:       fmt.Sprintf("%s/%s/%s/%s", f.Blueprint.Type(), f.Blueprint.Name(), f.Scale, f.Severity),
			Type:     string(f.Blueprint.Type()),
			Name:     f.Blueprint.Name(),
			Severity: string(f.Severity),
			Scale:    string(f.Scale),
		})
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(failures); err != nil {
		log.Errorf("error encoding response: %v", err)
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}
