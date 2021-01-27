package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type NetworkDuplication struct {
	Namespace             string
	AppNamespace          string
	DuplicationPercentage int
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
				},
			}),
		},
	})
}

func (n NetworkDuplication) Name() string {
	return "pod-network-duplication"
}
