package scenarios

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"net/http"
	"strconv"
)

type Form struct {
	Seed   int64
	Stages int
}

func parseScenarioParams(r *http.Request) (Form, error) {
	stages, err := strconv.ParseInt(r.FormValue("stages"), 10, 32)
	if err != nil {
		logger.Error(err)
		return Form{}, FormParseError
	}

	seed, err := strconv.ParseInt(r.FormValue("seed"), 10, 64)
	if err != nil {
		logger.Error(err)
		return Form{}, FormParseError
	}

	return Form{Seed: seed, Stages: int(stages)}, err
}
