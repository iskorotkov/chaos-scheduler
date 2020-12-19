package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type NetworkLatency struct {
	Namespace      string
	AppNamespace   string
	NetworkLatency int
}

func (p NetworkLatency) Engine(target targets.Target, duration time.Duration) experiments.Engine {
	return p.Instantiate(target.AppLabel, target.MainContainer, duration)
}

func (p NetworkLatency) Info() experiments.Info {
	return experiments.Info{
		Name:          "pod-network-latency",
		Lethal:        false,
		AffectingNode: false,
	}
}

func (p NetworkLatency) Instantiate(label string, container string, duration time.Duration) experiments.Engine {
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
					"NETWORK_INTERFACE":    "eth0",
					"TARGET_CONTAINER":     container,
					"NETWORK_LATENCY":      strconv.Itoa(p.NetworkLatency),
				},
			}),
		},
	})
}
