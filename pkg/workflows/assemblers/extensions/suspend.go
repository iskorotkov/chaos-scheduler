package extensions

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
)

type Suspend struct{}

func (d Suspend) Apply(stage generator.Stage, index int) []templates.Template {
	id := fmt.Sprintf("delay-%d", index+1)
	return []templates.Template{
		templates.NewSuspendTemplate(id, stage.Duration),
	}
}

//goland:noinspection GoUnusedExportedFunction
func UseSuspend() StageExtension {
	return Suspend{}
}
