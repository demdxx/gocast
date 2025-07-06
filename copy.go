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

	srcValue := reflect.ValueOf(src)
	if !srcValue.IsValid() {
		return dst, nil
	}

	// Fast path for simple types
	if srcValue.Type().Comparable() && srcValue.Kind() <= reflect.Complex128 {
		return src, nil // Simple types are immutable, so no copy needed
	}

	// Use a visited map to handle circular references
	visited := make(map[uintptr]reflect.Value)

	err := deepCopy(srcValue, reflect.ValueOf(&dst).Elem(), visited)
	if err != nil {
		return dst, err
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

// CopyWithOptions allows copying with custom options
type CopyOptions struct {
	IgnoreUnexportedFields bool
	MaxDepth               int
	IgnoreCircularRefs     bool
}

// TryCopyWithOptions creates a deep copy with custom options
func TryCopyWithOptions[T any](src T, opts CopyOptions) (T, error) {
	var dst T

	srcValue := reflect.ValueOf(src)
	if !srcValue.IsValid() {
		return dst, nil
	}

	// Fast path for simple types
	if srcValue.Type().Comparable() && srcValue.Kind() <= reflect.Complex128 {
		return src, nil
	}

	// Use a visited map to handle circular references
	visited := make(map[uintptr]reflect.Value)

	err := deepCopyWithOptions(srcValue, reflect.ValueOf(&dst).Elem(), visited, opts, 0)
	if err != nil {
		return dst, err
	}

	return dst, nil
}

func deepCopy(src, dst reflect.Value, visited map[uintptr]reflect.Value) error {
	// Handle nil or invalid source values
	if !src.IsValid() {
		return nil
	}

	// Fast path for simple types that don't need deep copying
	switch src.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		if dst.CanSet() {
			dst.Set(src)
		}
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

	// Handle different kinds of values that need deep copying
	switch src.Kind() {
	case reflect.Pointer:
		return copyPointer(src, dst, visited)
	case reflect.Struct:
		return copyStruct(src, dst, visited)
	case reflect.Slice:
		return copySlice(src, dst, visited)
	case reflect.Array:
		return copyArray(src, dst, visited)
	case reflect.Map:
		return copyMap(src, dst, visited)
	case reflect.Chan, reflect.Func:
		return errors.New("unsupported type: " + src.Type().String())
	default:
		// For other types, try direct assignment
		if dst.CanSet() {
			dst.Set(src)
		}
	}

	return nil
}

func deepCopyWithOptions(src, dst reflect.Value, visited map[uintptr]reflect.Value, opts CopyOptions, depth int) error {
	// Check max depth
	if opts.MaxDepth > 0 && depth >= opts.MaxDepth {
		if dst.CanSet() {
			dst.Set(reflect.Zero(dst.Type()))
		}
		return nil
	}

	// Handle nil or invalid source values
	if !src.IsValid() {
		return nil
	}

	// Fast path for simple types
	switch src.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		if dst.CanSet() {
			dst.Set(src)
		}
		return nil
	case reflect.Interface:
		if !src.IsNil() {
			srcElem := src.Elem()
			if dst.CanSet() {
				newDst := reflect.New(srcElem.Type()).Elem()
				if err := deepCopyWithOptions(srcElem, newDst, visited, opts, depth+1); err != nil {
					return err
				}
				dst.Set(newDst)
			}
		}
		return nil
	case reflect.Pointer:
		return copyPointerWithOptions(src, dst, visited, opts, depth)
	case reflect.Struct:
		return copyStructWithOptions(src, dst, visited, opts, depth)
	case reflect.Slice:
		return copySliceWithOptions(src, dst, visited, opts, depth)
	case reflect.Array:
		return copyArrayWithOptions(src, dst, visited, opts, depth)
	case reflect.Map:
		return copyMapWithOptions(src, dst, visited, opts, depth)
	case reflect.Chan, reflect.Func:
		return wrapError(ErrCopyUnsupportedType, src.Type().String())
	default:
		if dst.CanSet() {
			dst.Set(src)
		}
	}

	return nil
}

func copyPointer(src, dst reflect.Value, visited map[uintptr]reflect.Value) error {
	if src.IsNil() {
		return nil
	}

	// Check for circular references
	ptr := src.Pointer()
	if existingDst, ok := visited[ptr]; ok {
		if dst.CanSet() {
			dst.Set(existingDst)
		}
		return nil
	}

	// Create a new pointer of the same type as src
	newPtr := reflect.New(src.Type().Elem())
	if dst.CanSet() {
		dst.Set(newPtr)
	}

	// Add to visited BEFORE recursing to handle circular references properly
	visited[ptr] = newPtr

	// Recursively copy the pointed-to value
	return deepCopy(src.Elem(), newPtr.Elem(), visited)
}

