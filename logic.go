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
