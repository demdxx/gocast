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

// TryTo cast any input type into the target
func TryTo(v, to any, tags ...string) (any, error) {
	return TryToContext(context.Background(), v, reflect.TypeOf(to), tags...)
}

// TryToContext cast any input type into the target
func TryToContext(ctx context.Context, v, to any, tags ...string) (any, error) {
	if v == nil {
		return nil, ErrInvalidParams
	}
	return TryToTypeContext(ctx, v, reflect.TypeOf(to), tags...)
}

// TryTo cast any input type into the target
func To(v, to any, tags ...string) any {
	val, _ := TryTo(v, to, tags...)
	return val
}

// TryToType cast any input type into the target reflection
func TryToType(v any, t reflect.Type, tags ...string) (any, error) {
	return TryToTypeContext(context.Background(), v, t, tags...)
}

// TryToTypeContext cast any input type into the target reflection
func TryToTypeContext(ctx context.Context, v any, t reflect.Type, tags ...string) (any, error) {
	if v == nil {
		return nil, ErrInvalidParams
	}
	val := reflect.ValueOf(v)
	if t == nil { // In case of type is ANY make a copy of data
		if k := val.Kind(); k == reflect.Struct || k == reflect.Map || k == reflect.Slice || k == reflect.Array {
			return ReflectTryToTypeContext(ctx, val, val.Type(), true, tags...)
		}
		return v, nil
	}
	return ReflectTryToTypeContext(ctx, val, t, true, tags...)
}

// ToType cast any input type into the target reflection
func ToType(v any, t reflect.Type, tags ...string) any {
	val, _ := TryToType(v, t, tags...)
	return val
}

// ReflectTryToType converts reflection value to reflection type or returns error
func ReflectTryToType(v reflect.Value, t reflect.Type, recursive bool, tags ...string) (any, error) {
	return ReflectTryToTypeContext(context.Background(), v, t, recursive, tags...)
}

// ReflectTryToTypeContext converts reflection value to reflection type or returns error
func ReflectTryToTypeContext(ctx context.Context, v reflect.Value, t reflect.Type, recursive bool, tags ...string) (any, error) {
	v = reflectTarget(v)
	if v.Type() == t {
		if k := t.Kind(); k != reflect.Struct &&
			k != reflect.Map &&
			k != reflect.Array && k != reflect.Slice &&
			k != reflect.Interface && k != reflect.Pointer {
			return v.Interface(), nil
		}
	}
	var err error
	switch t.Kind() {
	case reflect.String:
		return Str(v.Interface()), nil
	case reflect.Bool:
		return Bool(v.Interface()), nil
	case reflect.Int:
		return Number[int](v.Interface()), nil
	case reflect.Int8:
		return Number[int8](v.Interface()), nil
	case reflect.Int16:
		return Number[int16](v.Interface()), nil
	case reflect.Int32:
		return Number[int32](v.Interface()), nil
	case reflect.Int64:
		return Number[int64](v.Interface()), nil
	case reflect.Uint:
		return Number[uint](v.Interface()), nil
	case reflect.Uint8:
		return Number[uint8](v.Interface()), nil
	case reflect.Uint16:
		return Number[uint16](v.Interface()), nil
	case reflect.Uint32:
		return Number[uint32](v.Interface()), nil
	case reflect.Uint64:
		return Number[uint64](v.Interface()), nil
	case reflect.Float32:
		return Number[float32](v.Interface()), nil
	case reflect.Float64:
		return Number[float64](v.Interface()), nil
	case reflect.Slice, reflect.Array:
		slice := reflect.New(t)
		if err = TryAnySliceContext(ctx, slice.Interface(), v.Interface(), tags...); err == nil {
			return slice.Elem().Interface(), nil
		}
	case reflect.Map:
		mp := reflect.MakeMap(t).Interface()
		if err = ToMapContext(ctx, mp, v.Interface(), recursive, tags...); err == nil {
			return mp, nil
		}
	case reflect.Interface:
		return v.Interface(), nil
	case reflect.Pointer:
		var (
			vl    any
			tElem = t.Elem()
		)
		if tElem.Kind() == reflect.Struct {
			newVal := reflect.New(tElem)
			if err = TryCopyStructContext(ctx, newVal.Interface(), v.Interface(), tags...); err == nil {
				return newVal.Interface(), nil
			}
		} else if vl, err = ReflectTryToTypeContext(ctx, v, tElem, true, tags...); err == nil {
			return reflect.ValueOf(vl).Addr().Interface(), nil
		}
	case reflect.Struct:
		newVal := reflect.New(t)
		if err = TryCopyStructContext(ctx, newVal.Interface(), v.Interface(), tags...); err == nil {
			return newVal.Elem().Interface(), nil
		}
	default:
		if v.Type() == t {
			newVal := reflect.New(t)
			reflect.Copy(newVal.Addr(), v.Addr())
			return newVal.Interface(), nil
		}
	}
	return nil, err
}

// ReflectToType converts reflection valut to reflection type or returns nil
func ReflectToType(v reflect.Value, t reflect.Type, tags ...string) any {
	return ReflectToTypeContext(context.Background(), v, t, tags...)
}

// ReflectToType converts reflection valut to reflection type or returns nil
func ReflectToTypeContext(ctx context.Context, v reflect.Value, t reflect.Type, tags ...string) any {
	val, _ := ReflectTryToTypeContext(ctx, v, t, true, tags...)
	return val
}

// TryCastValue source type into the target type
func TryCastValue[R any, T any](v T, recursive bool, tags ...string) (R, error) {
	return TryCastValueContext[R](context.Background(), v, recursive, tags...)
}

// TryCastValueContext source type into the target type
func TryCastValueContext[R any, T any](ctx context.Context, v T, recursive bool, tags ...string) (R, error) {
	var (
		rVal     R
		val, err = TryToContext(ctx, v, rVal, tags...)
	)
	if err != nil {
		return rVal, err
	}
	switch nval := val.(type) {
	case *R:
		return *nval, nil
	default:
		return val.(R), nil
	}
}

// TryCastRecursive source type into the target type with recursive data converting
func TryCastRecursive[R any, T any](v T, tags ...string) (R, error) {
	return TryCastRecursiveContext[R](context.Background(), v, tags...)
}

// TryCastRecursiveContext source type into the target type with recursive data converting
func TryCastRecursiveContext[R any, T any](ctx context.Context, v T, tags ...string) (R, error) {
	return TryCastValueContext[R](ctx, v, true, tags...)
}

// TryCast source type into the target type
func TryCast[R any, T any](v T, tags ...string) (R, error) {
	return TryCastContext[R](context.Background(), v, tags...)
}

// TryCastContext source type into the target type
func TryCastContext[R any, T any](ctx context.Context, v T, tags ...string) (R, error) {
	return TryCastValueContext[R](ctx, v, false, tags...)
}

// Cast source type into the target type
func Cast[R any, T any](v T, tags ...string) R {
	return CastContext[R](context.Background(), v, tags...)
}

// CastContext source type into the target type
func CastContext[R any, T any](ctx context.Context, v T, tags ...string) R {
	val, _ := TryCastContext[R](ctx, v, tags...)
	return val
}

// CastRecursive source type into the target type
func CastRecursive[R any, T any](v T, tags ...string) R {
	return CastRecursiveContext[R](context.Background(), v, tags...)
}

// CastRecursiveContext source type into the target type
func CastRecursiveContext[R any, T any](ctx context.Context, v T, tags ...string) R {
	val, _ := TryCastRecursiveContext[R](ctx, v, tags...)
	return val
}
