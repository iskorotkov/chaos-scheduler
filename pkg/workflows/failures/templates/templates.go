package templates

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"time"
)

type Template interface {
	Name() string
	Instantiate(target targets.Target, duration time.Duration) Engine
}
