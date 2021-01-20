package extensions

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"strings"
	"time"
)

type StageMonitor struct {
	image         string
	targetNs      string
	stageInterval time.Duration
	logger        *zap.SugaredLogger
}

func (s StageMonitor) Apply(stage generator.Stage, stageIndex int) []templates.Template {
	if s.image == "" {
		s.logger.Warn("stage monitor image wasn't specified; stage monitor creation skipped")
		return nil
	}

	podsToKill := make([]string, 0)
	ignoredPods := make([]string, 0)

	for _, action := range stage.Actions {
		if action.Info.Severity == metadata.SeverityCritical {
			if action.Info.Scale == metadata.ScaleNode {
				ignoredPods = append(ignoredPods, action.Target.Node)
			} else {
				podTolerance := fmt.Sprintf("%s=%d", action.Target.AppLabel, -1)
				podsToKill = append(podsToKill, podTolerance)
			}
		}
	}

	name := fmt.Sprintf("stage-monitor-%d", stageIndex+1)
	containerTemplate := templates.NewContainerTemplate(name, templates.Container{
		Name:  "stage-monitor",
		Image: s.image,
		Env: []v1.EnvVar{
			{Name: "APP_NS", Value: s.targetNs},
			{Name: "DURATION", Value: (stage.Duration + s.stageInterval).String()},
			{Name: "CRASH_TOLERANCE", Value: strings.Join(podsToKill, ";")},
			{Name: "IGNORED_PODS", Value: strings.Join(ignoredPods, ";")},
		},
	})

	return []templates.Template{containerTemplate}
}

func UseStageMonitor(image string, targetNs string, bufferTime time.Duration, logger *zap.SugaredLogger) StageExtension {
	return StageMonitor{image: image, targetNs: targetNs, stageInterval: bufferTime, logger: logger}
}
