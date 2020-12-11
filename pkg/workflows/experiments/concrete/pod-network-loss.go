package concrete

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"strconv"
	"time"
)

type PodNetworkLoss struct {
	Namespace      string
	AppNamespace   string
	LossPercentage int
}

func (p PodNetworkLoss) Info() experiments.Info {
	return experiments.Info{Lethal: false}
}

func (p PodNetworkLoss) Instantiate(label string, container string, duration time.Duration) experiments.Engine {
	return experiments.NewEngine(experiments.EngineParams{
		Name:        string(p.Type()),
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
				Type: p.Type(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION":           strconv.Itoa(int(duration.Seconds())),
					"NETWORK_INTERFACE":              "eth0",
					"TARGET_CONTAINER":               container,
					"NETWORK_PACKET_LOSS_PERCENTAGE": strconv.Itoa(p.LossPercentage),
				},
			}),
		},
	})
}

func (p PodNetworkLoss) Type() experiments.ExperimentType {
	return "pod-network-loss"
}
