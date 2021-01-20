package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type CPUHog struct {
	Namespace    string
	AppNamespace string
	Cores        int
}

func (c CPUHog) Instantiate(target targets.Target, duration time.Duration) failures.Engine {
	if c.Cores == 0 {
		c.Cores = 1
	}

	return failures.NewEngine(failures.EngineParams{
		Name:        c.Name(),
		Namespace:   c.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: failures.AppInfo{
			AppNS:    c.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []failures.Experiment{
			failures.NewExperiment(failures.ExperimentParams{
				Name: c.Name(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"TARGET_CONTAINER":     target.MainContainer,
					"CPU_CORES":            strconv.Itoa(c.Cores),
				},
			}),
		},
	})
}

func (c CPUHog) Name() string {
	return "pod-cpu-hog"
}
