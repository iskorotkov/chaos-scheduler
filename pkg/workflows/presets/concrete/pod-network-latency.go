package concrete

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/presets"
	"strconv"
)

type PodNetworkLatency struct {
	Namespace      string
	AppNamespace   string
	NetworkLatency int
}

func (p PodNetworkLatency) Info() presets.Info {
	return presets.Info{Lethal: false}
}

func (p PodNetworkLatency) Instantiate(label string, container string) presets.Engine {
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
					"NETWORK_INTERFACE": "eth0",
					"TARGET_CONTAINER":  container,
					"NETWORK_LATENCY":   strconv.Itoa(p.NetworkLatency),
				},
			}),
		},
	})
}

func (p PodNetworkLatency) Type() presets.ExperimentType {
	return "pod-network-latency"
}
