package extensions

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"go.uber.org/zap"
	"math/rand"
	"reflect"
	"time"
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

type Extensions struct {
	Action   []ActionExtension
	Stage    []StageExtension
	Workflow []WorkflowExtension
}

func (e Extensions) Generate(r *rand.Rand, _ int) reflect.Value {
	return reflect.ValueOf(Extensions{
		Action: []ActionExtension{
			// No action extensions implemented
		},
		Stage: []StageExtension{
			UseSuspend(),
			UseStageMonitor("stage-monitor", "target-ns", time.Duration(r.Intn(60)), &zap.SugaredLogger{}),
		},
		Workflow: []WorkflowExtension{
			UseSteps(),
		},
	})
}
