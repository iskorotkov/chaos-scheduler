package node

import (
	"strconv"
	"time"

	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
)

type CPUHog struct {
	Namespace    string
	AppNamespace string
	Cores        int
}

func (c CPUHog) Instantiate(target targets.Target, duration time.Duration) blueprints.Engine {
	if c.Cores == 0 {
		c.Cores = 2
	}

	return blueprints.NewEngine(blueprints.EngineParams{
		Name:        c.Name(),
		Namespace:   c.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: blueprints.AppInfo{
			AppNS:    c.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []blueprints.Experiment{
			blueprints.NewExperiment(blueprints.ExperimentParams{
				Name: c.Name(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"TARGET_NODES":         target.Node,
					"NODE_CPU_CORE":        strconv.Itoa(c.Cores),
				},
			}),
		},
	})
}

func (c CPUHog) Name() string {
	return "node-cpu-hog"
}

func (c CPUHog) Type() blueprints.BlueprintType {
	return blueprints.BlueprintTypeResources
}
