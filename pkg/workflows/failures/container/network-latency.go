package container

import (
	"strconv"
	"time"

	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
)

type NetworkLatency struct {
	Namespace              string
	AppNamespace           string
	NetworkLatency         int
	PodsAffectedPercentage int
}

func (p NetworkLatency) Name() string {
	return "pod-network-latency"
}

func (p NetworkLatency) Instantiate(target targets.Target, duration time.Duration) blueprints.Engine {
	return blueprints.NewEngine(blueprints.EngineParams{
		Name:        p.Name(),
		Namespace:   p.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: blueprints.AppInfo{
			AppNS:    p.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []blueprints.Experiment{
			blueprints.NewExperiment(blueprints.ExperimentParams{
				Name: p.Name(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"NETWORK_INTERFACE":    "eth0",
					"TARGET_CONTAINER":     target.MainContainer,
					"NETWORK_LATENCY":      strconv.Itoa(p.NetworkLatency),
					"PODS_AFFECTED_PERC":   strconv.Itoa(p.PodsAffectedPercentage),
					"CONTAINER_RUNTIME":    "containerd",
					"SOCKET_PATH":          "/run/containerd/containerd.sock",
				},
			}),
		},
	})
}
