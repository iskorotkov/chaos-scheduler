package execution

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/output"
	"net/http"
	"strings"
)

func ExecuteFromConfig(url string, config output.Config) error {
	workflow, err := output.GenerateFromConfig(config)
	if err != nil {
		return fmt.Errorf("couldn't format provided scenario: %v", err)
	}

	return ExecuteWorkflow(url, workflow)
}

func ExecuteWorkflow(url string, w string) error {
	r, err := http.Post(url, "text/yaml", strings.NewReader(w))
	if err != nil {
		return fmt.Errorf("couldn't post scenario to executor server: %v", err)
	}

	if r.StatusCode != http.StatusOK && r.StatusCode != http.StatusCreated {
		return fmt.Errorf("executor server returned invalid status code: %v", r.StatusCode)
	}

	return nil
}
