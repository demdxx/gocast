package gocast

import (
  "errors"
)

var (
  errInvalidParams         = errors.New("Invalid params")
  errUnsupportedType       = errors.New("Unsupported destination type")
  errUnsupportedSourceType = errors.New("Unsupported source type")
)
