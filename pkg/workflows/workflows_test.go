package workflows

import (
	"go.uber.org/zap"
	"math/rand"
	"testing"
	"testing/quick"
)

func TestCreateScenario(t *testing.T) {
	rng := rand.New(rand.NewSource(0))

	f := func(params ScenarioParams) bool {
		s, err := CreateScenario(params, zap.NewNop().Sugar())
		if err == ScenarioParamsError &&
			(params.Stages <= 0 || params.Stages > 100 || params.StageDuration <= 0) {
			return true
		} else if err != nil {
			return false
		}

		if len(s.Stages) == 0 {
			return false
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: rng}); err != nil {
		t.Error(err)
	}
}

func TestCreateWorkflow(t *testing.T) {
	rng := rand.New(rand.NewSource(0))

	f := func(sp ScenarioParams, wp WorkflowParams) bool {
		wf, err := CreateWorkflow(sp, wp, zap.NewNop().Sugar())
		if err == ScenarioParamsError &&
			(sp.Stages <= 0 || sp.Stages > 100 || sp.StageDuration <= 0) {
			return true
		} else if err != nil {
			return false
		}

		if wf.Namespace == "" ||
			wf.GenerateName == "" {
			return false
		}

		if len(wf.Labels) != 0 ||
			len(wf.Annotations) != 0 {
			return false
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: rng}); err != nil {
		t.Error(err)
	}
}
