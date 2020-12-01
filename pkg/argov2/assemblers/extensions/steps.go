package extensions

import "github.com/iskorotkov/chaos-scheduler/pkg/argov2/templates"

type Steps struct{}

func (s Steps) Apply(ids [][]string) Extension {
	return templates.NewStepsTemplate(ids)
}
