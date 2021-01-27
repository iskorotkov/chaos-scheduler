package assemble

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generate"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
)

type suspend struct{}

func (d suspend) Apply(stage generate.Stage, index int) []templates.Template {
	id := fmt.Sprintf("delay-%d", index+1)
	return []templates.Template{
		templates.NewSuspendTemplate(id, stage.Duration),
	}
}

//goland:noinspection GoUnusedExportedFunction
func UseSuspend() StageExtension {
	return suspend{}
}
