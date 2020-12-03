package extensions

import "github.com/iskorotkov/chaos-scheduler/pkg/workflows/scenarios"

type Extension interface {
	Id() string
}

type ActionExtension interface {
	Apply(action scenarios.PlannedAction, stageIndex, actionIndex int) Extension
}

type StageExtension interface {
	Apply(stage scenarios.Stage, stageIndex int) Extension
}

type WorkflowExtension interface {
	Apply(ids [][]string) Extension
}
