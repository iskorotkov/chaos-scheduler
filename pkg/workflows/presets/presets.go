package presets

type Info struct {
	Lethal bool
}

type EnginePreset interface {
	Type() ExperimentType
	Info() Info
}

type PodEnginePreset interface {
	EnginePreset
	Instantiate(label string) Engine
}

type ContainerEnginePreset interface {
	EnginePreset
	Instantiate(label string, container string) Engine
}

type List struct {
	PodPresets       []PodEnginePreset
	ContainerPresets []ContainerEnginePreset
}
