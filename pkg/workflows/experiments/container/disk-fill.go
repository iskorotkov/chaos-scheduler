package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type DiskFill struct {
	Namespace      string
	AppNamespace   string
	FillPercentage int
}

func (d DiskFill) Engine(target targets.Target, duration time.Duration) experiments.Engine {
	return d.Instantiate(target.AppLabel, target.MainContainer, duration)
}

func (d DiskFill) Info() experiments.Info {
	return experiments.Info{
		Name:          "disk-fill",
		Lethal:        false,
		AffectingNode: false,
	}
}

func (d DiskFill) Instantiate(label string, container string, duration time.Duration) experiments.Engine {
	return experiments.NewEngine(experiments.EngineParams{
		Name:        d.Info().Name,
		Namespace:   d.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: experiments.AppInfo{
			AppNS:    d.AppNamespace,
			AppLabel: label,
			AppKind:  "deployment",
		},
		Experiments: []experiments.Experiment{
			experiments.NewExperiment(experiments.ExperimentParams{
				Name: d.Info().Name,
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"TARGET_CONTAINER":     container,
					"FILL_PERCENTAGE":      strconv.Itoa(d.FillPercentage),
				},
			}),
		},
	})
}
