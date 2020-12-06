package extensions

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generators"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/presets"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"strings"
)

type StageMonitor struct {
	image    string
	targetNs string
}

func (s StageMonitor) Apply(stage generators.Stage, stageIndex int) Extension {
	if s.image == "" {
		logger.Warning("stage monitor image wasn't specified; stage monitor creation skipped")
		return nil
	}

	podsToKill := make([]string, 0)
	for _, action := range stage.Actions {
		if action.Info.Lethal {
			podTolerance := fmt.Sprintf("%s=%d", action.Target.Selector(), -1)
			podsToKill = append(podsToKill, podTolerance)
		}
	}

	name := fmt.Sprintf("stage-monitor-%d", stageIndex+1)
	crashTolerance := strings.Join(podsToKill, ";")

	return templates.NewContainerTemplate(name, templates.Container{
		Name:  "stage-monitor",
		Image: s.image,
		Env: []presets.EnvVar{
			{"APP_NS", s.targetNs},
			{"DURATION", "1m"},
			{"CRASH_TOLERANCE", crashTolerance},
		},
		Ports:   nil,
		Command: nil,
		Args:    nil,
	})
}

func UseStageMonitor(image string, targetNs string) StageExtension {
	return StageMonitor{image: image, targetNs: targetNs}
}
