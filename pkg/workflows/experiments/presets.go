package experiments

import "time"

type Info struct {
	Lethal bool
}

type EnginePreset interface {
	Type() ExperimentType
	Info() Info
}

type PodEnginePreset interface {
	EnginePreset
	Instantiate(label string, duration time.Duration) Engine
}

type ContainerEnginePreset interface {
	EnginePreset
	Instantiate(label string, container string, duration time.Duration) Engine
}

type List struct {
	PodPresets       []PodEnginePreset
	ContainerPresets []ContainerEnginePreset
}
