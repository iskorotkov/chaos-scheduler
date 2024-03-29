package assemble

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/rx"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generate"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
)

type monitor struct {
	image         string
	targetNs      string
	stageInterval time.Duration
	logger        *zap.SugaredLogger
}

func (s monitor) Generate(rand *rand.Rand, _ int) reflect.Value {
	return reflect.ValueOf(monitor{
		image:         rx.Rstr(rand, "image"),
		targetNs:      rx.Rstr(rand, "target-ns"),
		stageInterval: time.Duration(rand.Intn(200)) * time.Second,
		logger:        zap.NewNop().Sugar(),
	})
}

// Apply adds monitor app to the stage
func (s monitor) Apply(stage generate.Stage, stageIndex int) []templates.Template {
	if s.image == "" {
		s.logger.Warn("stage monitor image wasn't specified; stage monitor creation skipped")
		return nil
	}

	ignoredLabels := make([]string, 0)
	ignoredNodes := make([]string, 0)

	for _, action := range stage.Steps {
		if action.Severity == metadata.SeverityCritical {
			if action.Scale == metadata.ScaleNode {
				ignoredNodes = append(ignoredNodes, action.Target.Node)
			} else {
				ignoredLabels = append(ignoredLabels, action.Target.AppLabel)
			}
		}
	}

	name := fmt.Sprintf("monitor-%d", stageIndex+1)
	containerTemplate := templates.NewContainerTemplate(name, templates.Container{
		Name:  "monitor",
		Image: s.image,
		Env: []v1.EnvVar{
			{Name: "APP_NS", Value: s.targetNs},
			{Name: "DURATION", Value: (stage.Duration + s.stageInterval).String()},
			{Name: "IGNORED_PODS", Value: ""},
			{Name: "IGNORED_LABELS", Value: strings.Join(ignoredLabels, ";")},
			{Name: "IGNORED_NODES", Value: strings.Join(ignoredNodes, ";")},
		},
	})

	return []templates.Template{containerTemplate}
}

// UseMonitor returns a stage extension that adds a monitor app to the stage.
func UseMonitor(image, targetNs string, bufferTime time.Duration, logger *zap.SugaredLogger) StageExt {
	return monitor{image: image, targetNs: targetNs, stageInterval: bufferTime, logger: logger}
}
