package container

import (
	"strconv"
	"time"

	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
)

type NetworkDuplication struct {
	Namespace              string
	AppNamespace           string
	DuplicationPercentage  int
	PodsAffectedPercentage int
}

func (n NetworkDuplication) Instantiate(target targets.Target, duration time.Duration) blueprints.Engine {
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
					"TOTAL_CHAOS_DURATION":                  strconv.Itoa(int(duration.Seconds())),
					"NETWORK_INTERFACE":                     "eth0",
					"TARGET_CONTAINER":                      target.MainContainer,
					"NETWORK_PACKET_DUPLICATION_PERCENTAGE": strconv.Itoa(n.DuplicationPercentage),
					"PODS_AFFECTED_PERC":                    strconv.Itoa(n.PodsAffectedPercentage),
					"CONTAINER_RUNTIME":                     "containerd",
					"SOCKET_PATH":                           "/run/containerd/containerd.sock",
				},
			}),
		},
	})
}

func (n NetworkDuplication) Name() string {
	return "pod-network-duplication"
}

func (n NetworkDuplication) Type() blueprints.BlueprintType {
	return blueprints.BlueprintTypeNetwork
}
