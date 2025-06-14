package gocast

// Any is a type alias for interface{} and type converting functionally.
type Any struct {
	v any
}

// Any returns a new Any instance wrapping the provided value.
func (a *Any) Any() any {
	return a.v
}

// Str converts any type to string.
func (a *Any) Str() string {
	return Str(a.v)
}

// Int converts any type to int.
func (a *Any) Int() int {
	return Int(a.v)
}

// Int64 converts any type to int64.
func (a *Any) Int64() int64 {
	return Int64(a.v)
}

// Uint converts any type to uint.
func (a *Any) Uint() uint {
	return Uint(a.v)
}

// Uint64 converts any type to uint64.
func (a *Any) Uint64() uint64 {
	return Uint64(a.v)
}

// Float32 converts any type to float32.
func (a *Any) Float32() float32 {
	return Float32(a.v)
}

// Float64 converts any type to float64.
func (a *Any) Float64() float64 {
	return Float64(a.v)
}

// Bool converts any type to bool.
func (a *Any) Bool() bool {
	return Bool(a.v)
}

// Slice converts any type to []any.
func (a *Any) Slice() []any {
	return AnySlice[any](a.v)
}

// IsSlice checks if the value is a slice.
func (a *Any) IsSlice() bool {
	return IsSlice(a.v)
}

// IsMap checks if the value is a map.
func (a *Any) IsMap() bool {
	return IsMap(a.v)
}

// IsNil checks if the value is nil.
func (a *Any) IsNil() bool {
	return IsNil(a.v)
}

// IsEmpty checks if the value is empty.
func (a *Any) IsEmpty() bool {
	return IsEmpty(a.v)
}
