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
	switch v.(type) {
	case []interface{}:
		return v.([]interface{})
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
	switch v.(type) {
	case []string:
		return v.([]string)
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
	switch v.(type) {
	case []int:
		return v.([]int)
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
	switch v.(type) {
	case []float64:
		return v.([]float64)
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

func eachSlice(v interface{}, fi func(length int), f func(v interface{}, i int)) {
	switch v.(type) {
	case []interface{}:
		if fi != nil {
			fi(len(v.([]interface{})))
		}
		for i, v := range v.([]interface{}) {
			f(v, i)
		}
		break
		// String
	case []string:
		if fi != nil {
			fi(len(v.([]string)))
		}
		for i, v := range v.([]string) {
			f((interface{})(v), i)
		}
		break
		// Numeric
	case []int:
		if fi != nil {
			fi(len(v.([]int)))
		}
		for i, v := range v.([]int) {
			f((interface{})(v), i)
		}
		break
	case []int64:
		if fi != nil {
			fi(len(v.([]int64)))
		}
		for i, v := range v.([]int64) {
			f((interface{})(v), i)
		}
		break
	case []int32:
		if fi != nil {
			fi(len(v.([]int32)))
		}
		for i, v := range v.([]int32) {
			f((interface{})(v), i)
		}
		break
		// Unsigned numeric
	case []uint:
		if fi != nil {
			fi(len(v.([]uint)))
		}
		for i, v := range v.([]uint) {
			f((interface{})(v), i)
		}
		break
	case []uint64:
		if fi != nil {
			fi(len(v.([]uint64)))
		}
		for i, v := range v.([]uint64) {
			f((interface{})(v), i)
		}
		break
	case []uint32:
		if fi != nil {
			fi(len(v.([]uint32)))
		}
		for i, v := range v.([]uint32) {
			f((interface{})(v), i)
		}
		break
		// Float numeric
	case []float32:
		if fi != nil {
			fi(len(v.([]float32)))
		}
		for i, v := range v.([]float32) {
			f((interface{})(v), i)
		}
		break
	case []float64:
		if fi != nil {
			fi(len(v.([]float64)))
		}
		for i, v := range v.([]float64) {
			f((interface{})(v), i)
		}
		break
	case []bool:
		if fi != nil {
			fi(len(v.([]bool)))
		}
		for i, v := range v.([]bool) {
			f((interface{})(v), i)
		}
		break
	}
}
