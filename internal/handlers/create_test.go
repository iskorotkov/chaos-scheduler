package handlers

import (
	"context"
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/execute"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"testing/quick"
	"time"
)

func Test_create(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewSource(0))
	f := func(cfg config.Config, finder targets.TestTargetFinder, executor execute.TestExecutor) bool {
		request := httptest.NewRequest("GET", "/url", nil)

		seed := rand.Int63() - math.MaxInt64/2
		stages := int64(-10 + rand.Intn(120))
		request.PostForm = url.Values{
			"seed":   {strconv.FormatInt(seed, 10)},
			"stages": {strconv.FormatInt(stages, 10)},
		}

		ctx := request.Context()
		ctx = context.WithValue(ctx, "config", &cfg)
		ctx = context.WithValue(ctx, "finder", &finder)
		ctx = context.WithValue(ctx, "executor", &executor)
		request = request.WithContext(ctx)

		recorder := httptest.NewRecorder()

		create(recorder, request, zap.NewNop().Sugar())
		if recorder.Code != http.StatusOK {
			t.Logf("%d: %s", recorder.Code, recorder.Body)
			if recorder.Code == http.StatusBadRequest {
				return stages < 1 || stages > 100 || cfg.StageDuration < time.Second
			} else if recorder.Code == http.StatusInternalServerError {
				return len(finder.Targets) == 0 || finder.Err != nil || executor.Err != nil
			}
		}

		if finder.Err != nil || executor.Err != nil {
			t.Log("handler must return error when finder or executor return error")
			return false
		}

		if len(finder.Targets) == 0 {
			t.Log("handler must return error on zero targets")
			return false
		}

		if stages < 1 || stages > 100 {
			t.Log("handler must return error on incorrect params")
			return false
		}

		if cfg.StageDuration < time.Second {
			t.Log("handler must return error on incorrect stage duration")
			return false
		}

		var response createResponse
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Log(err)
			return false
		}

		t.Log("succeeded")
		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: r}); err != nil {
		t.Error(err)
	}
}
