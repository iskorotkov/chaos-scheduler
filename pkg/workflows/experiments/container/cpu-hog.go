package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type CPUHog struct {
	Namespace    string
	AppNamespace string
	Cores        int
}

func (c CPUHog) Engine(target targets.Target, duration time.Duration) experiments.Engine {
	return c.Instantiate(target.AppLabel, target.MainContainer, duration)
}

func (c CPUHog) Info() experiments.Info {
	return experiments.Info{
		Name:   "pod-cpu-hog",
		Lethal: false,
	}
}

func (c CPUHog) Instantiate(label string, container string, duration time.Duration) experiments.Engine {
	if c.Cores == 0 {
		c.Cores = 1
	}

	return experiments.NewEngine(experiments.EngineParams{
		Name:        c.Info().Name,
		Namespace:   c.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: experiments.AppInfo{
			AppNS:    c.AppNamespace,
			AppLabel: label,
			AppKind:  "deployment",
		},
		Experiments: []experiments.Experiment{
			experiments.NewExperiment(experiments.ExperimentParams{
				Name: c.Info().Name,
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"TARGET_CONTAINER":     container,
					"CPU_CORES":            strconv.Itoa(c.Cores),
				},
			}),
		},
	})
}
