package gocast

import (
	"reflect"
	"testing"
)

func TestToBool(t *testing.T) {
	var tests = []struct {
		src    interface{}
		target bool
	}{
		{src: 120, target: true},
		{src: 0, target: false},
		{src: uint64(121), target: true},
		{src: 122., target: true},
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

func BenchmarkToBool(b *testing.B) {
	values := []interface{}{120, uint64(122), "f", "true", "", []byte("t"), true, false, 0.}
	for n := 0; n < b.N; n++ {
		_ = ToBool(values[n%len(values)])
	}
}

func BenchmarkToBoolByReflect(b *testing.B) {
	var (
		baseValues = []interface{}{120, uint64(122), "f", "true", "", []byte("t"), true, false, 0.}
		values     = []reflect.Value{}
	)
	for _, v := range baseValues {
		values = append(values, reflect.ValueOf(v))
	}
	for n := 0; n < b.N; n++ {
		_ = ToBoolByReflect(values[n%len(values)])
	}
}
