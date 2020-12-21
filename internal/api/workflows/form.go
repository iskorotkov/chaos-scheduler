package workflows

import (
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type workflowParams struct {
	Seed   int64 `json:"seed"`
	Stages int   `json:"stages"`
}

func parseWorkflowParams(r *http.Request, logger *zap.SugaredLogger) (workflowParams, error) {
	stages, err := strconv.ParseInt(r.FormValue("stages"), 10, 32)
	if err != nil {
		logger.Errorw(err.Error(),
			"stages", r.FormValue("stages"))
		return workflowParams{}, formParseError
	}

	seed, err := strconv.ParseInt(r.FormValue("seed"), 10, 64)
	if err != nil {
		logger.Errorw(err.Error(),
			"seed", r.FormValue("seed"))
		return workflowParams{}, formParseError
	}

	return workflowParams{Seed: seed, Stages: int(stages)}, err
}
