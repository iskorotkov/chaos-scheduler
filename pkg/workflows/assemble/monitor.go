package assemble

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/rx"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generate"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

type monitor struct {
	image         string
	appLabel      string
	targetNs      string
	stageInterval time.Duration
	logger        *zap.SugaredLogger
}

func (s monitor) Generate(rand *rand.Rand, _ int) reflect.Value {
	return reflect.ValueOf(monitor{
		image:         rx.Rstr(rand, "image"),
		appLabel:      rx.Rstr(rand, "app-label"),
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

	ignoredDeployments := make([]string, 0)
	ignoredNodes := make([]string, 0)

	for _, action := range stage.Actions {
		if action.Severity == metadata.SeverityCritical {
			if action.Scale == metadata.ScaleNode {
				ignoredNodes = append(ignoredNodes, action.Target.Node)
			} else {
				ignoredDeployments = append(ignoredDeployments, action.Target.AppLabel)
			}
		}
	}

	name := fmt.Sprintf("monitor-%d", stageIndex+1)
	containerTemplate := templates.NewContainerTemplate(name, templates.Container{
		Name:  "monitor",
		Image: s.image,
		Env: []v1.EnvVar{
			{Name: "APP_NS", Value: s.targetNs},
			{Name: "APP_LABEL", Value: s.appLabel},
			{Name: "DURATION", Value: (stage.Duration + s.stageInterval).String()},
			{Name: "IGNORED_PODS", Value: ""},
			{Name: "IGNORED_DEPLOYMENTS", Value: strings.Join(ignoredDeployments, ";")},
			{Name: "IGNORED_NODES", Value: strings.Join(ignoredNodes, ";")},
		},
	})

	return []templates.Template{containerTemplate}
}

// UseMonitor returns a stage extension that adds a monitor app to the stage.
func UseMonitor(image, targetNs, appLabel string, bufferTime time.Duration, logger *zap.SugaredLogger) StageExt {
	return monitor{image: image, targetNs: targetNs, appLabel: appLabel, stageInterval: bufferTime, logger: logger}
}
