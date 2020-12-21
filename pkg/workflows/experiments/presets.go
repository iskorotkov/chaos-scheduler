package experiments

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"time"
)

type Info struct {
	Name          string `json:"name"`
	Lethal        bool   `json:"lethal"`
	AffectingNode bool   `json:"affectingNode"`
}

type Preset interface {
	Info() Info
	Engine(target targets.Target, duration time.Duration) Engine
}
