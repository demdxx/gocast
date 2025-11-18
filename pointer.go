package gocast

// PtrAsValue returns the value of `v` if `v` is not `nil`, else def
//
//go:inline
func PtrAsValue[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}

// Ptr returns the pointer of v
//
//go:inline
func Ptr[T any](v T) *T {
	return &v
}
