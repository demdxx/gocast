package gocast

import (
	"database/sql/driver"
	"reflect"
)

// IsEmptyChecker is an interface for types that can report their own emptiness.
// Implement this on your custom types to integrate with IsEmpty and IsEmptyByReflection.
type IsEmptyChecker interface {
	IsEmpty() bool
}

// IsEmpty checks value for empty state.
// It supports all primitive types, slices, maps, channels, functions, structs,
// and any type that implements the IsEmptyChecker interface.
func IsEmpty[T any](v T) bool {
	switch tv := any(v).(type) {
	case nil:
		return true
	case IsEmptyChecker:
		return tv.IsEmpty()
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
	case uintptr:
		return tv == 0
	case float32:
		return tv == 0
	case float64:
		return tv == 0
	case complex64:
		return tv == 0
	case complex128:
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
	case []string:
		return len(tv) == 0
	}
	return IsEmptyByReflection(reflect.ValueOf(v))
}

// IsEmptyByReflection checks emptiness via reflection.
// Handles all reflect kinds including struct (zero-value check), func/chan (nil check),
// named types (e.g. type MyInt int), and types implementing IsEmptyChecker.
func IsEmptyByReflection(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	// Allow types to self-report emptiness via the IsEmptyChecker interface.
	if v.CanInterface() {
		if checker, ok := v.Interface().(IsEmptyChecker); ok {
			return checker.IsEmpty()
		}
		// Also check pointer-receiver implementations for addressable values.
		if v.CanAddr() {
			if checker, ok := v.Addr().Interface().(IsEmptyChecker); ok {
				return checker.IsEmpty()
			}
		}
	}
	switch v.Kind() {
	case reflect.Interface, reflect.Pointer:
		if v.IsNil() {
			return true
		}
		return IsEmptyByReflection(v.Elem())
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
		return v.Len() == 0
	case reflect.Chan:
		return v.IsNil()
	case reflect.Func:
		return v.IsNil()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.Struct:
		return v.IsZero()
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
	if !r.IsValid() {
		return r
	}
	for kind := r.Kind(); kind == reflect.Ptr || kind == reflect.Interface; {
		if r.IsNil() {
			break
		}
		if r = r.Elem(); !r.IsValid() {
			break
		}
		kind = r.Kind()
	}
	return r
}

// IsNil reports whether v is nil.
// It returns true for a nil interface, and for nil values of any nilable kind:
// pointer, interface, map, slice, channel, and function.
func IsNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return true
	}
	switch rv.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return rv.IsNil()
	}
	return false
}
