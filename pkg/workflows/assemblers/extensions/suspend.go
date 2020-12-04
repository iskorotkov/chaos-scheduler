package extensions

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generators"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
)

type Suspend struct{}

func (d Suspend) Apply(stage generators.Stage, index int) Extension {
	id := fmt.Sprintf("delay-%d", index+1)
	return templates.NewSuspendTemplate(id, stage.Duration)
}

func UseSuspend() StageExtension {
	return Suspend{}
}
