package assemblers

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenarios"
	"time"
)

var (
	WorkflowTemplateError          = errors.New("couldn't read workflow template file")
	WorkflowTemplateUnmarshalError = errors.New("couldn't unmarshall workflow template")
	WorkflowTemplatePropertyError  = errors.New("couldn't find required template property")
	StagesError                    = errors.New("number of stages must be positive")
	ActionsError                   = errors.New("number of actions in every stage must be positive")
	TemplateParseError             = errors.New("couldn't parse template text")
	TemplateExecuteError           = errors.New("couldn't execute template text")
)

type Workflow map[string]interface{}

type Assembler interface {
	Assemble(scenario scenarios.Scenario) (Workflow, error)
}

type context struct {
	Name     string
	Duration time.Duration
	Stage    int
	Index    int
}
