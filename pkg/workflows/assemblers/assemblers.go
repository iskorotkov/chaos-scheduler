package assemblers

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/scenarios"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/workflow"
)

var (
	StagesError         = errors.New("number of stages must be positive")
	ActionsError        = errors.New("number of actions in every stage must be positive")
	ActionMarshallError = errors.New("couldn't marshall action to yaml")
)

type Assembler interface {
	Assemble(scenario scenarios.Scenario) (workflow.Workflow, error)
}
