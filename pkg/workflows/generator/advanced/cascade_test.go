package advanced

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func Test_addCascadeFailures(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewSource(0))
	f := func(fs []failures.Failure, ts []targets.Target, params phaseParams) bool {
		if len(fs) == 0 || len(ts) == 0 {
			t.Log("zero failures or targets provided")
			return true
		}

		gen, err := NewGenerator(fs, TestTargetSeeker{ts, nil}, zap.NewNop().Sugar())
		if err != nil {
			t.Log(err)
			return false
		}

		stages := gen.addCascadeFailures(ts, r, params)

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