package node

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"strconv"
	"time"
)

type IOStress struct {
	Namespace             string
	AppNamespace          string
	UtilizationPercentage int
}

func (i IOStress) Type() experiments.ExperimentType {
	return "node-io-stress"
}

func (i IOStress) Info() experiments.Info {
	return experiments.Info{Lethal: false}
}

func (i IOStress) Instantiate(label string, node string, duration time.Duration) experiments.Engine {
	if i.UtilizationPercentage == 0 {
		i.UtilizationPercentage = 10
	}

	return experiments.NewEngine(experiments.EngineParams{
		Name:        string(i.Type()),
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
				Type: i.Type(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION":              strconv.Itoa(int(duration.Seconds())),
					"TARGET_NODES":                      node,
					"FILESYSTEM_UTILIZATION_PERCENTAGE": strconv.Itoa(i.UtilizationPercentage),
				},
			}),
		},
	})
}
