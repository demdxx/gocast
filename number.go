package gocast

import (
	"reflect"
	"regexp"
	"strconv"

	///"strings"

	"golang.org/x/exp/constraints"
)

// Numeric data type
type Numeric interface {
	constraints.Integer | constraints.Float
}

// /added by simon for parse_string_to_number 2023.1.27
var reg_compare_float_str = regexp.MustCompile(`\.|e|E`)

func parse_string_to_number[R Numeric](str string) (R, error) {
	if reg_compare_float_str.MatchString(str) {
		rval, err := strconv.ParseFloat(str, 64)
		return R(rval), err
	}
	rval, err := strconv.ParseInt(str, 10, 64)
	return R(rval), err
}

///end of added.

// TryNumber converts from types which could be numbers
func TryNumber[R Numeric](v any) (R, error) {
	switch v := v.(type) {
	case nil:
		return R(0), nil
	case R:
		return v, nil
	case *R:
		return *v, nil
	}

	switch v := v.(type) {
	case string:
		///replaced by simon 2023.1.27
		return parse_string_to_number[R](v)
	case []byte:
		s := string(v)
		///replaced by simon 2023.1.27
		return parse_string_to_number[R](s)
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
	///added by simon for ~int,~string etc. 2023.1.27
	default:
		r_value := reflect.ValueOf(v)
		r_value = reflectTarget(r_value)
		switch r_value.Kind() {
		case reflect.String:
			str := r_value.String()
			return parse_string_to_number[R](str)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value := r_value.Int()
			return R(value), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			value := r_value.Uint()
			return R(value), nil
		case reflect.Float32, reflect.Float64:
			value := r_value.Float()
			return R(value), nil
		}
		///end of added.
	}
	return R(0), ErrUnsupportedNumericType
}

// Number converts from types which could be numbers or returns 0
func Number[R Numeric](v any) R {
	res, _ := TryNumber[R](v)
	return res
}
