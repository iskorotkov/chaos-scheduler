package workflows

import (
	"math/rand"
	"testing"
	"testing/quick"
	"time"

	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/execute"
	"go.uber.org/zap"
)

func paramsValid(params ScenarioParams) bool {
	return params.Stages.Single > 0 && params.Stages.Single <= 100 &&
		params.Stages.Similar > 0 && params.Stages.Similar <= 100 &&
		params.Stages.Mixed > 0 && params.Stages.Mixed <= 100 &&
		params.StageDuration >= time.Second
}

func TestCreateScenario(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(0))

	f := func(params ScenarioParams) bool {
		s, err := CreateScenario(params, zap.NewNop().Sugar())
		if err != nil {
			t.Log(err)
			return err == ErrScenarioParams && !paramsValid(params) ||
				err == ErrNotEnoughTargets && len(params.Targets) == 0
		}

		if !paramsValid(params) {
			t.Log("must return error when params are invalid")
			return false
		}

		if len(params.Targets) == 0 {
			t.Log("must return error when len(targets) == 0")
			return false
		}

		if len(s.Stages) == 0 {
			t.Log("scenario must contain at least one stage")
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
		wf, err := CreateWorkflow(sp, wp, zap.NewNop().Sugar())
		if err != nil {
			t.Log(err)
			return err == ErrScenarioParams && !paramsValid(sp) ||
				err == ErrNotEnoughTargets && len(sp.Targets) == 0
		}

		if !paramsValid(sp) {
			t.Log("must return error when params are invalid")
			return false
		}

		if len(sp.Targets) == 0 {
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
		executor := ep.Executor.(*execute.TestExecutor)

		wf, err := ExecuteWorkflow(sp, wp, ep, zap.NewNop().Sugar())
		if err != nil {
			t.Log(err)
			return err == ErrScenarioParams && !paramsValid(sp) ||
				err == ErrNotEnoughTargets && len(sp.Targets) == 0 ||
				err == ErrExecution && executor.Err != nil
		}

		if !paramsValid(sp) {
			t.Log("must return error when params are invalid")
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

		t.Log("succeeded")
		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: rng}); err != nil {
		t.Error(err)
	}
}
