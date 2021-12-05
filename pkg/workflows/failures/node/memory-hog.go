package node

import (
	"strconv"
	"time"

	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
)

type MemoryHog struct {
	Namespace        string
	AppNamespace     string
	MemoryPercentage int
}

func (m MemoryHog) Instantiate(target targets.Target, duration time.Duration) blueprints.Engine {
	if m.MemoryPercentage == 0 {
		m.MemoryPercentage = 90
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
					"TOTAL_CHAOS_DURATION":          strconv.Itoa(int(duration.Seconds())),
					"TARGET_NODES":                  target.Node,
					"MEMORY_CONSUMPTION_PERCENTAGE": strconv.Itoa(m.MemoryPercentage),
				},
			}),
		},
	})
}

func (m MemoryHog) Name() string {
	return "node-memory-hog"
}

func (m MemoryHog) Type() blueprints.BlueprintType {
	return blueprints.BlueprintTypeResources
}
