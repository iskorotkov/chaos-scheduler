package executors

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"net/http"
	"strings"
)

type RestExecutor struct {
	Url string
}

func (r RestExecutor) Execute(workflow string) error {
	response, err := http.Post(r.Url, "application/yaml", strings.NewReader(workflow))
	if err != nil {
		logger.Error(err)
		return ConnectionError
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		logger.Warning(fmt.Sprintf("executor service returned status '%s'", response.Status))
		return ResponseError
	}

	return nil
}

func NewRestExecutor(url string) RestExecutor {
	return RestExecutor{Url: url}
}
