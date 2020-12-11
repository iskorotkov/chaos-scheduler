package executors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"net/http"
)

type RestExecutor struct {
	Url string
}

func (r RestExecutor) Execute(wf templates.Workflow) (templates.Workflow, error) {
	req := struct {
		Workflow templates.Workflow `json:"workflow"`
	}{wf}

	marshalled, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		logger.Error(err)
		return templates.Workflow{}, MarshalError
	}

	response, err := http.Post(r.Url, "application/yaml", bytes.NewReader(marshalled))
	if err != nil {
		logger.Error(err)
		return templates.Workflow{}, ConnectionError
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		logger.Warning(fmt.Sprintf("executor service returned status '%s'", response.Status))
		return templates.Workflow{}, ResponseError
	}

	return templates.Workflow{}, nil
}

func NewRestExecutor(url string) RestExecutor {
	return RestExecutor{Url: url}
}
