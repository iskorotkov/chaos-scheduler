package executors

import (
	"bytes"
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"go.uber.org/zap"
	"net/http"
)

type RestExecutor struct {
	url    string
	logger *zap.SugaredLogger
}

func (r RestExecutor) Execute(wf templates.Workflow) (templates.Workflow, error) {
	req := struct {
		Workflow templates.Workflow `json:"workflow"`
	}{wf}

	marshalled, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		r.logger.Errorw(err.Error(),
			"url", r.url)
		return templates.Workflow{}, MarshalError
	}

	response, err := http.Post(r.url, "application/yaml", bytes.NewReader(marshalled))
	if err != nil {
		r.logger.Error(err.Error(),
			"url", r.url,
			"workflow", wf)
		return templates.Workflow{}, ConnectionError
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		r.logger.Warnw("executor service returned invalid status code",
			"code", response.Status)
		return templates.Workflow{}, ResponseError
	}

	return templates.Workflow{}, nil
}

//goland:noinspection GoUnusedExportedFunction
func NewRestExecutor(url string, logger *zap.SugaredLogger) RestExecutor {
	return RestExecutor{url: url, logger: logger}
}
