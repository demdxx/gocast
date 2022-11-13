package gocast

import (
	"reflect"
	"testing"
)

func TestToBool(t *testing.T) {
	var tests = []struct {
		src    any
		target bool
	}{
		{src: 120, target: true},
		{src: 0, target: false},
		{src: int8(121), target: true},
		{src: int16(121), target: true},
		{src: int32(121), target: true},
		{src: int64(121), target: true},
		{src: uint(121), target: true},
		{src: uint8(121), target: true},
		{src: uint16(121), target: true},
		{src: uint32(121), target: true},
		{src: uint64(121), target: true},
		{src: float64(122.), target: true},
		{src: float32(122.), target: true},
		{src: "t", target: true},
		{src: "", target: false},
		{src: []byte("true"), target: true},
		{src: []byte(""), target: false},
		{src: true, target: true},
		{src: false, target: false},
		{src: []int{0}, target: true},
		{src: []int{}, target: false},
	}
	for _, test := range tests {
		if v := ToBool(test.src); v != test.target {
			t.Errorf("target must be equal %v != %v", test.src, test.target)
		}
	}
}

func BenchmarkBool(b *testing.B) {
	values := []any{120, uint64(122), "f", "true", "", []byte("t"), true, false, 0.}
	for n := 0; n < b.N; n++ {
		_ = Bool(values[n%len(values)])
	}
}

func BenchmarkToBoolByReflect(b *testing.B) {
	var (
		baseValues = []any{120, uint64(122), "f", "true", "", []byte("t"), true, false, 0.}
		values     = []reflect.Value{}
	)
	for _, v := range baseValues {
		values = append(values, reflect.ValueOf(v))
	}
	for n := 0; n < b.N; n++ {
		_ = ReflectToBool(values[n%len(values)])
	}
}
