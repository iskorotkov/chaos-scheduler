package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type MemoryHog struct {
	Namespace         string
	AppNamespace      string
	MemoryConsumption int
}

func (m MemoryHog) Instantiate(target targets.Target, duration time.Duration) blueprints.Engine {
	if m.MemoryConsumption == 0 {
		m.MemoryConsumption = 500
	}

	return blueprints.NewEngine(blueprints.EngineParams{
		Name:        m.Name(),
		Namespace:   m.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: blueprints.AppInfo{
			AppNS:    m.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []blueprints.Experiment{
			blueprints.NewExperiment(blueprints.ExperimentParams{
				Name: m.Name(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"TARGET_CONTAINER":     target.MainContainer,
					"MEMORY_CONSUMPTION":   strconv.Itoa(m.MemoryConsumption),
				},
			}),
		},
	})
}

func (m MemoryHog) Name() string {
	return "pod-memory-hog"
}
