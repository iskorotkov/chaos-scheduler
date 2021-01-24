package roundrobin

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
	"testing"
	"testing/quick"
)

type TestTargetSeeker struct {
	targets []targets.Target
	error   error
}

func (t TestTargetSeeker) Targets() ([]targets.Target, error) {
	return t.targets, t.error
}

func TestRoundRobin_Generate(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	f := func(failures []failures.Failure, targets []targets.Target, params generator.Params) bool {
		gen := NewRoundRobin(failures, TestTargetSeeker{targets, nil}, zap.NewNop().Sugar())

		sc, err := gen.Generate(params)
		if err == generator.TargetsError && len(targets) == 0 {
			return true
		} else if err == generator.NonPositiveStagesError && params.Stages <= 0 {
			return true
		} else if err == generator.ZeroFailures && len(failures) == 0 {
			return true
		} else if err == generator.TooManyStagesError && params.Stages > 100 {
			return true
		} else if err != nil {
			t.Log(err)
			return false
		}

		if len(sc.Stages) != params.Stages {
			t.Log("scenario must contain specified number of stages")
			return false
		}

		for _, stage := range sc.Stages {
			if stage.Duration != params.StageDuration {
				t.Log("duration time must equal to specified number")
				return false
			}

			for _, action := range stage.Actions {
				if action.Name == "" ||
					action.Severity == "" ||
					action.Scale == "" {
					t.Log("action name, severity and scale must not be empty")
					return false
				}
			}
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: r}); err != nil {
		t.Error(err)
	}
}
