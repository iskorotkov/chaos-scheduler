package node

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type CPUHog struct {
	Namespace    string
	AppNamespace string
	Cores        int
}

func (c CPUHog) Instantiate(target targets.Target, duration time.Duration) templates.Engine {
	if c.Cores == 0 {
		c.Cores = 2
	}

	return templates.NewEngine(templates.EngineParams{
		Name:        c.Name(),
		Namespace:   c.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: templates.AppInfo{
			AppNS:    c.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []templates.Experiment{
			templates.NewExperiment(templates.ExperimentParams{
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
