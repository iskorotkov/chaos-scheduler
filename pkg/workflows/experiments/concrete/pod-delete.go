package concrete

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"strconv"
)

type PodDelete struct {
	Namespace    string
	AppNamespace string
	Duration     int
	Interval     int
	Force        bool
}

func (p PodDelete) Info() experiments.Info {
	return experiments.Info{Lethal: true}
}

func (p PodDelete) Instantiate(label string) experiments.Engine {
	return experiments.NewEngine(experiments.EngineParams{
		Name:        string(p.Type()),
		Namespace:   p.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: experiments.AppInfo{
			AppNS:    p.AppNamespace,
			AppLabel: label,
			AppKind:  "deployment",
		},
		Experiments: []experiments.Experiment{
			experiments.NewExperiment(experiments.ExperimentParams{Type: p.Type(), Env: map[string]string{
				"TOTAL_CHAOS_DURATION": strconv.Itoa(p.Duration),
				"CHAOS_INTERVAL":       strconv.Itoa(p.Interval),
				"FORCE":                strconv.FormatBool(p.Force),
			}}),
		},
	})
}

func (p PodDelete) Type() experiments.ExperimentType {
	return "pod-delete"
}
