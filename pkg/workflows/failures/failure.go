package failures

import (
	"fmt"
	api "github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/templates/pod"
	"math/rand"
	"reflect"
)

type Failure struct {
	Template templates.Template
	Scale    api.Scale
	Severity api.Severity
}

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
		Template: pod.Delete{
			Namespace:    rs("namespace"),
			AppNamespace: rs("app-namespace"),
			Interval:     r.Intn(60),
			Force:        r.Int()%2 == 0,
		},
		Scale:    scale[r.Intn(len(scale))],
		Severity: severity[r.Intn(len(severity))],
	})
}

func (f Failure) Name() string {
	return f.Template.Name()
}
