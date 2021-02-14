package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type CPUHog struct {
	Namespace              string
	AppNamespace           string
	Cores                  int
	PodsAffectedPercentage int
}

func (c CPUHog) Instantiate(target targets.Target, duration time.Duration) blueprints.Engine {
	if c.Cores == 0 {
		c.Cores = 1
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
					"TARGET_CONTAINER":     target.MainContainer,
					"CPU_CORES":            strconv.Itoa(c.Cores),
					"PODS_AFFECTED_PERC":   strconv.Itoa(c.PodsAffectedPercentage),
				},
			}),
		},
	})
}

func (c CPUHog) Name() string {
	return "pod-cpu-hog"
}
