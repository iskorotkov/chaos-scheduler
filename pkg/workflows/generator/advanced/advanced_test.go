package advanced

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
	"testing"
	"testing/quick"
	"time"
)

func TestGenerator_Generate(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewSource(0))
	f := func(failures []failures.Failure, targets []targets.Target) bool {
		seeker := TestTargetSeeker{targets, nil}

		budget := Budget{}
		modifiers := Modifiers{}
		retries := -5 + rand.Intn(20)

		var gen *Generator
		var err error
		if r.Intn(100) >= 20 {
			gen, err = NewGenerator(failures, seeker, zap.NewNop().Sugar(),
				WithBudget(budget),
				WithModifiers(modifiers),
				WithRetries(retries))
		} else {
			gen, err = NewGenerator(failures, seeker, zap.NewNop().Sugar())
		}

		if err != nil {
			t.Log(err)
			return err == generator.ZeroFailures && len(failures) == 0 ||
				err == generator.MaxFailuresError && budget.MaxFailures < 1 ||
				err == generator.MaxPointsError && budget.MaxPoints < 1 ||
				err == generator.RetriesError && retries < 1
		}

		stages := -10 + r.Intn(120)
		sc, err := gen.Generate(generator.Params{
			Stages:        stages,
			Seed:          r.Int63(),
			StageDuration: time.Duration(-10+r.Intn(200)) * time.Second,
		})
		if err != nil {
			t.Log(err)
			return err == generator.ZeroTargetsError && len(targets) == 0 ||
				err == generator.NonPositiveStagesError && stages <= 0 ||
				err == generator.TooManyStagesError && stages > 100
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
