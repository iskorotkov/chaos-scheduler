package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type NetworkLatency struct {
	Namespace      string
	AppNamespace   string
	NetworkLatency int
}

func (p NetworkLatency) Name() string {
	return "pod-network-latency"
}

func (p NetworkLatency) Instantiate(target targets.Target, duration time.Duration) failures.Engine {
	return failures.NewEngine(failures.EngineParams{
		Name:        p.Name(),
		Namespace:   p.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: failures.AppInfo{
			AppNS:    p.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []failures.Experiment{
			failures.NewExperiment(failures.ExperimentParams{
				Name: p.Name(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"NETWORK_INTERFACE":    "eth0",
					"TARGET_CONTAINER":     target.MainContainer,
					"NETWORK_LATENCY":      strconv.Itoa(p.NetworkLatency),
				},
			}),
		},
	})
}
