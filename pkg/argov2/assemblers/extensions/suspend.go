package extensions

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenarios"
)

type Suspend struct{}

func (d Suspend) Apply(stage scenarios.Stage, index int) Extension {
	id := fmt.Sprintf("delay-%d", index+1)
	return templates.NewSuspendTemplate(id, stage.Duration())
}

func UseSuspend() StageExtension {
	return Suspend{}
}
