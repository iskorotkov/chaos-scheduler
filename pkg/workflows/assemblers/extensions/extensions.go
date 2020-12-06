package extensions

import "github.com/iskorotkov/chaos-scheduler/pkg/workflows/generators"

type Extension interface {
	Id() string
}

type ActionExtension interface {
	Apply(action generators.Action, stageIndex, actionIndex int) Extension
}

type StageExtension interface {
	Apply(stage generators.Stage, stageIndex int) Extension
}

type WorkflowExtension interface {
	Apply(ids [][]string) Extension
}

type List struct {
	ActionExtensions   []ActionExtension
	StageExtensions    []StageExtension
	WorkflowExtensions []WorkflowExtension
}
