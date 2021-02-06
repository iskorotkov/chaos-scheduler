// Package execute handles execute of previously generated workflows.
package execute

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemble"
	"math/rand"
	"reflect"
)

var (
	ErrConnection = errors.New("couldn't post scenario to execution server")
	ErrResponse   = errors.New("execution server returned invalid status code")
)

// Executor executes a assemble.Workflow.
type Executor interface {
	// Execute passed assemble.Workflow.
	Execute(wf assemble.Workflow) (assemble.Workflow, error)
}

// TestExecutor mocks Executor.
type TestExecutor struct {
	// Workflow is a workflow to return from Execute.
	Workflow assemble.Workflow
	// Err is an error to return from Execute.
	Err error
	// SubmittedWorkflow is a workflow passed to Execute.
	SubmittedWorkflow assemble.Workflow
}

func (t TestExecutor) Generate(rand *rand.Rand, size int) reflect.Value {
	switch rand.Intn(10) {
	case 0:
		return reflect.ValueOf(TestExecutor{
			Workflow: assemble.Workflow{},
			Err:      ErrConnection,
		})
	case 1:
		return reflect.ValueOf(TestExecutor{
			Workflow: assemble.Workflow{},
			Err:      ErrResponse,
		})
	default:
		return reflect.ValueOf(TestExecutor{
			Workflow: assemble.Workflow{}.Generate(rand, size).Interface().(assemble.Workflow),
			Err:      nil,
		})
	}
}

// Execute ignores passed workflow.
func (t *TestExecutor) Execute(wf assemble.Workflow) (assemble.Workflow, error) {
	t.SubmittedWorkflow = wf
	return t.Workflow, t.Err
}
