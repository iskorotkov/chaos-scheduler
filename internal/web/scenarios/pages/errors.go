package pages

import "errors"

var (
	FormParseError         = errors.New("couldn't parse form data")
	ScenarioExecutionError = errors.New("couldn't execute scenario")
	MarshalError           = errors.New("couldn't marshall workflow to readable format")
	ScenarioParamsError    = errors.New("couldn't create scenario with given parameters")
	ScenarioGeneratorError = errors.New("couldn't generate scenario due to unknown reason")
)
