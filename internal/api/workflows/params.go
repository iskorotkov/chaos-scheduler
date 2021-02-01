package workflows

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/k8s"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"go.uber.org/zap"
	"time"
)

type scenarioParams struct {
	server        string
	namespace     string
	label         string
	stageDuration time.Duration
	seed          int64
	stages        int
	failures      []failures.Failure
}

func createScenarioParams(params scenarioParams, logger *zap.SugaredLogger) (workflows.ScenarioParams, error) {
	finder, err := k8s.NewFinder(logger.Named("targets"))
	if err != nil {
		logger.Error(err)
		return workflows.ScenarioParams{}, internalError
	}

	return workflows.ScenarioParams{
		Seed:          params.seed,
		Stages:        params.stages,
		AppNS:         params.namespace,
		AppLabel:      params.label,
		StageDuration: params.stageDuration,
		Failures:      params.failures,
		TargetFinder:  finder,
	}, nil
}
