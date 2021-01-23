package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/templates"
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

func (p NetworkLatency) Instantiate(target targets.Target, duration time.Duration) templates.Engine {
	return templates.NewEngine(templates.EngineParams{
		Name:        p.Name(),
		Namespace:   p.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: templates.AppInfo{
			AppNS:    p.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []templates.Experiment{
			templates.NewExperiment(templates.ExperimentParams{
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
