package extensions

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generators"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
)

type ActionExtension interface {
	Apply(action generators.Action, stageIndex, actionIndex int) []templates.Template
}

type StageExtension interface {
	Apply(stage generators.Stage, stageIndex int) []templates.Template
}

type WorkflowExtension interface {
	Apply(ids [][]string) []templates.Template
}

type List struct {
	ActionExtensions   []ActionExtension
	StageExtensions    []StageExtension
	WorkflowExtensions []WorkflowExtension
}
