package exporters

import "errors"

var (
	JsonMarshallError = errors.New("couldn't marshall workflow to json")
)

type Definition interface{}

type Exporter interface {
	Export(workflow Definition) (string, error)
}
