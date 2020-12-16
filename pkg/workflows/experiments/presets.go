package experiments

import "time"

type Info struct {
	Lethal bool
}

type EnginePreset interface {
	Type() ExperimentType
	Info() Info
}

type PodPreset interface {
	EnginePreset
	Instantiate(label string, duration time.Duration) Engine
}

type ContainerPreset interface {
	EnginePreset
	Instantiate(label string, container string, duration time.Duration) Engine
}

type NodePreset interface {
	EnginePreset
	Instantiate(label string, node string, duration time.Duration) Engine
}

type PresetsList struct {
	ContainerPresets []ContainerPreset
	PodPresets       []PodPreset
	NodePreset       []NodePreset
}
