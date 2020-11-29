package templates

import (
	"time"
)

type SuspendSection struct {
	Duration string `yaml:"duration" json:"duration"`
}

type SuspendTemplate struct {
	Name    string         `yaml:"name" json:"name"`
	Suspend SuspendSection `yaml:"suspend" json:"suspend"`
}

func NewSuspendTemplate(name string, duration time.Duration) SuspendTemplate {
	suspendFor := SuspendSection{Duration: duration.String()}
	return SuspendTemplate{name, suspendFor}
}
