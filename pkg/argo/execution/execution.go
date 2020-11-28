package execution

import (
	"errors"
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/output"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"net/http"
	"strings"
)

var (
	ConnectionError = errors.New("couldn't post scenario to executor server")
	ResponseError   = errors.New("executor server returned invalid status code")
	FormatError     = errors.New("couldn't format provided scenario")
)

func ExecuteFromConfig(url string, config output.Config) error {
	workflow, err := output.GenerateFromConfig(config)
	if err != nil {
		logger.Error(err)
		return FormatError
	}

	return ExecuteWorkflow(url, workflow)
}

func ExecuteWorkflow(url string, w string) error {
	r, err := http.Post(url, "application/yaml", strings.NewReader(w))
	if err != nil {
		logger.Error(err)
		return ConnectionError
	}

	if r.StatusCode != http.StatusOK && r.StatusCode != http.StatusCreated {
		logger.Warning(fmt.Sprintf("executor service returned status '%s'", r.Status))
		return ResponseError
	}

	return nil
}
