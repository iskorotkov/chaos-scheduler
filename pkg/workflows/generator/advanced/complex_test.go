package advanced

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
	"testing"
	"testing/quick"
)

func Test_addComplexFailures(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewSource(0))
	f := func(fs []failures.Failure, targets []targets.Target, params phaseParams) bool {
		if len(fs) == 0 || len(targets) == 0 {
			t.Log("zero failures or targets provided")
			return true
		}

		gen, err := NewGenerator(fs, TestTargetSeeker{targets, nil}, zap.NewNop().Sugar())
		if err != nil {
			t.Log(err)
			return false
		}

		stages := gen.addComplexFailures(targets, r, params)

		for _, stage := range stages {
			if stage.Duration != params.StageDuration {
				t.Log("stage duration must be equal to the specified value")
				return false
			}

			if len(stage.Actions) == 0 {
				t.Log("stage must contain at least 1 action")
				return false
			}

			if len(stage.Actions) > gen.budget.MaxFailures {
				t.Log("total number os actions per stage must be less or equal to budget's max value")
				return false
			}

			points := Cost(0)
			for _, action := range stage.Actions {
				if action.Name == "" ||
					action.Scale == "" ||
					action.Severity == "" {
					t.Log("action name, scale, severity must not be empty")
					return false
				}

				points += calculateCost(gen.modifiers, failures.Failure{
					Template: nil,
					Scale:    action.Scale,
					Severity: action.Severity,
				})
			}

			if points > gen.budget.MaxPoints {
				t.Log("total points of each stage must be less or equal to budget points")
				return false
			}
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: r}); err != nil {
		t.Error(err)
	}
}
