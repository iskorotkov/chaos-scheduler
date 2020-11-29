package assemblers

import (
	"errors"
	"time"
)

var (
	WorkflowTemplateError          = errors.New("couldn't read workflow template file")
	WorkflowTemplateUnmarshalError = errors.New("couldn't unmarshall workflow template")
	WorkflowTemplatePropertyError  = errors.New("couldn't find required template property")
	StagesError                    = errors.New("number of stages must be positive")
	ActionsError                   = errors.New("number of actions in every stage must be positive")
)

type PlannedAction interface {
	Id() string
	Template() string
	Duration() time.Duration
}

type Stage []PlannedAction

type Scenario []Stage

type Workflow map[string]interface{}

type Assembler interface {
	Assemble(scenario Scenario) (Workflow, error)
}
