package gocast

import "reflect"

// Len returns size of slice, array or map
func Len[T any](val T) int {
	switch s := any(val).(type) {
	case []any:
		return len(s)
	case []string:
		return len(s)
	case []int:
		return len(s)
	case map[string]any:
		return len(s)
	case map[string]string:
		return len(s)
	default:
		obj := reflectTarget(reflect.ValueOf(val))
		if obj.Kind() == reflect.Array || obj.Kind() == reflect.Slice || obj.Kind() == reflect.Map {
			return obj.Len()
		}
	}
	return -1
}
