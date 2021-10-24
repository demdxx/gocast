package gocast

import (
	"database/sql/driver"
	"reflect"
)

func reflectTarget(r reflect.Value) reflect.Value {
	for reflect.Ptr == r.Kind() || reflect.Interface == r.Kind() {
		r = r.Elem()
	}
	return r
}

// IsEmpty value
func IsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	switch tv := v.(type) {
	case bool:
		return tv
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
	}
	return IsEmptyByReflection(reflect.ValueOf(v))
}

// IsEmptyByReflection value
func IsEmptyByReflection(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.Slice, reflect.Map, reflect.String:
		return v.Len() < 1
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	}
	return false
}

func getValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	if vl, ok := v.(driver.Valuer); ok {
		v, _ = vl.Value()
	}
	return v
}
