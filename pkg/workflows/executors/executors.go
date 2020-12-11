package executors

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
)

var (
	MarshalError    = errors.New("couldn't marshall workflow to required format")
	ConnectionError = errors.New("couldn't post scenario to executor server")
	ResponseError   = errors.New("executor server returned invalid status code")
)

type Executor interface {
	Execute(wf templates.Workflow) (templates.Workflow, error)
}
