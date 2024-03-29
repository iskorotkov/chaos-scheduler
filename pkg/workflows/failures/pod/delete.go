package pod

import (
	"strconv"
	"time"

	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
)

type Delete struct {
	Namespace              string
	AppNamespace           string
	Interval               int
	PodsAffectedPercentage int
	Force                  bool
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
					"PODS_AFFECTED_PERC":   strconv.Itoa(p.PodsAffectedPercentage),
				},
			}),
		},
	})
}

func (p Delete) Name() string {
	return "pod-delete"
}

func (p Delete) Type() blueprints.BlueprintType {
	return blueprints.BlueprintTypeRestart
}
