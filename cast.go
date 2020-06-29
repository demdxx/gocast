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

// To cast any input type into the target
func To(v, to interface{}, tags string) (interface{}, error) {
	if v == nil || to == nil {
		return nil, errInvalidParams
	}
	return ToT(v, reflect.ValueOf(to).Type(), tags)
}

// ToT cast any input type into the target reflection
func ToT(v interface{}, t reflect.Type, tags string) (interface{}, error) {
	if v == nil || t == nil {
		return nil, errInvalidParams
	}
	vl, err := ToType(reflect.ValueOf(v), t, tags)
	return vl, err
}

// ToType from reflection value to reflection type
func ToType(v reflect.Value, t reflect.Type, tags string) (interface{}, error) {
	if reflect.Ptr == v.Kind() {
		v = reflectTarget(v)
	}

	var err error
	switch t.Kind() {
	case reflect.String:
		return ToString(v.Interface()), nil
	case reflect.Bool:
		return ToBool(v.Interface()), nil
	case reflect.Int:
		return ToInt(v.Interface()), nil
	case reflect.Int8:
		return (int8)(ToInt(v.Interface())), nil
	case reflect.Int16:
		return (int16)(ToInt(v.Interface())), nil
	case reflect.Int32:
		return ToInt32(v.Interface()), nil
	case reflect.Int64:
		return ToInt64(v.Interface()), nil
	case reflect.Uint:
		return ToUint(v.Interface()), nil
	case reflect.Uint8:
		return (uint8)(ToUint(v.Interface())), nil
	case reflect.Uint16:
		return (uint16)(ToUint(v.Interface())), nil
	case reflect.Uint32:
		return ToUint32(v.Interface()), nil
	case reflect.Uint64:
		return ToUint64(v.Interface()), nil
	case reflect.Float32:
		return ToFloat32(v.Interface()), nil
	case reflect.Float64:
		return ToFloat64(v.Interface()), nil
	case reflect.Slice:
		slice := reflect.New(t)
		if err = ToSlice(slice.Interface(), v.Interface(), tags); nil == err {
			return slice.Elem().Interface(), nil
		}
	case reflect.Map:
		mp := reflect.MakeMap(t)
		if err = ToMap(mp.Interface(), v.Interface(), tags, true); nil == err {
			return mp.Interface(), nil
		}
	case reflect.Ptr:
		var vl interface{}
		if vl, err = ToType(v, t.Elem(), tags); nil == err {
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
		if err = ToStruct(st, v.Interface(), tags); nil == err {
			return st, nil
		}
	}
	return nil, err
}
