package workflows

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/execution"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
	"testing"
	"testing/quick"
	"time"
)

func paramsValid(params ScenarioParams) bool {
	return params.Stages > 0 && params.Stages <= 100 && params.StageDuration >= time.Second
}

func TestCreateScenario(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(0))

	f := func(params ScenarioParams) bool {
		targetFinder := params.TargetFinder.(*targets.TestTargetFinder)

		s, err := CreateScenario(params, zap.NewNop().Sugar())
		if err != nil {
			t.Log(err)
			return err == ScenarioParamsError && !paramsValid(params) ||
				err == NotEnoughTargetsError && len(targetFinder.Targets) == 0 ||
				err == TargetsFetchError && targetFinder.Err != nil
		}

		if targetFinder.Err != nil {
			t.Log("must return error when target finder returns error")
			return false
		}

		if !paramsValid(params) {
			t.Log("must return error when params are invalid")
			return false
		}

		if len(targetFinder.Targets) == 0 {
			t.Log("must return error when len(targets) == 0")
			return false
		}

		if len(s.Stages) == 0 {
			t.Log("scenario must contain at least one stage")
			return false
		}

		if targetFinder.SubmittedLabel != params.AppLabel ||
			targetFinder.SubmittedNamespace != params.AppNS {
			t.Log("app label and namespace must equal to values from params")
			return false
		}

		t.Log("succeeded")
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
		targetFinder := sp.TargetFinder.(*targets.TestTargetFinder)

		wf, err := CreateWorkflow(sp, wp, zap.NewNop().Sugar())
		if err != nil {
			t.Log(err)
			return err == ScenarioParamsError && !paramsValid(sp) ||
				err == NotEnoughTargetsError && len(targetFinder.Targets) == 0 ||
				err == TargetsFetchError && targetFinder.Err != nil
		}

		if targetFinder.Err != nil {
			t.Log("must return error when target finder returns error")
			return false
		}

		if !paramsValid(sp) {
			t.Log("must return error when params are invalid")
			return false
		}

		if len(targetFinder.Targets) == 0 {
			t.Log("must return error when len(targets) == 0")
			return false
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

		if targetFinder.SubmittedLabel != sp.AppLabel ||
			targetFinder.SubmittedNamespace != sp.AppNS {
			t.Log("app label and namespace must equal to values from params")
			return false
		}

		t.Log("succeeded")
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
		targetFinder := sp.TargetFinder.(*targets.TestTargetFinder)
		executor := ep.Executor.(*execution.TestExecutor)

		wf, err := ExecuteWorkflow(sp, wp, ep, zap.NewNop().Sugar())
		if err != nil {
			t.Log(err)
			return err == ScenarioParamsError && !paramsValid(sp) ||
				err == NotEnoughTargetsError && len(targetFinder.Targets) == 0 ||
				err == TargetsFetchError && targetFinder.Err != nil ||
				err == ExecutionError && executor.Err != nil
		}

		if targetFinder.Err != nil {
			t.Log("must return error when target finder returns error")
			return false
		}

		if !paramsValid(sp) {
			t.Log("must return error when params are invalid")
			return false
		}

		if len(targetFinder.Targets) == 0 {
			t.Log("must return error when len(targets) == 0")
			return false
		}

		if wf.Namespace == "" ||
			wf.GenerateName == "" {
			t.Log("namespace and generateName must not be empty")
			return false
		}

		if len(executor.SubmittedWorkflow.Labels) != 0 ||
			len(executor.SubmittedWorkflow.Annotations) != 0 {
			t.Log("labels and annotations must contain zero items")
			return false
		}

		if targetFinder.SubmittedLabel != sp.AppLabel ||
			targetFinder.SubmittedNamespace != sp.AppNS {
			t.Log("app label and namespace must equal to values from params")
			return false
		}

		t.Log("succeeded")
		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: rng}); err != nil {
		t.Error(err)
	}
}
