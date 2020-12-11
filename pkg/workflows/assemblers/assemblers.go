package assemblers

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generators"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
)

var (
	StagesError         = errors.New("number of stages must be positive")
	ActionsError        = errors.New("number of actions in every stage must be positive")
	ActionMarshallError = errors.New("couldn't marshall action to yaml")
)

type Assembler interface {
	Assemble(scenario generators.Scenario) (templates.Workflow, error)
}
