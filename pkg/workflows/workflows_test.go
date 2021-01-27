package workflows

import (
	"go.uber.org/zap"
	"math/rand"
	"testing"
	"testing/quick"
)

func TestCreateScenario(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("test requires connection to a running Kubernetes cluster")
	}

	rng := rand.New(rand.NewSource(0))

	f := func(params ScenarioParams) bool {
		s, err := CreateScenario(params, zap.NewNop().Sugar())
		if err == ScenarioParamsError &&
			(params.Stages <= 0 || params.Stages > 100 || params.StageDuration <= 0) {
			t.Log("invalid scenario params")
			return true
		} else if err == TargetsSeekerError {
			t.Skip("can't create target seeker in this environment; probably Kubernetes cluster isn't running")
			return true
		} else if err != nil {
			t.Log(err)
			return false
		}

		if len(s.Stages) == 0 {
			t.Log("scenario must contain at least one stage")
			return false
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: rng}); err != nil {
		t.Error(err)
	}
}

func TestCreateWorkflow(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("test requires connection to a running Kubernetes cluster")
	}

	rng := rand.New(rand.NewSource(0))

	f := func(sp ScenarioParams, wp WorkflowParams) bool {
		wf, err := CreateWorkflow(sp, wp, zap.NewNop().Sugar())
		if err != nil {
			if err == TargetsSeekerError {
				t.Skip("can't create target seeker in this environment; probably Kubernetes cluster isn't running")
			}

			t.Log(err)
			return err == ScenarioParamsError && (sp.Stages <= 0 || sp.Stages > 100 || sp.StageDuration <= 0)
		}

		if wf.Namespace == "" ||
			wf.GenerateName == "" {
			t.Log("namespace and generateName must not be empty")
			return false
		}

		if len(wf.Labels) != 0 ||
			len(wf.Annotations) != 0 {
			t.Log("labels and annotations must contain zero items")
			return false
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: rng}); err != nil {
		t.Error(err)
	}
}
