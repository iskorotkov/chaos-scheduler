package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type NetworkDuplication struct {
	Namespace             string
	AppNamespace          string
	DuplicationPercentage int
}

func (n NetworkDuplication) Instantiate(target targets.Target, duration time.Duration) templates.Engine {
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
