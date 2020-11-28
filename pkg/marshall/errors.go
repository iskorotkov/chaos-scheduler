package marshall

import "errors"

//goland:noinspection GoNameStartsWithPackageName
var (
	PropertyConversionError = errors.New("couldn't convert property to string")
	MarshallError           = errors.New("couldn't marshall object")
	UnmarshallError         = errors.New("couldn't unmarshall object")
)
