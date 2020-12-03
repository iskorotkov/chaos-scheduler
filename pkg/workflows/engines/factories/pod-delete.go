package factories

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/engines"
	"strconv"
)

type PodDeleteFactory struct {
	Namespace       string
	TargetNamespace string
	Duration        int
	Interval        int
	Force           bool
}

func (p PodDeleteFactory) Create(target string) engines.Engine {
	return engines.NewEngine(engines.EngineParams{
		Name:        string(p.Type()),
		Namespace:   p.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo:     engines.AppInfo{AppNS: p.TargetNamespace, AppLabel: "app=server", AppKind: "deployment"},
		Experiments: []engines.Experiment{
			engines.NewExperiment(engines.ExperimentParams{Type: p.Type(), Env: map[string]string{
				"TOTAL_CHAOS_DURATION": strconv.Itoa(p.Duration),
				"CHAOS_INTERVAL":       strconv.Itoa(p.Interval),
				"FORCE":                strconv.FormatBool(p.Force),
			}}),
		},
	})
}

func (p PodDeleteFactory) Type() engines.ExperimentType {
	return "pod-delete"
}
