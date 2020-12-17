package experiments

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"time"
)

type Info struct {
	Name   string
	Lethal bool
}

type Preset interface {
	Info() Info
	Engine(target targets.Target, duration time.Duration) Engine
}

type PodPreset interface {
	Preset
	Instantiate(label string, duration time.Duration) Engine
}

type ContainerPreset interface {
	Preset
	Instantiate(label string, container string, duration time.Duration) Engine
}

type NodePreset interface {
	Preset
	Instantiate(label string, node string, duration time.Duration) Engine
}
