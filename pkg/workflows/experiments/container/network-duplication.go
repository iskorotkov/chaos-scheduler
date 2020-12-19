package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type NetworkDuplication struct {
	Namespace             string
	AppNamespace          string
	DuplicationPercentage int
}

func (n NetworkDuplication) Engine(target targets.Target, duration time.Duration) experiments.Engine {
	return n.Instantiate(target.AppLabel, target.MainContainer, duration)
}

func (n NetworkDuplication) Info() experiments.Info {
	return experiments.Info{
		Name:          "pod-network-duplication",
		Lethal:        false,
		AffectingNode: false,
	}
}

func (n NetworkDuplication) Instantiate(label string, container string, duration time.Duration) experiments.Engine {
	return experiments.NewEngine(experiments.EngineParams{
		Name:        n.Info().Name,
		Namespace:   n.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: experiments.AppInfo{
			AppNS:    n.AppNamespace,
			AppLabel: label,
			AppKind:  "deployment",
		},
		Experiments: []experiments.Experiment{
			experiments.NewExperiment(experiments.ExperimentParams{
				Name: n.Info().Name,
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION":                  strconv.Itoa(int(duration.Seconds())),
					"NETWORK_INTERFACE":                     "eth0",
					"TARGET_CONTAINER":                      container,
					"NETWORK_PACKET_DUPLICATION_PERCENTAGE": strconv.Itoa(n.DuplicationPercentage),
				},
			}),
		},
	})
}
