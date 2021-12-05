// Package blueprints describes failure blueprints used for instantiating concrete failures with targets.
package blueprints

import (
	"time"

	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
)

const (
	BlueprintTypeNetwork   = "network"
	BlueprintTypeResources = "resources"
	BlueprintTypeIO        = "io"
	BlueprintTypeRestart   = "restart"
)

type BlueprintType string

// Blueprint describes a failure blueprint used for instantiating concrete failures with a target.
type Blueprint interface {
	Name() string
	Type() BlueprintType
	// Instantiate returns a concrete failure with a target.
	Instantiate(target targets.Target, duration time.Duration) Engine
}
