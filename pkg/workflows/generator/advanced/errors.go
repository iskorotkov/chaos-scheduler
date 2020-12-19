package advanced

import "errors"

var (
	NonPositiveStagesError = errors.New("number of stages must be positive")
	TooManyStagesError     = errors.New("number of stages can't be that high")
	ZeroFailures           = errors.New("can't create scenario out of 0 failures")
	TargetsError           = errors.New("couldn't get list of targets")
	MaxFailuresError       = errors.New("max number of failures must be positive")
	MaxPointsError         = errors.New("max points per stage must be positive")
	RetriesError           = errors.New("retries must be non negative")
	LowTargetsError        = errors.New("can't select enough targets for scenario generation")
)
