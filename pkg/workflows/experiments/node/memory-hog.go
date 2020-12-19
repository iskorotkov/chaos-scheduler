package node

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type MemoryHog struct {
	Namespace        string
	AppNamespace     string
	MemoryPercentage int
}

func (m MemoryHog) Engine(target targets.Target, duration time.Duration) experiments.Engine {
	return m.Instantiate(target.AppLabel, target.Node, duration)
}

func (m MemoryHog) Info() experiments.Info {
	return experiments.Info{
		Name:          "node-memory-hog",
		Lethal:        false,
		AffectingNode: true,
	}
}

func (m MemoryHog) Instantiate(label string, node string, duration time.Duration) experiments.Engine {
	if m.MemoryPercentage == 0 {
		m.MemoryPercentage = 90
	}

	return experiments.NewEngine(experiments.EngineParams{
		Name:        m.Info().Name,
		Namespace:   m.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: experiments.AppInfo{
			AppNS:    m.AppNamespace,
			AppLabel: label,
			AppKind:  "deployment",
		},
		Experiments: []experiments.Experiment{
			experiments.NewExperiment(experiments.ExperimentParams{
				Name: m.Info().Name,
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"TARGET_NODES":         node,
					"MEMORY_PERCENTAGE":    strconv.Itoa(m.MemoryPercentage),
				},
			}),
		},
	})
}
