package generator

import "errors"

var (
	NonPositiveStagesError = errors.New("number of stages must be positive")
	TooManyStagesError     = errors.New("number of stages can't be that high")
	ZeroFailures           = errors.New("can't create scenario out of 0 failures")
	TargetsError           = errors.New("couldn't get list of targets")
)
