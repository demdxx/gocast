package gocast

// Or returns the first non-empty value
func Or[T any](a, b T, v ...T) T {
	if !IsEmpty(a) {
		return a
	}
	if IsEmpty(b) {
		for _, val := range v {
			if !IsEmpty(val) {
				return val
			}
		}
	}
	return b
}

// IfThen returns a if cond is true, else b
func IfThen[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

// PtrAsValue returns the value of `v` if `v` is not `nil`, else def
func PtrAsValue[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}
