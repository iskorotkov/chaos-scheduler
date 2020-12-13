package scenarios

import (
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type form struct {
	Seed   int64
	Stages int
}

func parseScenarioParams(r *http.Request, logger *zap.SugaredLogger) (form, error) {
	stages, err := strconv.ParseInt(r.FormValue("stages"), 10, 32)
	if err != nil {
		logger.Errorw(err.Error(),
			"stages", r.FormValue("stages"))
		return form{}, formParseError
	}

	seed, err := strconv.ParseInt(r.FormValue("seed"), 10, 64)
	if err != nil {
		logger.Errorw(err.Error(),
			"seed", r.FormValue("seed"))
		return form{}, formParseError
	}

	return form{Seed: seed, Stages: int(stages)}, err
}
