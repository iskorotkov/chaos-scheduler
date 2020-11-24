package argo

import (
	"bytes"
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenario"
	"gopkg.in/yaml.v2"
	"net/http"
)

type Executor struct {
	host string
	port int
}

func (e Executor) Execute(scenario scenario.Scenario) error {
	url := fmt.Sprintf("%s:%d", e.host, e.port)
	content, err := yaml.Marshal(scenario)
	if err != nil {
		return fmt.Errorf("couldn't marshal scenario to yaml: %v", err)
	}

	r, err := http.Post(url, "text/yaml", bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("couldn't post scenario to executor server: %v", err)
	}

	if r.StatusCode != http.StatusOK && r.StatusCode != http.StatusCreated {
		return fmt.Errorf("executor server returned invalid status code: %v", r.StatusCode)
	}

	return nil
}

func NewExecutor(host string, port int) Executor {
	return Executor{host, port}
}
