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
	"context"
	"reflect"
)

// TrySlice converts one type of array to other or resturns error
func TrySlice[R any, S any](src []S, tags ...string) (res []R, err error) {
	return TrySliceContext[R](context.Background(), src, tags...)
}

// TrySliceContext converts one type of array to other or resturns error
func TrySliceContext[R any, S any](ctx context.Context, src []S, tags ...string) (res []R, err error) {
	res = make([]R, len(src))
	switch srcArr := any(src).(type) {
	case []R:
		copy(res, srcArr)
	default:
		for i, v := range src {
			var newVal R
			if newVal, err = TryCastContext[R](ctx, v); err != nil {
				return nil, err
			} else {
				res[i] = newVal
			}
		}
	}
	return res, nil
}

// Slice converts one type of array to other or resturns nil if not compatible
func Slice[R any, S any](src []S, tags ...string) []R {
	return SliceContext[R](context.Background(), src, tags...)
}

// SliceContext converts one type of array to other or resturns nil if not compatible
func SliceContext[R any, S any](ctx context.Context, src []S, tags ...string) []R {
	res, _ := TrySliceContext[R](ctx, src, tags...)
	return res
}

// ToInterfaceSlice converts any input slice into Interface type slice
//
// Deprecated: Use Slice[any](v) instead
func ToInterfaceSlice(v any) []any {
	switch sv := v.(type) {
	case []any:
		return sv
	default:
		var result []any = nil
		eachSlice(v, func(length int) {
			if length > 0 {
				result = make([]any, length)
			}
		}, func(v any, i int) {
			result[i] = v
		})
		return result
	}
}

// ToStringSlice converts any input slice into String type slice
//
// Deprecated: Use Slice[string](v) instead
func ToStringSlice(v any) []string {
	switch sv := v.(type) {
	case []string:
		return sv
	default:
		var result []string = nil
		eachSlice(v, func(length int) {
			if length > 0 {
				result = make([]string, length)
			}
		}, func(v any, i int) {
			result[i] = Str(v)
		})
		return result
	}
}

// ToIntSlice converts any input slice into Int type slice
//
// Deprecated: Use Slice[int](v) instead
func ToIntSlice(v any) []int {
	switch sv := v.(type) {
	case []int:
		return sv
	default:
		var result []int = nil
		eachSlice(v, func(length int) {
			if length > 0 {
				result = make([]int, length)
			}
		}, func(v any, i int) {
			result[i] = Int(v)
		})
		return result
	}
}

// ToFloat64Slice converts any input slice into Float64 type slice
//
// Deprecated: Use Slice[float64](v) instead
func ToFloat64Slice(v any) []float64 {
	switch sv := v.(type) {
	case []float64:
		return sv
	default:
		var result []float64 = nil
		eachSlice(v, func(length int) {
			if length > 0 {
				result = make([]float64, length)
			}
		}, func(v any, i int) {
			result[i] = Float64(v)
		})
		return result
	}
}

// ToSlice converts any input slice into destination type slice
//
// Deprecated: Use Slice[type](v) or TrySlice[type](v) instead
func ToSlice(dst, src any, tags ...string) error {
	return TryAnySliceContext(context.Background(), dst, src, tags...)
}

// TryAnySliceContext converts any input slice into destination type slice
func TryAnySliceContext(ctx context.Context, dst, src any, tags ...string) error {
	if dst == nil || src == nil {
		return ErrInvalidParams
	}

	dstSlice := reflectTarget(reflect.ValueOf(dst))
	if k := dstSlice.Kind(); k != reflect.Slice && k != reflect.Array {
		return ErrInvalidParams
	}

	srcSlice := reflectTarget(reflect.ValueOf(src))
	if k := srcSlice.Kind(); k != reflect.Slice && k != reflect.Array {
		return ErrInvalidParams
	}

	dstElemType := dstSlice.Type().Elem()

	if dstSlice.Len() < srcSlice.Len() {
		newv := reflect.MakeSlice(dstSlice.Type(), srcSlice.Len(), srcSlice.Len())
		reflect.Copy(newv, dstSlice)
		dstSlice.Set(newv)
		dstSlice.SetLen(srcSlice.Len())
	}

	for i := 0; i < srcSlice.Len(); i++ {
		srcItem := srcSlice.Index(i)
		dstItem := dstSlice.Index(i)
		if setter, _ := dstItem.Interface().(CastSetter); setter != nil {
			if dstItem.Kind() == reflect.Pointer && dstItem.IsNil() {
				dstItem.Set(reflect.New(dstItem.Type().Elem()))
				setter, _ = dstItem.Interface().(CastSetter)
			}
			if err := setter.CastSet(ctx, srcItem.Interface()); err != nil {
				return err
			}
			continue
		} else if dstItem.CanAddr() {
			if setter, _ := dstItem.Addr().Interface().(CastSetter); setter != nil {
				if err := setter.CastSet(ctx, srcItem.Interface()); err != nil {
					return err
				}
				continue
			}
		}
		if v, err := ReflectTryToTypeContext(ctx, srcItem, dstElemType, true, tags...); err == nil {
			dstItem.Set(reflect.ValueOf(v))
		} else {
			return err
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func eachSlice(v any, fi func(length int), f func(v any, i int)) bool {
	switch sv := v.(type) {
	case []any:
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
			f((any)(v), i)
		}
		// Numeric
	case []int:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((any)(v), i)
		}
	case []int64:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((any)(v), i)
		}
	case []int32:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((any)(v), i)
		}
	case []int16:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((any)(v), i)
		}
	case []int8:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((any)(v), i)
		}
		// Unsigned numeric
	case []uint:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((any)(v), i)
		}
	case []uint64:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((any)(v), i)
		}
	case []uint32:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((any)(v), i)
		}
	case []uint16:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((any)(v), i)
		}
	case []uint8:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((any)(v), i)
		}
		// Float numeric
	case []float32:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((any)(v), i)
		}
	case []float64:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((any)(v), i)
		}
	case []bool:
		if fi != nil {
			fi(len(sv))
		}
		for i, v := range sv {
			f((any)(v), i)
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
