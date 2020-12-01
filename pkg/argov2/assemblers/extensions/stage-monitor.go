package extensions

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenarios"
	"os"
)

type StageMonitor struct{}

func (s StageMonitor) Apply(_ scenarios.Stage, stageIndex int) Extension {
	name := fmt.Sprintf("stage-monitor-%d", stageIndex+1)
	image := os.Getenv("STAGE_MONITOR_IMAGE")
	if image == "" {
		logger.Warning("stage monitor image wasn't specified; no stage monitor will be created")
		return nil
	}

	return templates.NewContainerTemplate(name, templates.Container{
		Name:  "stage-monitor",
		Image: image,
		Env: []templates.EnvVar{
			{"TARGET_NAMESPACE", "chaos-app"},
			{"DURATION", "1m"},
		},
		Ports:   nil,
		Command: nil,
		Args:    nil,
	})
}
