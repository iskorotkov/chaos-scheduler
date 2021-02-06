package generate

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func Test_addCascadeFailures(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewSource(0))
	f := func(params Params) bool {
		if len(params.Failures) == 0 || len(params.Targets) == 0 {
			t.Log("zero failures or targets provided")
			return true
		}

		stages := addCascadeFailures(params, rand.New(rand.NewSource(0)))

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

			firstAction := stage.Actions[0]

			for _, action := range stage.Actions {
				if action.Name != firstAction.Name ||
					action.Scale != firstAction.Scale ||
					action.Severity != firstAction.Severity ||
					!reflect.DeepEqual(action.Engine.Metadata, firstAction.Engine.Metadata) {
					t.Log("action name, scale, severity and engine of all actions in every stage must be equal")
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
