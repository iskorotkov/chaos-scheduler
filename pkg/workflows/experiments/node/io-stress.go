package node

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type IOStress struct {
	Namespace             string
	AppNamespace          string
	UtilizationPercentage int
}

func (i IOStress) Engine(target targets.Target, duration time.Duration) experiments.Engine {
	return i.Instantiate(target.AppLabel, target.Node, duration)
}

func (i IOStress) Info() experiments.Info {
	return experiments.Info{
		Name:          "node-io-stress",
		Lethal:        false,
		AffectingNode: true,
	}
}

func (i IOStress) Instantiate(label string, node string, duration time.Duration) experiments.Engine {
	if i.UtilizationPercentage == 0 {
		i.UtilizationPercentage = 10
	}

	return experiments.NewEngine(experiments.EngineParams{
		Name:        i.Info().Name,
		Namespace:   i.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: experiments.AppInfo{
			AppNS:    i.AppNamespace,
			AppLabel: label,
			AppKind:  "deployment",
		},
		Experiments: []experiments.Experiment{
			experiments.NewExperiment(experiments.ExperimentParams{
				Name: i.Info().Name,
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION":              strconv.Itoa(int(duration.Seconds())),
					"TARGET_NODES":                      node,
					"FILESYSTEM_UTILIZATION_PERCENTAGE": strconv.Itoa(i.UtilizationPercentage),
				},
			}),
		},
	})
}