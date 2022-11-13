package gocast

import (
	"errors"
)

// Error list...
var (
	errInvalidParams                 = errors.New("invalid params")
	errUnsupportedType               = errors.New("unsupported destination type")
	errUnsupportedSourceType         = errors.New("unsupported source type")
	ErrUnsettableValue               = errors.New("can't set value")
	ErrUnsupportedNumericType        = errors.New("unsupported numeric type")
	ErrStructFieldNameUndefined      = errors.New("struct field name undefined")
	ErrStructFieldValueCantBeChanged = errors.New("struct field value cant be changed")
)
