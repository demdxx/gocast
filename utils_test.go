package gocast

import "testing"

func TestIsEmpty(t *testing.T) {
	var tests = []struct {
		src    interface{}
		target bool
	}{
		{src: 120, target: false},
		{src: uint64(121), target: false},
		{src: 122., target: false},
		{src: "123", target: false},
		{src: "", target: true},
		{src: []byte("125."), target: false},
		{src: []byte(""), target: true},
		{src: true, target: true},
		{src: false, target: false},
		{src: []interface{}{}, target: true},
		{src: []interface{}{1}, target: false},
	}

	for _, test := range tests {
		if v := IsEmpty(test.src); v != test.target {
			t.Errorf("target must be equal %v != %v", test.src, test.target)
		}
	}
}

func BenchmarkIsEmpty(b *testing.B) {
	values := []interface{}{120, int64(122), "123", "120.0", "120.", []byte("125."), true, false}
	for n := 0; n < b.N; n++ {
		_ = IsEmpty(values[n%len(values)])
	}
}
