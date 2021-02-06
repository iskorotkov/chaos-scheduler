package handlers

import (
	"go.uber.org/zap"
	"math"
	"math/rand"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"testing/quick"
)

func Test_parseForm(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewSource(0))
	f := func() bool {
		request := httptest.NewRequest("GET", "/url", nil)

		seed := rand.Int63() - math.MaxInt64/2
		stages := -10 + rand.Intn(120)

		variant := rand.Intn(10)
		invalidVariants := 3
		switch variant {
		case 0:
			request.PostForm = url.Values{}
		case 1:
			request.PostForm = url.Values{
				"stages": {strconv.FormatInt(int64(stages), 10)},
			}
		case 2:
			request.PostForm = url.Values{
				"seed": {strconv.FormatInt(seed, 10)},
			}
		default:
			request.PostForm = url.Values{
				"seed":   {strconv.FormatInt(seed, 10)},
				"stages": {strconv.FormatInt(int64(stages), 10)},
			}
		}

		form, ok := parseForm(request, zap.NewNop().Sugar())
		if !ok {
			if variant < invalidVariants {
				t.Log("not all values were provided")
				return true
			} else {
				t.Log("couldn't parse form")
				return false
			}
		}

		if variant < invalidVariants {
			t.Log("must return error on invalid form")
			return false
		}

		if form.Stages != stages || form.Seed != seed {
			t.Log("stages and seed must be equal to provided values")
			return false
		}

		t.Log("succeeded")
		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: r}); err != nil {
		t.Error(err)
	}
}
