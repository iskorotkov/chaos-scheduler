package extensions

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
)

type ActionExtension interface {
	Apply(action generator.Action, stageIndex, actionIndex int) []templates.Template
}

type StageExtension interface {
	Apply(stage generator.Stage, stageIndex int) []templates.Template
}

type WorkflowExtension interface {
	Apply(ids [][]string) []templates.Template
}

type List struct {
	ActionExtensions   []ActionExtension
	StageExtensions    []StageExtension
	WorkflowExtensions []WorkflowExtension
}
