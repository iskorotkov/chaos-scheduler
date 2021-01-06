package workflows

import "errors"

var (
	formParseError          = errors.New("couldn't parse form data")
	scenarioParamsError     = errors.New("couldn't create scenario with given parameters")
	targetsSeekerError      = errors.New("couldn't create targets seeker")
	workflowGenerationError = errors.New("couldn't generate scenario due to unknown reason")
	targetsError            = errors.New("not enough targets present")
	failuresError           = errors.New("not enough failures provided to scenario generator")
)