func copyPointerWithOptions(src, dst reflect.Value, visited map[uintptr]reflect.Value, opts CopyOptions, depth int) error {
	if src.IsNil() {
		return nil
	}

	ptr := src.Pointer()
	if existingDst, ok := visited[ptr]; ok {
		if opts.IgnoreCircularRefs {
			return nil
		}
		if dst.CanSet() {
			dst.Set(existingDst)
		}
		return nil
	}

	newPtr := reflect.New(src.Type().Elem())
	if dst.CanSet() {
		dst.Set(newPtr)
	}

	visited[ptr] = newPtr
	return deepCopyWithOptions(src.Elem(), newPtr.Elem(), visited, opts, depth+1)
}

func copyStruct(src, dst reflect.Value, visited map[uintptr]reflect.Value) error {
	for i := 0; i < src.NumField(); i++ {
		if !dst.Field(i).CanSet() {
			continue // Skip unexported fields
		}
		if err := deepCopy(src.Field(i), dst.Field(i), visited); err != nil {
			return err
		}
	}
	return nil
}

func copyStructWithOptions(src, dst reflect.Value, visited map[uintptr]reflect.Value, opts CopyOptions, depth int) error {
	for i := 0; i < src.NumField(); i++ {
		field := src.Type().Field(i)

		// Skip unexported fields if option is set
		if opts.IgnoreUnexportedFields && !field.IsExported() {
			continue
		}

		if !dst.Field(i).CanSet() {
			continue
		}

		if err := deepCopyWithOptions(src.Field(i), dst.Field(i), visited, opts, depth+1); err != nil {
			return err
		}
	}
	return nil
}

func copySlice(src, dst reflect.Value, visited map[uintptr]reflect.Value) error {
	if src.IsNil() {
		return nil
	}

	newSlice := reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
	dst.Set(newSlice)

	for i := 0; i < src.Len(); i++ {
		if err := deepCopy(src.Index(i), dst.Index(i), visited); err != nil {
			return err
		}
	}
	return nil
}

func copySliceWithOptions(src, dst reflect.Value, visited map[uintptr]reflect.Value, opts CopyOptions, depth int) error {
	if src.IsNil() {
		return nil
	}

	newSlice := reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
	dst.Set(newSlice)

	for i := 0; i < src.Len(); i++ {
		if err := deepCopyWithOptions(src.Index(i), dst.Index(i), visited, opts, depth+1); err != nil {
			return err
		}
	}
	return nil
}

func copyArray(src, dst reflect.Value, visited map[uintptr]reflect.Value) error {
	for i := 0; i < src.Len(); i++ {
		if err := deepCopy(src.Index(i), dst.Index(i), visited); err != nil {
			return err
		}
	}
	return nil
}

func copyArrayWithOptions(src, dst reflect.Value, visited map[uintptr]reflect.Value, opts CopyOptions, depth int) error {
	for i := 0; i < src.Len(); i++ {
		if err := deepCopyWithOptions(src.Index(i), dst.Index(i), visited, opts, depth+1); err != nil {
			return err
		}
	}
	return nil
}

func copyMap(src, dst reflect.Value, visited map[uintptr]reflect.Value) error {
	if src.IsNil() {
		return nil
	}

	newMap := reflect.MakeMap(src.Type())
	dst.Set(newMap)

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
	return nil
}

func copyMapWithOptions(src, dst reflect.Value, visited map[uintptr]reflect.Value, opts CopyOptions, depth int) error {
	if src.IsNil() {
		return nil
	}

	newMap := reflect.MakeMap(src.Type())
	dst.Set(newMap)

	iter := src.MapRange()
	for iter.Next() {
		srcKey := iter.Key()
		srcValue := iter.Value()

		dstKey := reflect.New(srcKey.Type()).Elem()
		if err := deepCopyWithOptions(srcKey, dstKey, visited, opts, depth+1); err != nil {
			return err
		}

		dstValue := reflect.New(srcValue.Type()).Elem()
		if err := deepCopyWithOptions(srcValue, dstValue, visited, opts, depth+1); err != nil {
			return err
		}

		dst.SetMapIndex(dstKey, dstValue)
	}
	return nil
}

// CopySlice creates a deep copy of a slice
func CopySlice[T any](src []T) []T {
	if src == nil {
		return nil
	}
	dst := make([]T, len(src))
	for i, v := range src {
		copied, _ := TryCopy(v)
		dst[i] = copied
	}
	return dst
}

// CopyMap creates a deep copy of a map
func CopyMap[K comparable, V any](src map[K]V) map[K]V {
	if src == nil {
		return nil
	}
	dst := make(map[K]V, len(src))
	for k, v := range src {
		copiedKey, _ := TryCopy(k)
		copiedValue, _ := TryCopy(v)
		dst[copiedKey] = copiedValue
	}
	return dst
}

// MustCopy is like TryCopy but panics on error
func MustCopy[T any](src T) T {
	dst, err := TryCopy(src)
	if err != nil {
		panic(err)
	}
	return dst
}

// ShallowCopy creates a shallow copy (for performance when deep copy is not needed)
func ShallowCopy[T any](src T) T {
	return src
}

// CopyInterface copies an interface{} value with type preservation
func CopyInterface(src any) (any, error) {
	if src == nil {
		return nil, nil
	}

	return TryAnyCopy(src)
}
