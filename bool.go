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
	"bytes"
	"reflect"
	"strings"
)

var bytesType = reflect.TypeOf([]byte(nil))

// ReflectToBool returns boolean from reflection
func ReflectToBool(v reflect.Value) bool {
	if !v.IsValid() {
		return false
	}
	switch v.Kind() {
	case reflect.String:
		val := v.String()
		return val == "1" || val == "true" || val == "t"
	case reflect.Slice:
		if v.Type() == bytesType {
			bv := v.Interface().([]byte)
			return len(bv) != 0 && (false ||
				bytes.Equal(bv, []byte("1")) ||
				bytes.Equal(bv, []byte("true")) ||
				bytes.Equal(bv, []byte("t")))
		}
		return v.Len() != 0
	case reflect.Array, reflect.Map:
		return v.Len() != 0
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return v.Float() != 0
	}
	return false
}

// Bool from any other basic types
func Bool(v any) bool {
	switch bv := v.(type) {
	case nil:
		return false
	case string:
		return bv == "1" || bv == "T" || bv == "t" || strings.EqualFold(bv, "true")
	case []byte:
		return len(bv) != 0 && (false ||
			bytes.Equal(bv, []byte("1")) ||
			bytes.Equal(bv, []byte("T")) ||
			bytes.Equal(bv, []byte("t")) ||
			bytes.EqualFold(bv, []byte("true")))
	case bool:
		return bv
	case int:
		return bv != 0
	case int8:
		return bv != 0
	case int16:
		return bv != 0
	case int32:
		return bv != 0
	case int64:
		return bv != 0
	case uint:
		return bv != 0
	case uint8:
		return bv != 0
	case uint16:
		return bv != 0
	case uint32:
		return bv != 0
	case uint64:
		return bv != 0
	case uintptr:
		return bv != 0
	case float32:
		return bv != 0
	case float64:
		return bv != 0
	case []int:
		return len(bv) != 0
	case []int8:
		return len(bv) != 0
	case []int16:
		return len(bv) != 0
	case []int32:
		return len(bv) != 0
	case []uint64:
		return len(bv) != 0
	case []uint:
		return len(bv) != 0
	case []uint16:
		return len(bv) != 0
	case []uint32:
		return len(bv) != 0
	case []int64:
		return len(bv) != 0
	case []float32:
		return len(bv) != 0
	case []float64:
		return len(bv) != 0
	case []bool:
		return len(bv) != 0
	case []any:
		return len(bv) != 0
	}
	return ReflectToBool(reflect.ValueOf(v))
}

// ToBool from any other basic types
//
// Deprecated: Use Bool(v) instead
func ToBool(v any) bool {
	return Bool(v)
}
