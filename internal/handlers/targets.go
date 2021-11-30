package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
)

func getAvailableTargets(w http.ResponseWriter, cfg *config.Config, finder targets.TargetFinder, log *zap.SugaredLogger) {
	type Target struct {
		ID        string `json:"id,omitempty"`
		Type      string `json:"type,omitempty"`
		Name      string `json:"name,omitempty"`
		Namespace string `json:"namespace,omitempty"`
		Count     int    `json:"count,omitempty"`
	}

	log.Infof("getting list of available targets")

	tt, err := finder.List(cfg.AppNS, cfg.AppLabel)
	if err != nil {
		log.Errorf("error getting list of targets: %v", err)
		http.Error(w, "error getting list of targets", http.StatusInternalServerError)
		return
	}

	// Fill map (deployment name => list of pods).
	targetsMap := make(map[string][]Target)
	for _, t := range tt {
		targetsMap[t.AppLabelValue] = append(targetsMap[t.AppLabelValue], Target{
			ID:        t.ID(),
			Type:      "deployment",
			Name:      t.AppLabelValue,
			Namespace: cfg.AppNS,
			Count:     1,
		})
	}

	// Convert map to slice.
	var targetsSlice []Target
	for _, pods := range targetsMap {
		// Take any pod in deployment and set correct count.
		pod := pods[0]
		pod.Count = len(pods)

		targetsSlice = append(targetsSlice, pod)
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(targetsSlice); err != nil {
		log.Errorf("error encoding response: %v", err)
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}
