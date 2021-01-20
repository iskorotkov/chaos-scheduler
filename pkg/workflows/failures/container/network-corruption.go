package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type NetworkCorruption struct {
	Namespace            string
	AppNamespace         string
	CorruptionPercentage int
}

func (n NetworkCorruption) Instantiate(target targets.Target, duration time.Duration) failures.Engine {
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
					"TOTAL_CHAOS_DURATION":                 strconv.Itoa(int(duration.Seconds())),
					"NETWORK_INTERFACE":                    "eth0",
					"TARGET_CONTAINER":                     target.MainContainer,
					"NETWORK_PACKET_CORRUPTION_PERCENTAGE": strconv.Itoa(n.CorruptionPercentage),
				},
			}),
		},
	})
}

func (n NetworkCorruption) Name() string {
	return "pod-network-corruption"
}
