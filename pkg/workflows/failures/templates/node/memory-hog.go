package node

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type MemoryHog struct {
	Namespace        string
	AppNamespace     string
	MemoryPercentage int
}

func (m MemoryHog) Instantiate(target targets.Target, duration time.Duration) templates.Engine {
	if m.MemoryPercentage == 0 {
		m.MemoryPercentage = 90
	}

	return templates.NewEngine(templates.EngineParams{
		Name:        m.Name(),
		Namespace:   m.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: templates.AppInfo{
			AppNS:    m.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []templates.Experiment{
			templates.NewExperiment(templates.ExperimentParams{
				Name: m.Name(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"TARGET_NODES":         target.Node,
					"MEMORY_PERCENTAGE":    strconv.Itoa(m.MemoryPercentage),
				},
			}),
		},
	})
}

func (m MemoryHog) Name() string {
	return "node-memory-hog"
}
