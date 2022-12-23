package gocast

import (
	"errors"
)

type errorWrapper struct {
	err error
	msg string
}

func (w *errorWrapper) Error() string { return w.msg + ": " + w.err.Error() }
func (w *errorWrapper) Unwrap() error { return w.err }

func wrapError(err error, msg string) error {
	return &errorWrapper{err: err, msg: msg}
}

// Error list...
var (
	ErrInvalidParams                 = errors.New("invalid params")
	ErrUnsupportedType               = errors.New("unsupported destination type")
	ErrUnsupportedSourceType         = errors.New("unsupported source type")
	ErrUnsettableValue               = errors.New("can't set value")
	ErrUnsupportedNumericType        = errors.New("unsupported numeric type")
	ErrStructFieldNameUndefined      = errors.New("struct field name undefined")
	ErrStructFieldValueCantBeChanged = errors.New("struct field value cant be changed")
)
