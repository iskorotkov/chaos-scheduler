package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type NetworkDuplication struct {
	Namespace             string
	AppNamespace          string
	DuplicationPercentage int
}

func (n NetworkDuplication) Instantiate(target targets.Target, duration time.Duration) failures.Engine {
	return failures.NewEngine(failures.EngineParams{
		Name:        n.Name(),
		Namespace:   n.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: failures.AppInfo{
			AppNS:    n.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []failures.Experiment{
			failures.NewExperiment(failures.ExperimentParams{
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
