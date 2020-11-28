package marshall

import (
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
)

func ToJson(in interface{}) ([]byte, error) {
	b, err := json.Marshal(in)
	if err != nil {
		logger.Error(err)
		return nil, MarshallError
	}

	return b, nil
}

func FromJson(data []byte) (Tree, error) {
	var t Tree
	err := json.Unmarshal(data, &t)
	if err != nil {
		logger.Error(err)
		return nil, UnmarshallError
	}

	return t, nil
}
