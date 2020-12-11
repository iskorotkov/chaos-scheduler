package extensions

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generators"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	v1 "k8s.io/api/core/v1"
	"strings"
	"time"
)

type StageMonitor struct {
	image         string
	targetNs      string
	stageInterval time.Duration
}

func (s StageMonitor) Apply(stage generators.Stage, stageIndex int) []templates.Template {
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

	containerTemplate := templates.NewContainerTemplate(name, templates.Container{
		Name:  "stage-monitor",
		Image: s.image,
		Env: []v1.EnvVar{
			{Name: "APP_NS", Value: s.targetNs},
			{Name: "DURATION", Value: (stage.Duration + s.stageInterval).String()},
			{Name: "CRASH_TOLERANCE", Value: crashTolerance},
		},
	})

	return []templates.Template{containerTemplate}
}

func UseStageMonitor(image string, targetNs string, bufferTime time.Duration) StageExtension {
	return StageMonitor{image: image, targetNs: targetNs, stageInterval: bufferTime}
}
