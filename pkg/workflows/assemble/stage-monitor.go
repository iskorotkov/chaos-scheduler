package assemble

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

type stageMonitor struct {
	image         string
	targetNs      string
	stageInterval time.Duration
	logger        *zap.SugaredLogger
}

func (s stageMonitor) Apply(stage generator.Stage, stageIndex int) []templates.Template {
	if s.image == "" {
		s.logger.Warn("stage monitor image wasn't specified; stage monitor creation skipped")
		return nil
	}

	podsToKill := make([]string, 0)
	ignoredNodes := make([]string, 0)

	for _, action := range stage.Actions {
		if action.Severity == metadata.SeverityCritical {
			if action.Scale == metadata.ScaleNode {
				ignoredNodes = append(ignoredNodes, action.Target.Node)
			} else {
				podsToKill = append(podsToKill, fmt.Sprintf("%s=%d", action.Target.AppLabel, -1))
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
			{Name: "IGNORED_NODES", Value: strings.Join(ignoredNodes, ";")},
		},
	})

	return []templates.Template{containerTemplate}
}

func UseStageMonitor(image string, targetNs string, bufferTime time.Duration, logger *zap.SugaredLogger) StageExtension {
	return stageMonitor{image: image, targetNs: targetNs, stageInterval: bufferTime, logger: logger}
}