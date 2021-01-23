package workflows

import "errors"

var (
	formParseError = errors.New("couldn't parse form data")
	paramsError    = errors.New("params are invalid")
	internalError  = errors.New("internal error occurred")
)
