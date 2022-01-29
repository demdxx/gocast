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
)

// ToInterfaceSlice converts any input slice into Interface type slice
func ToInterfaceSlice(v interface{}) []interface{} {
	switch sv := v.(type) {
	case []interface{}:
		return sv
	default:
		var result []interface{} = nil
		eachSlice(v, func(length int) {
			if length > 0 {
				result = make([]interface{}, length)
			}
		}, func(v interface{}, i int) {
			result[i] = v
		})
		return result
	}
}

// ToStringSlice converts any input slice into String type slice
func ToStringSlice(v interface{}) []string {
	switch sv := v.(type) {
	case []string:
		return sv
	default:
		var result []string = nil
		eachSlice(v, func(length int) {
			if length > 0 {
				result = make([]string, length)
			}
		}, func(v interface{}, i int) {
			result[i] = ToString(v)
		})
		return result
	}
}

// ToIntSlice converts any input slice into Int type slice
func ToIntSlice(v interface{}) []int {
	switch sv := v.(type) {
	case []int:
		return sv
	default:
		var result []int = nil
		eachSlice(v, func(length int) {
			if length > 0 {
				result = make([]int, length)
			}
		}, func(v interface{}, i int) {
			result[i] = ToInt(v)
		})
		return result
	}
}

// ToFloat64Slice converts any input slice into Float64 type slice
func ToFloat64Slice(v interface{}) []float64 {
	switch sv := v.(type) {
	case []float64:
		return sv
	default:
		var result []float64 = nil
		eachSlice(v, func(length int) {
			if length > 0 {
				result = make([]float64, length)
			}
		}, func(v interface{}, i int) {
			result[i] = ToFloat64(v)
		})
		return result
	}
}

// ToSlice converts any input slice into destination type slice
func ToSlice(dst, src interface{}, tags string) error {
	if dst == nil || src == nil {
		return errInvalidParams
	}

	dstSlice := reflectTarget(reflect.ValueOf(dst))
	if reflect.Slice != dstSlice.Kind() {
		return errInvalidParams
	}

	srcSlice := reflectTarget(reflect.ValueOf(src))
	if reflect.Slice != srcSlice.Kind() {
		return errInvalidParams
	}

	dstElemType := dstSlice.Type().Elem()

	if dstSlice.Len() < srcSlice.Len() {
		newv := reflect.MakeSlice(dstSlice.Type(), srcSlice.Len(), srcSlice.Len())
		reflect.Copy(newv, dstSlice)
		dstSlice.Set(newv)
		dstSlice.SetLen(srcSlice.Len())
	}

	for i := 0; i < srcSlice.Len(); i++ {
		it := srcSlice.Index(i)
		if v, err := ToType(it, dstElemType, tags); err == nil {
			val := reflect.ValueOf(v)
			if dstElemType.Kind() != val.Kind() {
				val = val.Elem()
			}
			dstSlice.Index(i).Set(val)
		} else {
			return err
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func eachSlice(v interface{}, fi func(length int), f func(v interface{}, i int)) bool {
	switch sv := v.(type) {
	case []interface{}:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f(v, i)
		}
		// String
	case []string:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((interface{})(v), i)
		}
		// Numeric
	case []int:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((interface{})(v), i)
		}
	case []int64:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((interface{})(v), i)
		}
	case []int32:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((interface{})(v), i)
		}
		// Unsigned numeric
	case []uint:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((interface{})(v), i)
		}
	case []uint64:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((interface{})(v), i)
		}
	case []uint32:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((interface{})(v), i)
		}
		// Float numeric
	case []float32:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((interface{})(v), i)
		}
	case []float64:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((interface{})(v), i)
		}
	case []bool:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((interface{})(v), i)
		}
	default:
		rVal := reflect.ValueOf(sv)
		if k := rVal.Kind(); k == reflect.Slice || k == reflect.Array {
			if fi != nil {
				fi(rVal.Len())
			}
			for i := 0; i < rVal.Len(); i++ {
				f(rVal.Index(i).Interface(), i)
			}
		} else {
			return false
		}
	}
	return true
}
