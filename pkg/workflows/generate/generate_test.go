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
			return err == ZeroTargetsError && len(params.Targets) == 0 ||
				err == NonPositiveStagesError && params.Stages <= 0 ||
				err == TooManyStagesError && params.Stages > 100 ||
				err == ZeroFailures && len(params.Failures) == 0 ||
				err == MaxFailuresError && params.Budget.MaxFailures < 1 ||
				err == MaxPointsError && params.Budget.MaxPoints < 1 ||
				err == RetriesError && params.Retries < 1 ||
				err == StageDurationError && params.StageDuration < time.Second
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
