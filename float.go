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
)

// ReflectToFloat64 returns float64 from reflection
func ReflectToFloat64(v reflect.Value) float64 {
	switch v.Kind() {
	case reflect.String:
		val, _ := strconv.ParseFloat(v.String(), 64)
		return val
	case reflect.Bool:
		if v.Bool() {
			return 1.
		}
		return 0.
	case reflect.Slice, reflect.Array:
		switch v.Type().Elem().Kind() {
		case reflect.Uint8:
			str := string(v.Interface().([]byte))
			fval, _ := strconv.ParseFloat(str, 64)
			return fval
		default:
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	}
	return 0.
}

// Float from any other basic type, float64 by default
//
//go:inline
func Float(v any) float64 {
	return Number[float64](v)
}

// Float64 from any other basic type
//
//go:inline
func Float64(v any) float64 {
	return Number[float64](v)
}

// Float32 from any other basic type
//
//go:inline
func Float32(v any) float32 {
	return Number[float32](v)
}
