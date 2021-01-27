package generate

import (
	"math/rand"
	"testing"
	"testing/quick"
)

func Test_addIsolatedFailures(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewSource(0))
	f := func(params Params) bool {
		if len(params.Failures) == 0 || len(params.Targets) == 0 {
			t.Log("zero failures or targets provided")
			return true
		}

		stages := addIsolatedFailures(params)

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
