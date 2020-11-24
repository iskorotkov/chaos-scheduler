package scenario

import (
	"testing"
	"testing/quick"
)

var _ Failure = MockFailure(0)

type MockFailure byte

func (m MockFailure) Name() string {
	return "name"
}

func (m MockFailure) Yaml() string {
	return "yaml"
}

func TestGeneration(t *testing.T) {
	f := []Failure{MockFailure(0), MockFailure(1)}
	c := Config{f, 5, nil}

	s, err := NewScenario(c)
	if err != nil {
		t.Fatalf("scenario generation returned error: %v", err)
	}

	if len(s) != 5 {
		t.Fatalf("scenario has invalid number of stages: %v", err)
	}

	for _, stage := range s {
		if len(stage) == 0 {
			t.Fatalf("one of the stages has zero actions: %v", err)
		}
	}
}

func TestGeneration2(t *testing.T) {
	f := func(stages int) bool {
		f := []Failure{MockFailure(0), MockFailure(1)}
		c := Config{f, stages, nil}

		s, err := NewScenario(c)
		if err != nil {
			return false
		}

		if len(s) != stages {
			return false
		}

		for _, stage := range s {
			if len(stage) == 0 {
				return false
			}
		}

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

