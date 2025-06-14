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
//
//go:inline
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
//
//go:inline
func Slice[R any, S any](src []S, tags ...string) []R {
	return SliceContext[R](context.Background(), src, tags...)
}

// SliceContext converts one type of array to other or resturns nil if not compatible
func SliceContext[R any, S any](ctx context.Context, src []S, tags ...string) []R {
	res, _ := TrySliceContext[R](ctx, src, tags...)
	return res
}

// TryAnySlice converts any input slice into destination type slice as return value
func TryAnySlice[R any](src any, tags ...string) (res []R, err error) {
	return TryAnySliceContext[R](context.Background(), src, tags...)
}

// TryAnySliceContext converts any input slice into destination type slice as return value
func TryAnySliceContext[R any](ctx context.Context, src any, tags ...string) ([]R, error) {
	res := []R{}
	err := TryToAnySliceContext(ctx, &res, src, tags...)
	return res, err
}

// AnySlice converts any input slice into destination type slice as return value
func AnySlice[R any](src any, tags ...string) []R {
	return AnySliceContext[R](context.Background(), src, tags...)
}

// AnySliceContext converts any input slice into destination type slice as return value
func AnySliceContext[R any](ctx context.Context, src any, tags ...string) []R {
	res := []R{}
	_ = TryToAnySliceContext(ctx, &res, src, tags...)
	return res
}

// TryToAnySlice converts any input slice into destination type slice
//
//go:inline
func TryToAnySlice(dst, src any, tags ...string) error {
	return TryToAnySliceContext(context.Background(), dst, src, tags...)
}

// TryToAnySliceContext converts any input slice into destination type slice
func TryToAnySliceContext(ctx context.Context, dst, src any, tags ...string) error {
	if dst == nil || src == nil {
		if dst == nil {
			return wrapError(ErrInvalidParams, "TryToAnySliceContext `destenation` parameter is nil")
		}
		return wrapError(ErrInvalidParams, "TryToAnySliceContext `source` parameter is nil")
	}

	dstSlice := reflectTarget(reflect.ValueOf(dst))
	if k := dstSlice.Kind(); k != reflect.Slice && k != reflect.Array {
		return wrapError(ErrInvalidParams, "TryToAnySliceContext `destenation` parameter is not a slice or array")
	}

	srcSlice := reflectTarget(reflect.ValueOf(src))
	if k := srcSlice.Kind(); k != reflect.Slice && k != reflect.Array {
		return wrapError(ErrInvalidParams, "TryToAnySliceContext `source` parameter is not a slice or array")
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
			if v == nil {
				dstItem.Set(reflect.Zero(dstElemType))
			} else {
				dstItem.Set(reflect.ValueOf(v))
			}
		} else {
			return err
		}
	}

	return nil
}

// IsSlice returns true if v is a slice or array
func IsSlice(v any) bool {
	switch v.(type) {
	// Check default types first for performance
	case []any, []string, []bool,
		[]int, []int8, []int16, []int32, []int64,
		[]uint, []uint8, []uint16, []uint32, []uint64,
		[]float32, []float64:
		return true
	default:
		refValue := reflect.ValueOf(v)
		kind := refValue.Kind()
		return kind == reflect.Slice || kind == reflect.Array
	}
}
