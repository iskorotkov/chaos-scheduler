package templates

import (
	v1 "k8s.io/api/core/v1"
)

type Container v1.Container

func NewContainerTemplate(name string, container Container) Template {
	c := v1.Container(container)
	return Template{
		Name:      name,
		Container: &c,
	}
}
