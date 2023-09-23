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
	case nil:
		return R(0), nil
	case R:
		return v, nil
	case *R:
		return *v, nil
	}
	switch v := v.(type) {
	case string:
		if strings.Contains(v, ".") || strings.Contains(v, "e") || strings.Contains(v, "E") {
			rval, err := strconv.ParseFloat(v, 64)
			return R(rval), err
		}
		rval, err := strconv.ParseInt(v, 10, 64)
		return R(rval), err
	case []byte:
		s := string(v)
		if strings.Contains(s, ".") || strings.Contains(s, "e") || strings.Contains(s, "E") {
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
	}
	return R(0), ErrUnsupportedNumericType
}

// Number converts from types which could be numbers or returns 0
func Number[R Numeric](v any) R {
	res, _ := TryNumber[R](v)
	return res
}

// IsNumeric returns true if input is a numeric
func IsNumericStr(s string) bool {
	if len(s) > 0 && (s[0] == '-' || s[0] == '+') {
		s = s[1:]
	}
	if len(s) == 0 || s[0] == '-' || s[0] == '+' {
		return false
	}
	switch {
	case strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X"):
		return IsStrContainsOf(s[2:], "0123456789abcdefABCDEF")
	case strings.HasPrefix(s, "0o") || strings.HasPrefix(s, "0O") || (s[0] == '0' && len(s) > 1 && s[1] >= '0' && s[1] <= '7'):
		return IsStrContainsOf(s[2:], "01234567")
	case strings.HasPrefix(s, "0b") || strings.HasPrefix(s, "0B"):
		return IsStrContainsOf(s[2:], "01")
	}
	return IsNumeric10Str(s)
}

// IsNumeric10Str returns true if input string is a numeric in base 10
func IsNumeric10Str(s string) bool {
	dot := false
	esimbol := -1
	isStart := false
	for i, c := range s {
		if c < '0' || c > '9' {
			if c == '.' && !dot && esimbol == -1 {
				dot = true
				continue
			}
			if i == 0 && (c == '-' || c == '+') {
				continue
			}
			if (c == 'e' || c == 'E') && isStart && esimbol == -1 {
				esimbol = i
				continue
			}
			if esimbol != -1 && (c == '-' || c == '+') && esimbol == i-1 {
				continue
			}
			return false
		} else {
			isStart = true
		}
	}
	return isStart && esimbol != len(s)-1
}

// IsNumericOnlyStr returns true if input string is a numeric only
func IsNumericOnlyStr(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
