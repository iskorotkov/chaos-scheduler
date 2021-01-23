package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type NetworkCorruption struct {
	Namespace            string
	AppNamespace         string
	CorruptionPercentage int
}

func (n NetworkCorruption) Instantiate(target targets.Target, duration time.Duration) templates.Engine {
	return templates.NewEngine(templates.EngineParams{
		Name:        n.Name(),
		Namespace:   n.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: templates.AppInfo{
			AppNS:    n.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []templates.Experiment{
			templates.NewExperiment(templates.ExperimentParams{
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
