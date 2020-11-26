package communication

import (
	"fmt"
	"net/http"
	"strings"
)

type Executor struct {
	host string
	port int
}

func (e Executor) Execute(config FormatConfig) error {
	formatted, err := GenerateWorkflow(config)
	if err != nil {
		return fmt.Errorf("couldn't format provided scenario: %v", err)
	}

	url := fmt.Sprintf("%s:%d", e.host, e.port)
	r, err := http.Post(url, "text/yaml", strings.NewReader(formatted))
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
