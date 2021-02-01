package workflows

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
	"testing"
	"testing/quick"
)

func TestCreateScenario(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(0))

	f := func(params ScenarioParams) bool {
		s, err := CreateScenario(params, zap.NewNop().Sugar())
		if (err == ScenarioParamsError) ==
			(params.Stages <= 0 || params.Stages > 100 || params.StageDuration <= 0) {
			t.Log(err)
			return true
		} else if (err == NotEnoughTargetsError) ==
			(len(params.TargetFinder.(targets.TestTargetFinder).Targets) == 0) {
			t.Log(err)
			return true
		} else if err != nil {
			t.Log(err)
			return false
		}

		if len(s.Stages) == 0 {
			t.Log("scenario must contain at least one stage")
			return false
		}

		targetFinder := params.TargetFinder.(targets.TestTargetFinder)
		if targetFinder.SubmittedLabel != params.AppLabel ||
			targetFinder.SubmittedNamespace != params.AppNS {
			t.Log("app label and namespace must equal to values from params")
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

	rng := rand.New(rand.NewSource(0))

	f := func(sp ScenarioParams, wp WorkflowParams) bool {
		wf, err := CreateWorkflow(sp, wp, zap.NewNop().Sugar())
		if (err == NotEnoughTargetsError) ==
			(len(sp.TargetFinder.(targets.TestTargetFinder).Targets) == 0) {
			t.Log(err)
			return true
		} else if err != nil {
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

		targetFinder := sp.TargetFinder.(targets.TestTargetFinder)
		if targetFinder.SubmittedLabel != sp.AppLabel ||
			targetFinder.SubmittedNamespace != sp.AppNS {
			t.Log("app label and namespace must equal to values from params")
			return false
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: rng}); err != nil {
		t.Error(err)
	}
}

func TestExecuteWorkflow(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(0))

	f := func(sp ScenarioParams, wp WorkflowParams, ep ExecutionParams) bool {
		wf, err := ExecuteWorkflow(sp, wp, ep, zap.NewNop().Sugar())
		if (err == NotEnoughTargetsError) ==
			(len(sp.TargetFinder.(targets.TestTargetFinder).Targets) == 0) {
			t.Log(err)
			return true
		} else if err != nil {
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

		targetFinder := sp.TargetFinder.(targets.TestTargetFinder)
		if targetFinder.SubmittedLabel != sp.AppLabel ||
			targetFinder.SubmittedNamespace != sp.AppNS {
			t.Log("app label and namespace must equal to values from params")
			return false
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: rng}); err != nil {
		t.Error(err)
	}
}
