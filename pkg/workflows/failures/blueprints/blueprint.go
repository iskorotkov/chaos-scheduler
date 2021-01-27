package blueprints

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"time"
)

type Blueprint interface {
	Name() string
	Instantiate(target targets.Target, duration time.Duration) Engine
}
