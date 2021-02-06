package generate

import (
	"math/rand"
	"testing"
	"testing/quick"
	"time"
)

func TestGenerator_Generate(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewSource(0))
	f := func(params Params) bool {
		sc, err := Generate(params)
		if err != nil {
			t.Log(err)
			return err == ErrZeroTargets && len(params.Targets) == 0 ||
				err == ErrNonPositiveStages && params.Stages <= 0 ||
				err == ErrTooManyStages && params.Stages > 100 ||
				err == ErrZeroFailures && len(params.Failures) == 0 ||
				err == ErrMaxFailures && params.Budget.MaxFailures < 1 ||
				err == ErrMaxPoints && params.Budget.MaxPoints < 1 ||
				err == ErrStageDuration && params.StageDuration < time.Second
		}

		if len(sc.Stages) == 0 {
			t.Log("scenario must contain at least one stage")
			return false
		}

		for _, stage := range sc.Stages {
			if stage.Duration.Seconds() <= 0 {
				t.Log("stage duration must be >= 1s")
				return false
			}

			if len(stage.Actions) == 0 {
				t.Log("stage must contain at least one action")
				return false
			}

			for _, action := range stage.Actions {
				if action.Name == "" ||
					action.Scale == "" ||
					action.Severity == "" {
					t.Log("action name, scale and severity must be non-empty")
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
