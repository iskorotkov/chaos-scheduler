package experiments

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"time"
)

type Info struct {
	Name          string
	Lethal        bool
	AffectingNode bool
}

type Preset interface {
	Info() Info
	Engine(target targets.Target, duration time.Duration) Engine
}
