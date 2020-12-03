package templates

import "github.com/iskorotkov/chaos-scheduler/pkg/workflows/engines"

type Container struct {
	Name    string           `yaml:"name" json:"name"`
	Image   string           `yaml:"image" json:"image"`
	Env     []engines.EnvVar `yaml:"env" json:"env"`
	Ports   []string         `yaml:"ports" json:"ports"`
	Command []string         `yaml:"command" json:"command"`
	Args    []string         `yaml:"args" json:"args"`
}

type ContainerTemplate struct {
	Name      string    `yaml:"name" json:"name"`
	Container Container `yaml:"container" json:"container"`
}

func (c ContainerTemplate) Id() string {
	return c.Name
}

func NewContainerTemplate(name string, container Container) ContainerTemplate {
	return ContainerTemplate{name, container}
}
