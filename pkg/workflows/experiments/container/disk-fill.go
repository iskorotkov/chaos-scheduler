package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"strconv"
	"time"
)

type DiskFill struct {
	Namespace      string
	AppNamespace   string
	FillPercentage int
}

func (d DiskFill) Info() experiments.Info {
	return experiments.Info{Lethal: false}
}

func (d DiskFill) Instantiate(label string, container string, duration time.Duration) experiments.Engine {
	return experiments.NewEngine(experiments.EngineParams{
		Name:        string(d.Type()),
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
				Type: d.Type(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"TARGET_CONTAINER":     container,
					"FILL_PERCENTAGE":      strconv.Itoa(d.FillPercentage),
				},
			}),
		},
	})
}

func (d DiskFill) Type() experiments.ExperimentType {
	return "disk-fill"
}
