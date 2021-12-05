package generate

import (
	"math/rand"
	"testing"
	"testing/quick"

	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
)

func Test_addComplexFailures(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewSource(0))
	f := func(params Params) bool {
		if len(params.Failures) == 0 || len(params.Targets) == 0 {
			t.Log("zero failures or targets provided")
			return true
		}

		stages := addComplexFailures(params, rand.New(rand.NewSource(0)), r)

		for _, stage := range stages {
			if stage.Duration != params.StageDuration {
				t.Log("stage duration must be equal to the specified value")
				return false
			}

			if len(stage.Actions) == 0 {
				t.Log("stage must contain at least 1 action")
				return false
			}

			if len(stage.Actions) > params.Budget.MaxFailures {
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

				points += calculateCost(params.Modifiers, failures.Failure{
					Blueprint: nil,
					Scale:     action.Scale,
					Severity:  action.Severity,
				})
			}

			if points > params.Budget.MaxPoints {
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
