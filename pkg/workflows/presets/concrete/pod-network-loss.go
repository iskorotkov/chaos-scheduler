package concrete

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/presets"
	"strconv"
)

type PodNetworkLoss struct {
	Namespace      string
	AppNamespace   string
	LossPercentage int
}

func (p PodNetworkLoss) Info() presets.Info {
	return presets.Info{Lethal: false}
}

func (p PodNetworkLoss) Instantiate(label string, container string) presets.Engine {
	return presets.NewEngine(presets.EngineParams{
		Name:        string(p.Type()),
		Namespace:   p.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: presets.AppInfo{
			AppNS:    p.AppNamespace,
			AppLabel: label,
			AppKind:  "deployment",
		},
		Experiments: []presets.Experiment{
			presets.NewExperiment(presets.ExperimentParams{
				Type: p.Type(),
				Env: map[string]string{
					"NETWORK_INTERFACE":              "eth0",
					"TARGET_CONTAINER":               container,
					"NETWORK_PACKET_LOSS_PERCENTAGE": strconv.Itoa(p.LossPercentage),
				},
			}),
		},
	})
}

func (p PodNetworkLoss) Type() presets.ExperimentType {
	return "pod-network-loss"
}
