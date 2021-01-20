package experiments

import (
	"github.com/iskorotkov/chaos-scheduler/api/metadata"
)

type Failure struct {
	Preset   Preset
	Scale    metadata.Scale
	Severity metadata.Severity
}
