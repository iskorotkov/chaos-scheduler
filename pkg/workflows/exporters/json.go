package exporters

import (
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
)

type JsonExporter byte

func (j JsonExporter) Export(workflow Definition) (string, error) {
	req := struct {
		Workflow Definition `json:"workflow"`
	}{workflow}

	marshalled, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		logger.Error(err)
		return "", JsonMarshallError
	}

	return string(marshalled), nil
}

func NewJsonExporter() Exporter {
	return JsonExporter(0)
}
