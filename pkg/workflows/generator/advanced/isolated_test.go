package advanced

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
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

func Test_addIsolatedFailures(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewSource(0))
	f := func(failures []failures.Failure, targets []targets.Target, params phaseParams) bool {
		if len(failures) == 0 || len(targets) == 0 {
			t.Log("zero failures or targets provided")
			return true
		}

		gen, err := NewGenerator(failures, TestTargetSeeker{targets, nil}, zap.NewNop().Sugar())
		if err != nil {
			t.Log(err)
			return false
		}

		stages := gen.addIsolatedFailures(targets, r, params)

		for _, stage := range stages {
			if stage.Duration != params.StageDuration {
				t.Log("stage duration must be equal to the specified value")
				return false
			}

			if len(stage.Actions) != 1 {
				t.Log("stage must contain exactly 1 action")
				return false
			}

			for _, action := range stage.Actions {
				if action.Name == "" ||
					action.Scale == "" ||
					action.Severity == "" {
					t.Log("action name, scale, severity must not be empty")
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
