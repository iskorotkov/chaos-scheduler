package executors

import "errors"

var (
	ConnectionError = errors.New("couldn't post scenario to executor server")
	ResponseError   = errors.New("executor server returned invalid status code")
)

type Executor interface {
	Execute(workflow string) error
}
