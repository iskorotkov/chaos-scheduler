package exporters

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/marshall"
)

type JsonExporter byte

func (j JsonExporter) Export(workflow Definition) (string, error) {
	req := struct {
		Workflow Definition `json:"workflow"`
	}{workflow}

	json, err := marshall.ToJson(req)
	if err != nil {
		logger.Error(err)
		return "", JsonMarshallError
	}

	return string(json), nil
}

func NewJsonExporter() Exporter {
	return JsonExporter(0)
}
