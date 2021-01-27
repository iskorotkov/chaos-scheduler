package pod

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
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

func (p Delete) Instantiate(target targets.Target, duration time.Duration) blueprints.Engine {
	return blueprints.NewEngine(blueprints.EngineParams{
		Name:        p.Name(),
		Namespace:   p.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: blueprints.AppInfo{
			AppNS:    p.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []blueprints.Experiment{
			blueprints.NewExperiment(blueprints.ExperimentParams{
				Name: p.Name(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"CHAOS_INTERVAL":       strconv.Itoa(p.Interval),
					"FORCE":                strconv.FormatBool(p.Force),
				},
			}),
		},
	})
}

func (p Delete) Name() string {
	return "pod-delete"
}
