package pod

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type Delete struct {
	Namespace    string
	AppNamespace string
	Interval     int
	Force        bool
}

func (p Delete) Engine(target targets.Target, duration time.Duration) experiments.Engine {
	return p.Instantiate(target.Selector(), duration)
}

func (p Delete) Info() experiments.Info {
	return experiments.Info{
		Name:   "pod-delete",
		Lethal: true,
	}
}

func (p Delete) Instantiate(label string, duration time.Duration) experiments.Engine {
	return experiments.NewEngine(experiments.EngineParams{
		Name:        p.Info().Name,
		Namespace:   p.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: experiments.AppInfo{
			AppNS:    p.AppNamespace,
			AppLabel: label,
			AppKind:  "deployment",
		},
		Experiments: []experiments.Experiment{
			experiments.NewExperiment(experiments.ExperimentParams{
				Name: p.Info().Name,
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"CHAOS_INTERVAL":       strconv.Itoa(p.Interval),
					"FORCE":                strconv.FormatBool(p.Force),
				},
			}),
		},
	})
}
