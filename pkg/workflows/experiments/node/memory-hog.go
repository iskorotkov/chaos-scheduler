package node

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"strconv"
	"time"
)

type MemoryHog struct {
	Namespace        string
	AppNamespace     string
	MemoryPercentage int
}

func (m MemoryHog) Type() experiments.ExperimentType {
	return "node-memory-hog"
}

func (m MemoryHog) Info() experiments.Info {
	return experiments.Info{Lethal: false}
}

func (m MemoryHog) Instantiate(label string, node string, duration time.Duration) experiments.Engine {
	if m.MemoryPercentage == 0 {
		m.MemoryPercentage = 90
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
					"TARGET_NODES":         node,
					"MEMORY_PERCENTAGE":    strconv.Itoa(m.MemoryPercentage),
				},
			}),
		},
	})
}
