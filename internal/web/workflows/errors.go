package workflows

import "errors"

var (
	formParseError          = errors.New("couldn't parse form data")
	scenarioParamsError     = errors.New("couldn't create scenario with given parameters")
	targetsSeekerError      = errors.New("couldn't create targets seeker")
	workflowGenerationError = errors.New("couldn't generate scenario due to unknown reason")
)
