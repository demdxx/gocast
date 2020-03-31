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
	t := reflect.ValueOf(v)
	switch t.Kind() {
	case reflect.Slice, reflect.Map, reflect.String:
		return t.Len() < 1
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return t.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return t.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return t.Float() == 0
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
