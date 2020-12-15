package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"strconv"
	"time"
)

type NetworkCorruption struct {
	Namespace            string
	AppNamespace         string
	CorruptionPercentage int
}

func (n NetworkCorruption) Info() experiments.Info {
	return experiments.Info{Lethal: false}
}

func (n NetworkCorruption) Instantiate(label string, container string, duration time.Duration) experiments.Engine {
	return experiments.NewEngine(experiments.EngineParams{
		Name:        string(n.Type()),
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
				Type: n.Type(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION":                 strconv.Itoa(int(duration.Seconds())),
					"NETWORK_INTERFACE":                    "eth0",
					"TARGET_CONTAINER":                     container,
					"NETWORK_PACKET_CORRUPTION_PERCENTAGE": strconv.Itoa(n.CorruptionPercentage),
				},
			}),
		},
	})
}

func (n NetworkCorruption) Type() experiments.ExperimentType {
	return "pod-network-corruption"
}