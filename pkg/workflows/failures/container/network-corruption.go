package container

import (
	"strconv"
	"time"

	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
)

type NetworkCorruption struct {
	Namespace              string
	AppNamespace           string
	CorruptionPercentage   int
	PodsAffectedPercentage int
}

func (n NetworkCorruption) Instantiate(target targets.Target, duration time.Duration) blueprints.Engine {
	return blueprints.NewEngine(blueprints.EngineParams{
		Name:        n.Name(),
		Namespace:   n.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: blueprints.AppInfo{
			AppNS:    n.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []blueprints.Experiment{
			blueprints.NewExperiment(blueprints.ExperimentParams{
				Name: n.Name(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION":                 strconv.Itoa(int(duration.Seconds())),
					"NETWORK_INTERFACE":                    "eth0",
					"TARGET_CONTAINER":                     target.MainContainer,
					"NETWORK_PACKET_CORRUPTION_PERCENTAGE": strconv.Itoa(n.CorruptionPercentage),
					"PODS_AFFECTED_PERC":                   strconv.Itoa(n.PodsAffectedPercentage),
					"CONTAINER_RUNTIME":                    "containerd",
					"SOCKET_PATH":                          "/run/containerd/containerd.sock",
				},
			}),
		},
	})
}

func (n NetworkCorruption) Name() string {
	return "pod-network-corruption"
}

func (n NetworkCorruption) Type() blueprints.BlueprintType {
	return blueprints.BlueprintTypeNetwork
}
