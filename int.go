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
	"reflect"
	"strconv"
	"strings"
)

// ToInt64ByReflect returns int64 from reflection
func ToInt64ByReflect(v reflect.Value) int64 {
	switch v.Kind() {
	case reflect.String:
		var val int64
		if strings.Contains(v.String(), ".") {
			fval, _ := strconv.ParseFloat(v.String(), 64)
			val = int64(fval)
		} else {
			val, _ = strconv.ParseInt(v.String(), 10, 64)
		}
		return val
	case reflect.Bool:
		if v.Bool() {
			return 1
		}
		return 0
	case reflect.Slice:
		switch v.Type().Elem().Kind() {
		case reflect.Uint8:
			var val int64
			str := string(v.Interface().([]byte))
			if strings.Contains(str, ".") {
				fval, _ := strconv.ParseFloat(str, 64)
				val = int64(fval)
			} else {
				val, _ = strconv.ParseInt(str, 10, 64)
			}
			return val
		default:
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return int64(v.Float())
	}
	return 0
}

// ToInt64 from any other basic types
func ToInt64(v interface{}) int64 {
	switch iv := v.(type) {
	case int:
		return int64(iv)
	case int8:
		return int64(iv)
	case int16:
		return int64(iv)
	case int32:
		return int64(iv)
	case int64:
		return iv
	case uint:
		return int64(iv)
	case uint8:
		return int64(iv)
	case uint16:
		return int64(iv)
	case uint32:
		return int64(iv)
	case uint64:
		return int64(iv)
	case uintptr:
		return int64(iv)
	case float32:
		return int64(iv)
	case float64:
		return int64(iv)
	case bool:
		if iv {
			return 1
		}
		return 0
	case string:
		var val int64
		if strings.Contains(iv, ".") {
			fval, _ := strconv.ParseFloat(iv, 64)
			val = int64(fval)
		} else {
			val, _ = strconv.ParseInt(iv, 10, 64)
		}
		return val
	case []byte:
		var val int64
		str := string(iv)
		if strings.Contains(str, ".") {
			fval, _ := strconv.ParseFloat(str, 64)
			val = int64(fval)
		} else {
			val, _ = strconv.ParseInt(str, 10, 64)
		}
		return val
	}
	return 0
}

// ToInt32 from any other basic types
func ToInt32(v interface{}) int32 {
	return int32(ToInt64(v))
}

// ToInt16 from any other basic types
func ToInt16(v interface{}) int16 {
	return int16(ToInt64(v))
}

// ToInt from any other basic types
func ToInt(v interface{}) int {
	return int(ToInt64(v))
}

// ToUint64ByReflect returns uint64 from reflection
func ToUint64ByReflect(v reflect.Value) uint64 {
	switch v.Kind() {
	case reflect.String:
		var val uint64
		if strings.Contains(v.Interface().(string), ".") {
			fval, _ := strconv.ParseFloat(v.Interface().(string), 64)
			val = uint64(fval)
		} else {
			val, _ = strconv.ParseUint(v.Interface().(string), 10, 64)
		}
		return val
	case reflect.Bool:
		if v.Bool() {
			return 1
		}
		return 0
	case reflect.Slice:
		switch v.Type().Elem().Kind() {
		case reflect.Uint8:
			var val uint64
			str := string(v.Interface().([]byte))
			if strings.Contains(str, ".") {
				fval, _ := strconv.ParseFloat(str, 64)
				val = uint64(fval)
			} else {
				val, _ = strconv.ParseUint(str, 10, 64)
			}
			return val
		default:
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uint64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return uint64(v.Float())
	}
	return 0
}

// ToUint64 from any other basic types
func ToUint64(v interface{}) uint64 {
	switch iv := v.(type) {
	case int:
		return uint64(iv)
	case int8:
		return uint64(iv)
	case int16:
		return uint64(iv)
	case int32:
		return uint64(iv)
	case int64:
		return uint64(iv)
	case uint:
		return uint64(iv)
	case uint8:
		return uint64(iv)
	case uint16:
		return uint64(iv)
	case uint32:
		return uint64(iv)
	case uint64:
		return iv
	case uintptr:
		return uint64(iv)
	case float32:
		return uint64(iv)
	case float64:
		return uint64(iv)
	case bool:
		if iv {
			return 1
		}
		return 0
	case string:
		var val uint64
		if strings.Contains(iv, ".") {
			fval, _ := strconv.ParseFloat(iv, 64)
			val = uint64(fval)
		} else {
			val, _ = strconv.ParseUint(iv, 10, 64)
		}
		return val
	case []byte:
		var val uint64
		str := string(iv)
		if strings.Contains(str, ".") {
			fval, _ := strconv.ParseFloat(str, 64)
			val = uint64(fval)
		} else {
			val, _ = strconv.ParseUint(str, 10, 64)
		}
		return val
	}
	return 0
}

// ToUint32 from any other basic types
func ToUint32(v interface{}) uint32 {
	return uint32(ToUint64(v))
}

// ToUint32 from any other basic types
func ToUint16(v interface{}) uint32 {
	return uint32(ToUint64(v))
}

// ToUint from any other basic types
func ToUint(v interface{}) uint {
	return uint(ToUint64(v))
}
