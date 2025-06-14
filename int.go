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

// ReflectToInt64 returns int64 from reflection
func ReflectToInt64(v reflect.Value) int64 {
	switch v.Kind() {
	case reflect.String:
		var val int64
		var strVal = v.String()
		if strings.Contains(strVal, ".") || strings.Contains(strVal, "e") || strings.Contains(strVal, "E") {
			fval, _ := strconv.ParseFloat(strVal, 64)
			val = int64(fval)
		} else {
			val, _ = strconv.ParseInt(strVal, 10, 64)
		}
		return val
	case reflect.Bool:
		if v.Bool() {
			return 1
		}
		return 0
	case reflect.Slice, reflect.Array:
		switch v.Type().Elem().Kind() {
		case reflect.Uint8:
			var val int64
			str := string(v.Interface().([]byte))
			if strings.Contains(str, ".") || strings.Contains(str, "e") || strings.Contains(str, "E") {
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
	case reflect.Slice, reflect.Array:
		switch v.Type().Elem().Kind() {
		case reflect.Uint8:
			var val uint64
			str := string(v.Interface().([]byte))
			if strings.Contains(str, ".") || strings.Contains(str, "e") || strings.Contains(str, "E") {
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

// Int from any other basic type
//
//go:inline
func Int(v any) int { return Number[int](v) }

// Int8 from any other basic type
//
//go:inline
func Int8(v any) int8 { return Number[int8](v) }

// Int16 from any other basic type
//
//go:inline
func Int16(v any) int16 { return Number[int16](v) }

// Int32 from any other basic type
//
//go:inline
func Int32(v any) int32 { return Number[int32](v) }

// Int64 from any other basic type
//
//go:inline
func Int64(v any) int64 { return Number[int64](v) }

// Uint from any other basic type
//
//go:inline
func Uint(v any) uint { return Number[uint](v) }

// Uint8 from any other basic type
//
//go:inline
func Uint8(v any) uint8 { return Number[uint8](v) }

// Uint16 from any other basic type
//
//go:inline
func Uint16(v any) uint16 { return Number[uint16](v) }

// Uint32 from any other basic type
//
//go:inline
func Uint32(v any) uint32 { return Number[uint32](v) }

// Uint64 from any other basic type
//
//go:inline
func Uint64(v any) uint64 { return Number[uint64](v) }
