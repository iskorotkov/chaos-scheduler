package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"strconv"
	"time"
)

type MemoryHog struct {
	Namespace         string
	AppNamespace      string
	MemoryConsumption int
}

func (m MemoryHog) Info() experiments.Info {
	return experiments.Info{Lethal: false}
}

func (m MemoryHog) Instantiate(label string, container string, duration time.Duration) experiments.Engine {
	if m.MemoryConsumption == 0 {
		m.MemoryConsumption = 500
	}

	return experiments.NewEngine(experiments.EngineParams{
		Name:        string(m.Type()),
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
				Type: m.Type(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"TARGET_CONTAINER":     container,
					"MEMORY_CONSUMPTION":   strconv.Itoa(m.MemoryConsumption),
				},
			}),
		},
	})
}

func (m MemoryHog) Type() experiments.ExperimentType {
	return "pod-memory-hog"
}
