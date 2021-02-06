// Package failures describes failures structure and hierarchy.
package failures

import (
	"fmt"
	api "github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/blueprints"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/pod"
	"math/rand"
	"reflect"
)

// Failure represents a real-world failure.
type Failure struct {
	// Blueprint is a template used for instantiating concrete failures with targets.
	Blueprint blueprints.Blueprint
	// Scale describes area of damage.
	Scale api.Scale
	// Severity describes level of damage.
	Severity api.Severity
}

// Generate returns random Failure.
func (f Failure) Generate(r *rand.Rand, _ int) reflect.Value {
	scale := []api.Scale{
		api.ScaleContainer,
		api.ScalePod,
		api.ScaleDeploymentPart,
		api.ScaleDeployment,
		api.ScaleNode,
		api.ScaleCluster,
	}
	severity := []api.Severity{
		api.SeverityHarmless,
		api.SeverityLight,
		api.SeveritySevere,
		api.SeverityCritical,
	}

	rs := func(s string) string {
		return fmt.Sprintf("%s-%d", s, r.Int())
	}

	return reflect.ValueOf(Failure{
		Blueprint: pod.Delete{
			Namespace:    rs("namespace"),
			AppNamespace: rs("app-namespace"),
			Interval:     r.Intn(60),
			Force:        r.Int()%2 == 0,
		},
		Scale:    scale[r.Intn(len(scale))],
		Severity: severity[r.Intn(len(severity))],
	})
}

// Name returns a Failure name.
func (f Failure) Name() string {
	return f.Blueprint.Name()
}
