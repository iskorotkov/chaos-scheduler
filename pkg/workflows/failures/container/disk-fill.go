package container

import (
	"strconv"
	"time"

	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
)

type DiskFill struct {
	Namespace              string
	AppNamespace           string
	FillPercentage         int
	PodsAffectedPercentage int
}

func (d DiskFill) Instantiate(target targets.Target, duration time.Duration) blueprints.Engine {
	return blueprints.NewEngine(blueprints.EngineParams{
		Name:        d.Name(),
		Namespace:   d.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: blueprints.AppInfo{
			AppNS:    d.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []blueprints.Experiment{
			blueprints.NewExperiment(blueprints.ExperimentParams{
				Name: d.Name(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"TARGET_CONTAINER":     target.MainContainer,
					"FILL_PERCENTAGE":      strconv.Itoa(d.FillPercentage),
					"PODS_AFFECTED_PERC":   strconv.Itoa(d.PodsAffectedPercentage),
				},
			}),
		},
	})
}

func (d DiskFill) Name() string {
	return "disk-fill"
}

func (d DiskFill) Type() blueprints.BlueprintType {
	return blueprints.BlueprintTypeIO
}
