// Package blueprints describes failure blueprints used for instantiating concrete failures with targets.
package blueprints

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"time"
)

// Blueprint describes a failure blueprint used for instantiating concrete failures with a target.
type Blueprint interface {
	Name() string
	// Instantiate returns a concrete failure with a target.
	Instantiate(target targets.Target, duration time.Duration) Engine
}
