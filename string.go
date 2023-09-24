// Copyright (c) 2014 Dmitry Ponomarev <demdxx@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package gocast

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// TryReflectStr converts reflection value to string
func TryReflectStr(v reflect.Value) (string, error) {
	if !v.IsValid() {
		return ``, nil
	}
	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Slice, reflect.Array:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return string(v.Bytes()), nil
		}
	case reflect.Bool:
		if v.Bool() {
			return "true", nil
		}
		return "false", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'G', -1, 64), nil
	}
	return fmt.Sprintf("%v", v.Interface()), nil
}

// ReflectStr converts reflection value to string
func ReflectStr(v reflect.Value) string {
	s, _ := TryReflectStr(v)
	return s
}

// TryReflectToString converts reflection value to string
//
// Deprecated: Use TryReflectStr instead
func TryReflectToString(v reflect.Value) (string, error) {
	return TryReflectStr(v)
}

// ReflectToString converts reflection value to string
//
// Deprecated: Use TryReflectStr instead
func ReflectToString(v reflect.Value) string {
	return ReflectStr(v)
}

// TryStr from any type
func TryStr(v any) (string, error) {
	switch val := v.(type) {
	case nil:
		return ``, nil
	case string:
		return val, nil
	case []byte:
		return string(val), nil
	case bool:
		if val {
			return "true", nil
		}
		return "false", nil
	case int:
		return strconv.FormatInt(int64(val), 10), nil
	case int8:
		return strconv.FormatInt(int64(val), 10), nil
	case int16:
		return strconv.FormatInt(int64(val), 10), nil
	case int32:
		return strconv.FormatInt(int64(val), 10), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case uint:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint64:
		return strconv.FormatUint(val, 10), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'G', -1, 64), nil
	case float64:
		return strconv.FormatFloat(val, 'G', -1, 64), nil
	case reflect.Value:
		return TryReflectStr(reflectTarget(val))
	}
	val := reflectTarget(reflect.ValueOf(v))
	return fmt.Sprintf("%v", val.Interface()), nil
}

// Str returns string value from any type
func Str(v any) string {
	s, _ := TryStr(v)
	return s
}

// ToString from any type
//
// Deprecated: Use Str instead
func ToString(v any) string {
	return Str(v)
}

// TryToString returns string value from any type
//
// Deprecated: Use TryStr instead
func TryToString(v any) (string, error) {
	return TryStr(v)
}

// IsStr returns true if value is string
func IsStr(v any) bool {
	if v == nil {
		return false
	}
	_, ok := v.(string)
	return ok
}

// IsStrContainsOf returns true if input string contains only chars from subset
func IsStrContainsOf(s, subset string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !strings.ContainsRune(subset, c) {
			return false
		}
	}
	return true
}
