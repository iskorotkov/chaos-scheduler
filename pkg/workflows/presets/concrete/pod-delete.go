package concrete

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/presets"
	"strconv"
)

type PodDelete struct {
	Namespace    string
	AppNamespace string
	Duration     int
	Interval     int
	Force        bool
}

func (p PodDelete) Instantiate(label string) presets.Engine {
	return presets.NewEngine(presets.EngineParams{
		Name:        string(p.Type()),
		Namespace:   p.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: presets.AppInfo{
			AppNS:    p.AppNamespace,
			AppLabel: label,
			AppKind:  "deployment",
		},
		Experiments: []presets.Experiment{
			presets.NewExperiment(presets.ExperimentParams{Type: p.Type(), Env: map[string]string{
				"TOTAL_CHAOS_DURATION": strconv.Itoa(p.Duration),
				"CHAOS_INTERVAL":       strconv.Itoa(p.Interval),
				"FORCE":                strconv.FormatBool(p.Force),
			}}),
		},
	})
}

func (p PodDelete) Type() presets.ExperimentType {
	return "pod-delete"
}
