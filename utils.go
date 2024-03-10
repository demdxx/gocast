package gocast

import (
	"database/sql/driver"
	"reflect"
)

// IsEmpty checks value for empty state
func IsEmpty[T any](v T) bool {
	switch tv := any(v).(type) {
	case nil:
		return true
	case bool:
		return !tv
	case int:
		return tv == 0
	case int8:
		return tv == 0
	case int16:
		return tv == 0
	case int32:
		return tv == 0
	case int64:
		return tv == 0
	case uint:
		return tv == 0
	case uint8:
		return tv == 0
	case uint16:
		return tv == 0
	case uint32:
		return tv == 0
	case uint64:
		return tv == 0
	case float32:
		return tv == 0
	case float64:
		return tv == 0
	case string:
		return tv == ""
	case []int:
		return len(tv) == 0
	case []int8:
		return len(tv) == 0
	case []int16:
		return len(tv) == 0
	case []int32:
		return len(tv) == 0
	case []int64:
		return len(tv) == 0
	case []uint:
		return len(tv) == 0
	case []uint8:
		return len(tv) == 0
	case []uint16:
		return len(tv) == 0
	case []uint32:
		return len(tv) == 0
	case []uint64:
		return len(tv) == 0
	case []float32:
		return len(tv) == 0
	case []float64:
		return len(tv) == 0
	case []any:
		return len(tv) == 0
	case []bool:
		return len(tv) == 0
	}
	return IsEmptyByReflection(reflect.ValueOf(v))
}

// IsEmptyByReflection value
func IsEmptyByReflection(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	}
	return false
}

func getValue(v any) any {
	if v == nil {
		return nil
	}
	if vl, ok := v.(driver.Valuer); ok {
		v, _ = vl.Value()
	}
	return v
}

func reflectTarget(r reflect.Value) reflect.Value {
	for kind := r.Kind(); kind == reflect.Ptr || kind == reflect.Interface; {
		r = r.Elem()
		kind = r.Kind()
	}
	return r
}
