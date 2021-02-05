package execution

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemble"
	"math/rand"
	"reflect"
)

var (
	ConnectionError = errors.New("couldn't post scenario to execution server")
	ResponseError   = errors.New("execution server returned invalid status code")
)

type Executor interface {
	Execute(wf assemble.Workflow) (assemble.Workflow, error)
}

type TestExecutor struct {
	Workflow          assemble.Workflow
	Err               error
	SubmittedWorkflow assemble.Workflow
}

func (t TestExecutor) Generate(rand *rand.Rand, size int) reflect.Value {
	switch rand.Intn(10) {
	case 0:
		return reflect.ValueOf(TestExecutor{
			Workflow: assemble.Workflow{},
			Err:      ConnectionError,
		})
	case 1:
		return reflect.ValueOf(TestExecutor{
			Workflow: assemble.Workflow{},
			Err:      ResponseError,
		})
	default:
		return reflect.ValueOf(TestExecutor{
			Workflow: assemble.Workflow{}.Generate(rand, size).Interface().(assemble.Workflow),
			Err:      nil,
		})
	}
}

func (t *TestExecutor) Execute(wf assemble.Workflow) (assemble.Workflow, error) {
	t.SubmittedWorkflow = wf
	return t.Workflow, t.Err
}
