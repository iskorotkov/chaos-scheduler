package pod

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type IOStress struct {
	Namespace             string
	AppNamespace          string
	UtilizationPercentage int
}

func (i IOStress) Instantiate(target targets.Target, duration time.Duration) blueprints.Engine {
	if i.UtilizationPercentage == 0 {
		i.UtilizationPercentage = 10
	}

	return blueprints.NewEngine(blueprints.EngineParams{
		Name:        i.Name(),
		Namespace:   i.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: blueprints.AppInfo{
			AppNS:    i.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []blueprints.Experiment{
			blueprints.NewExperiment(blueprints.ExperimentParams{
				Name: i.Name(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION":              strconv.Itoa(int(duration.Seconds())),
					"FILESYSTEM_UTILIZATION_PERCENTAGE": strconv.Itoa(i.UtilizationPercentage),
				},
			}),
		},
	})
}

func (i IOStress) Name() string {
	return "pod-io-stress"
}
