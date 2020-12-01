package extensions

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenarios"
)

type StageMonitor struct {
	Image string
}

func (s StageMonitor) Apply(_ scenarios.Stage, stageIndex int) Extension {
	if s.Image == "" {
		logger.Warning("stage monitor image wasn't specified; stage monitor creation skipped")
		return nil
	}

	name := fmt.Sprintf("stage-monitor-%d", stageIndex+1)
	return templates.NewContainerTemplate(name, templates.Container{
		Name:  "stage-monitor",
		Image: s.Image,
		Env: []templates.EnvVar{
			{"TARGET_NAMESPACE", "chaos-app"},
			{"DURATION", "1m"},
		},
		Ports:   nil,
		Command: nil,
		Args:    nil,
	})
}

func UseStageMonitor(image string) StageExtension {
	return StageMonitor{image}
}
