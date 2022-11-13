package gocast

import (
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

// Numeric data type
type Numeric interface {
	constraints.Integer | constraints.Float
}

// TryNumber converts from types which could be numbers
func TryNumber[R Numeric](v any) (R, error) {
	switch v := v.(type) {
	case string:
		if strings.Contains(v, ".") {
			rval, err := strconv.ParseFloat(v, 64)
			return R(rval), err
		}
		rval, err := strconv.ParseInt(v, 10, 64)
		return R(rval), err
	case []byte:
		s := string(v)
		if strings.Contains(s, ".") {
			rval, err := strconv.ParseFloat(s, 64)
			return R(rval), err
		}
		rval, err := strconv.ParseInt(s, 10, 64)
		return R(rval), err
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case int:
		return R(v), nil
	case int8:
		return R(v), nil
	case int16:
		return R(v), nil
	case int32:
		return R(v), nil
	case int64:
		return R(v), nil
	case uint:
		return R(v), nil
	case uint8:
		return R(v), nil
	case uint16:
		return R(v), nil
	case uint32:
		return R(v), nil
	case uintptr:
		return R(v), nil
	case uint64:
		return R(v), nil
	case float32:
		return R(v), nil
	case float64:
		return R(v), nil
	case nil:
		return R(0), nil
	}
	return R(0), ErrUnsupportedNumericType
}

// Number converts from types which could be numbers or returns 0
func Number[R Numeric](v any) R {
	res, _ := TryNumber[R](v)
	return res
}
