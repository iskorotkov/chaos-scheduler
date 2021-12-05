package container

import (
	"strconv"
	"time"

	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
)

type MemoryHog struct {
	Namespace              string
	AppNamespace           string
	MemoryConsumption      int
	PodsAffectedPercentage int
}

func (m MemoryHog) Instantiate(target targets.Target, duration time.Duration) blueprints.Engine {
	if m.MemoryConsumption == 0 {
		m.MemoryConsumption = 500
	}

	return blueprints.NewEngine(blueprints.EngineParams{
		Name:        m.Name(),
		Namespace:   m.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: blueprints.AppInfo{
			AppNS:    m.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []blueprints.Experiment{
			blueprints.NewExperiment(blueprints.ExperimentParams{
				Name: m.Name(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"TARGET_CONTAINER":     target.MainContainer,
					"MEMORY_CONSUMPTION":   strconv.Itoa(m.MemoryConsumption),
					"PODS_AFFECTED_PERC":   strconv.Itoa(m.PodsAffectedPercentage),
					"CONTAINER_RUNTIME":    "containerd",
					"SOCKET_PATH":          "/run/containerd/containerd.sock",
				},
			}),
		},
	})
}

func (m MemoryHog) Name() string {
	return "pod-memory-hog"
}
