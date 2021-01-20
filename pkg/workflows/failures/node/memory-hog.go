package node

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type MemoryHog struct {
	Namespace        string
	AppNamespace     string
	MemoryPercentage int
}

func (m MemoryHog) Instantiate(target targets.Target, duration time.Duration) failures.Engine {
	if m.MemoryPercentage == 0 {
		m.MemoryPercentage = 90
	}

	return failures.NewEngine(failures.EngineParams{
		Name:        m.Name(),
		Namespace:   m.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: failures.AppInfo{
			AppNS:    m.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []failures.Experiment{
			failures.NewExperiment(failures.ExperimentParams{
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
