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

// TryTo cast any input type into the target
func TryTo(v, to any, tags ...string) (any, error) {
	if v == nil || to == nil {
		return nil, ErrInvalidParams
	}
	return TryToType(v, reflect.ValueOf(to).Type(), tags...)
}

// TryTo cast any input type into the target
func To(v, to any, tags ...string) any {
	val, _ := TryTo(v, to, tags...)
	return val
}

// TryToType cast any input type into the target reflection
func TryToType(v any, t reflect.Type, tags ...string) (any, error) {
	if v == nil || t == nil {
		return nil, ErrInvalidParams
	}
	vl, err := ReflectTryToType(reflect.ValueOf(v), t, true, tags...)
	return vl, err
}

// ToType cast any input type into the target reflection
func ToType(v any, t reflect.Type, tags ...string) any {
	val, _ := TryToType(v, t, tags...)
	return val
}

// ReflectTryToType converts reflection value to reflection type or returns error
func ReflectTryToType(v reflect.Value, t reflect.Type, recursive bool, tags ...string) (any, error) {
	v = reflectTarget(v)
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
	case reflect.Slice:
		slice := reflect.New(t)
		if err = ToSlice(slice.Interface(), v.Interface(), tags...); err == nil {
			return slice.Elem().Interface(), nil
		}
	case reflect.Map:
		mp := reflect.MakeMap(t).Interface()
		if err = ToMap(mp, v.Interface(), recursive, tags...); err == nil {
			return mp, nil
		}
	case reflect.Interface:
		return v.Interface(), nil
	case reflect.Ptr:
		var vl any
		if vl, err = ReflectTryToType(v, t.Elem(), true, tags...); err == nil {
			if reflect.Ptr != t.Elem().Kind() {
				return vl, nil
			}
			vlPtr := reflect.Zero(t)
			val := reflect.ValueOf(vl)
			vlPtr.Set(val)
			return vlPtr.Interface(), nil
		}
	case reflect.Struct:
		st := reflect.New(t).Interface()
		if err = TryCopyStruct(st, v.Interface(), tags...); err == nil {
			return st, nil
		}
	}
	return nil, err
}

// ReflectToType converts reflection valut to reflection type or returns nil
func ReflectToType(v reflect.Value, t reflect.Type, tags ...string) any {
	val, _ := ReflectTryToType(v, t, true, tags...)
	return val
}

// TryCastValue source type into the target type
func TryCastValue[R any, T any](v T, recursive bool, tags ...string) (R, error) {
	var (
		rVal     R
		val, err = TryTo(v, rVal, tags...)
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
	return TryCastValue[R](v, true, tags...)
}

// TryCast source type into the target type
func TryCast[R any, T any](v T, tags ...string) (R, error) {
	return TryCastValue[R](v, false, tags...)
}

// Cast source type into the target type
func Cast[R any, T any](v T, tags ...string) R {
	val, _ := TryCast[R](v, tags...)
	return val
}

// CastRecursive source type into the target type
func CastRecursive[R any, T any](v T, tags ...string) R {
	val, _ := TryCastRecursive[R](v, tags...)
	return val
}
