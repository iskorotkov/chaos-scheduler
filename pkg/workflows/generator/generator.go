package generator

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

var (
	NonPositiveStagesError = errors.New("number of stages must be positive")
	TooManyStagesError     = errors.New("number of stages can't be that high")
	ZeroFailures           = errors.New("can't create scenario out of 0 failures")
	ZeroTargetsError       = errors.New("no targets available")
	TargetsError           = errors.New("couldn't get list of targets")
	MaxFailuresError       = errors.New("max number of failures must be positive")
	MaxPointsError         = errors.New("max points per stage must be positive")
	RetriesError           = errors.New("retries must be non negative")
)

type Action struct {
	Name     string            `json:"name"`
	Severity metadata.Severity `json:"severity"`
	Scale    metadata.Scale    `json:"scale"`
	Engine   templates.Engine  `json:"engine"`
	Target   targets.Target    `json:"target"`
}

func (a Action) Generate(r *rand.Rand, size int) reflect.Value {
	severity := []metadata.Severity{
		metadata.SeverityHarmless,
		metadata.SeverityLight,
		metadata.SeveritySevere,
		metadata.SeverityCritical,
	}
	scale := []metadata.Scale{
		metadata.ScaleContainer,
		metadata.ScalePod,
		metadata.ScaleDeploymentPart,
		metadata.ScaleDeployment,
		metadata.ScaleNode,
		metadata.ScaleCluster,
	}

	return reflect.ValueOf(Action{
		Name:     strconv.FormatUint(r.Uint64(), 10),
		Severity: severity[r.Intn(len(severity))],
		Scale:    scale[r.Intn(len(scale))],
		Engine:   templates.Engine{}.Generate(r, size).Interface().(templates.Engine),
		Target:   targets.Target{}.Generate(r, size).Interface().(targets.Target),
	})
}

type Stage struct {
	Actions  []Action      `json:"actions"`
	Duration time.Duration `json:"duration"`
}

func (s Stage) Generate(rand *rand.Rand, size int) reflect.Value {
	var actions []Action
	for i := 0; i < rand.Intn(10); i++ {
		actions = append(actions, Action{}.Generate(rand, size).Interface().(Action))
	}

	return reflect.ValueOf(Stage{
		Actions:  actions,
		Duration: time.Duration(30 + rand.Intn(60)),
	})
}

type Scenario struct {
	Stages []Stage `json:"stages"`
}

func (s Scenario) Generate(rand *rand.Rand, size int) reflect.Value {
	var stages []Stage
	for i := 0; i < rand.Intn(10); i++ {
		stages = append(stages, Stage{}.Generate(rand, size).Interface().(Stage))
	}

	return reflect.ValueOf(Scenario{Stages: stages})
}

type Params struct {
	Stages        int
	Seed          int64
	StageDuration time.Duration
}

func (p Params) Generate(rand *rand.Rand, _ int) reflect.Value {
	return reflect.ValueOf(Params{
		Stages:        -10 + rand.Intn(120),
		Seed:          rand.Int63(),
		StageDuration: time.Duration(-10+rand.Int63n(200)) + time.Second,
	})
}

type Generator interface {
	Generate(params Params) (Scenario, error)
}
