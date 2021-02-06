package handlers

import (
	"context"
	"encoding/json"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
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

func Test_preview(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewSource(0))
	f := func(cfg config.Config, finder targets.TestTargetFinder) bool {
		request := httptest.NewRequest("GET", "/url", nil)

		seed := rand.Int63() - math.MaxInt64/2
		stages := -10 + rand.Intn(120)
		request.PostForm = url.Values{
			"seed":   {strconv.FormatInt(seed, 10)},
			"stages": {strconv.FormatInt(int64(stages), 10)},
		}

		ctx := request.Context()
		ctx = context.WithValue(ctx, "config", &cfg)
		ctx = context.WithValue(ctx, "finder", &finder)
		request = request.WithContext(ctx)

		recorder := httptest.NewRecorder()

		preview(recorder, request, zap.NewNop().Sugar())
		if recorder.Code != 200 {
			t.Logf("%d: %s", recorder.Code, recorder.Body)
			if recorder.Code == http.StatusBadRequest {
				return stages < 1 || stages > 100 || cfg.StageDuration < time.Second
			} else if recorder.Code == http.StatusInternalServerError {
				return len(finder.Targets) == 0 || finder.Err != nil
			}
		}

		if finder.Err != nil {
			t.Log("handler must return error when finder returns error")
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

		var response previewResponse
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Log(err)
			return false
		}

		if len(response.Scenario.Stages) != 3*stages {
			t.Log("number of stages must match")
			return false
		}

		for _, stage := range response.Scenario.Stages {
			if stage.Duration != cfg.StageDuration {
				t.Log("stage duration must match")
				return false
			}

			for _, action := range stage.Actions {
				if action.Engine.Metadata.Namespace != cfg.ChaosNS {
					t.Log("actions namespace must match chaos namespace")
					return false
				}
			}
		}

		t.Log("succeeded")
		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: r}); err != nil {
		t.Error(err)
	}
}
