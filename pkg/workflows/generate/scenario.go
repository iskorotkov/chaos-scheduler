package generate

import (
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
)

// Step is a single scenario step.
type Step struct {
	Name     string                   `json:"name"`
	Type     blueprints.BlueprintType `json:"type"`
	Severity metadata.Severity        `json:"severity"`
	Scale    metadata.Scale           `json:"scale"`
	Engine   blueprints.Engine        `json:"engine"`
	Target   targets.Target           `json:"target"`
}

func (a Step) Generate(r *rand.Rand, size int) reflect.Value {
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

	return reflect.ValueOf(Step{
		Name:     strconv.FormatUint(r.Uint64(), 10),
		Type:     blueprints.BlueprintType(strconv.FormatUint(r.Uint64(), 10)),
		Severity: severity[r.Intn(len(severity))],
		Scale:    scale[r.Intn(len(scale))],
		Engine:   blueprints.Engine{}.Generate(r, size).Interface().(blueprints.Engine),
		Target:   targets.Target{}.Generate(r, size).Interface().(targets.Target),
	})
}

// Stage is a set of actions executed in parallel during specified time.
type Stage struct {
	Steps    []Step        `json:"steps"`
	Duration time.Duration `json:"duration"`
}

func (s Stage) Generate(rand *rand.Rand, size int) reflect.Value {
	var actions []Step
	for i := 0; i < rand.Intn(10); i++ {
		actions = append(actions, Step{}.Generate(rand, size).Interface().(Step))
	}

	return reflect.ValueOf(Stage{
		Steps:    actions,
		Duration: time.Duration(30 + rand.Intn(60)),
	})
}

// Scenario is a complete test scenario.
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
