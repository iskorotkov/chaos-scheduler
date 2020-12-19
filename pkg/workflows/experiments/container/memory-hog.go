package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type MemoryHog struct {
	Namespace         string
	AppNamespace      string
	MemoryConsumption int
}

func (m MemoryHog) Engine(target targets.Target, duration time.Duration) experiments.Engine {
	return m.Instantiate(target.AppLabel, target.MainContainer, duration)
}

func (m MemoryHog) Info() experiments.Info {
	return experiments.Info{
		Name:          "pod-memory-hog",
		Lethal:        false,
		AffectingNode: false,
	}
}

func (m MemoryHog) Instantiate(label string, container string, duration time.Duration) experiments.Engine {
	if m.MemoryConsumption == 0 {
		m.MemoryConsumption = 500
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
					"TARGET_CONTAINER":     container,
					"MEMORY_CONSUMPTION":   strconv.Itoa(m.MemoryConsumption),
				},
			}),
		},
	})
}
