package gocast

import (
	"errors"
	"reflect"
)

// TryCopy creates a deep copy of the provided value using reflection or
// returns an error if the types do not match.
// It returns a new value of the same type as the source.
// If the source value is nil, it returns the zero value of the type.
func TryCopy[T any](src T) (T, error) {
	var dst T

	switch v := any(src).(type) {
	case nil:
	case int:
		*any(&dst).(*int) = v
	case int8:
		*any(&dst).(*int8) = v
	case int16:
		*any(&dst).(*int16) = v
	case int32:
		*any(&dst).(*int32) = v
	case int64:
		*any(&dst).(*int64) = v
	case uint:
		*any(&dst).(*uint) = v
	case uint8:
		*any(&dst).(*uint8) = v
	case uint16:
		*any(&dst).(*uint16) = v
	case uint32:
		*any(&dst).(*uint32) = v
	case uint64:
		*any(&dst).(*uint64) = v
	case uintptr:
		*any(&dst).(*uintptr) = v
	case string:
		*any(&dst).(*string) = v
	case bool:
		*any(&dst).(*bool) = v
	case float32:
		*any(&dst).(*float32) = v
	case float64:
		*any(&dst).(*float64) = v
	default:
		// Use a visited map to handle circular references
		visited := make(map[uintptr]reflect.Value)

		err := deepCopy(reflect.ValueOf(src), reflect.ValueOf(&dst).Elem(), visited)
		if err != nil {
			return dst, err
		}
	}

	return dst, nil
}

// Copy creates a deep copy of the provided value using reflection.
// It returns a new value of the same type as the source.
// If the source value is nil, it returns the zero value of the type.
//
//go:inline
func Copy[T any](src T) T {
	dst, err := TryCopy(src)
	if err != nil {
		panic(err)
	}
	return dst
}

// TryAnyCopy creates a deep copy of the provided value using reflection.
func TryAnyCopy(src any) (any, error) {
	// Use a visited map to handle circular references
	visited := make(map[uintptr]reflect.Value)

	dst := reflect.New(reflect.TypeOf(src)).Elem()
	err := deepCopy(reflect.ValueOf(src), dst, visited)
	if err != nil {
		return nil, err
	}

	return dst.Interface(), nil
}

// AnyCopy creates a deep copy of the provided value using reflection.
//
//go:inline
func AnyCopy(src any) any {
	dst, err := TryAnyCopy(src)
	if err != nil {
		panic(err)
	}
	return dst
}

func deepCopy(src, dst reflect.Value, visited map[uintptr]reflect.Value) error {
	// Handle nil or invalid source values
	if !src.IsValid() {
		return nil
	}

	// Get the underlying value if src is an interface
	if src.Kind() == reflect.Interface && !src.IsNil() {
		srcElem := src.Elem()
		if dst.CanSet() {
			newDst := reflect.New(srcElem.Type()).Elem()
			if err := deepCopy(srcElem, newDst, visited); err != nil {
				return err
			}
			dst.Set(newDst)
		}
		return nil
	}

	// Handle different kinds of values
	switch src.Kind() {
	case reflect.Pointer:
		if src.IsNil() {
			// Nothing to do for nil pointer
			return nil
		}

		// Check for circular references
		if src.Kind() == reflect.Pointer {
			ptr := src.Pointer()
			if v, ok := visited[ptr]; ok {
				dst.Set(v)
				return nil
			}

			// Add to visited before recursing
			visited[ptr] = dst
		}

		// Create a new pointer of the same type as src
		dst.Set(reflect.New(src.Type().Elem()))
		// Recursively copy the pointed-to value
		return deepCopy(src.Elem(), dst.Elem(), visited)

	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			if !dst.Field(i).CanSet() {
				continue // Skip unexported fields
			}
			if err := deepCopy(src.Field(i), dst.Field(i), visited); err != nil {
				return err
			}
		}

	case reflect.Slice:
		if src.IsNil() {
			return nil
		}
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
		for i := 0; i < src.Len(); i++ {
			if err := deepCopy(src.Index(i), dst.Index(i), visited); err != nil {
				return err
			}
		}

	case reflect.Map:
		if src.IsNil() {
			return nil
		}
		dst.Set(reflect.MakeMap(src.Type()))
		iter := src.MapRange()
		for iter.Next() {
			srcKey := iter.Key()
			srcValue := iter.Value()

			dstKey := reflect.New(srcKey.Type()).Elem()
			if err := deepCopy(srcKey, dstKey, visited); err != nil {
				return err
			}

			dstValue := reflect.New(srcValue.Type()).Elem()
			if err := deepCopy(srcValue, dstValue, visited); err != nil {
				return err
			}

			dst.SetMapIndex(dstKey, dstValue)
		}

	case reflect.Chan, reflect.Func:
		return errors.New("unsupported type: " + src.Type().String())

	default:
		// Simple types can be directly set
		if dst.CanSet() {
			dst.Set(src)
		}
	}

	return nil
}
