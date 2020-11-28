package marshall

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"gopkg.in/yaml.v2"
)

type Tree map[string]interface{}

func ToYaml(in interface{}) ([]byte, error) {
	b, err := yaml.Marshal(in)
	if err != nil {
		logger.Error(err)
		return nil, MarshallError
	}

	return b, nil
}

func FromYaml(b []byte) (Tree, error) {
	var m Tree
	err := yaml.Unmarshal(b, &m)
	if err != nil {
		logger.Error(err)
		return nil, UnmarshallError
	}

	for k, v := range m {
		m[k], err = unmarshall(v)
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func unmarshall(elem interface{}) (interface{}, error) {
	m, ok := elem.(map[interface{}]interface{})
	if !ok {
		return elem, nil
	}

	res := make(Tree)
	for k, v := range m {
		str, ok := k.(string)
		if !ok {
			return nil, PropertyConversionError
		}

		var err error
		res[str], err = unmarshall(v)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
