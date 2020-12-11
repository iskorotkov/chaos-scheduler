package concrete

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"strconv"
	"time"
)

type PodNetworkLatency struct {
	Namespace      string
	AppNamespace   string
	NetworkLatency int
}

func (p PodNetworkLatency) Info() experiments.Info {
	return experiments.Info{Lethal: false}
}

func (p PodNetworkLatency) Instantiate(label string, container string, duration time.Duration) experiments.Engine {
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
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"NETWORK_INTERFACE":    "eth0",
					"TARGET_CONTAINER":     container,
					"NETWORK_LATENCY":      strconv.Itoa(p.NetworkLatency),
				},
			}),
		},
	})
}

func (p PodNetworkLatency) Type() experiments.ExperimentType {
	return "pod-network-latency"
}
