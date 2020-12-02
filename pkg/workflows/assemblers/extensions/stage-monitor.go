package extensions

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenarios"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
)

type StageMonitor struct {
	image    string
	targetNs string
}

func (s StageMonitor) Apply(_ scenarios.Stage, stageIndex int) Extension {
	if s.image == "" {
		logger.Warning("stage monitor image wasn't specified; stage monitor creation skipped")
		return nil
	}

	name := fmt.Sprintf("stage-monitor-%d", stageIndex+1)
	return templates.NewContainerTemplate(name, templates.Container{
		Name:  "stage-monitor",
		Image: s.image,
		Env: []templates.EnvVar{
			{"TARGET_NAMESPACE", s.targetNs},
			{"DURATION", "1m"},
		},
		Ports:   nil,
		Command: nil,
		Args:    nil,
	})
}

func UseStageMonitor(image string, targetNs string) StageExtension {
	return StageMonitor{image: image, targetNs: targetNs}
}
