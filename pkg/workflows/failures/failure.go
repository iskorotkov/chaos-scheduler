package failures

import (
	api "github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"time"
)

type EngineTemplate interface {
	Name() string
	Instantiate(target targets.Target, duration time.Duration) Engine
}

type Failure struct {
	Template EngineTemplate
	Scale    api.Scale
	Severity api.Severity
}

func (f Failure) Name() string {
	return f.Template.Name()
}
