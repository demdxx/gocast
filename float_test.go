package gocast

import "testing"

func TestToFloat(t *testing.T) {
	var tests = []struct {
		src    interface{}
		target float64
	}{
		{src: 120, target: 120},
		{src: uint64(121), target: 121},
		{src: 122., target: 122},
		{src: "123", target: 123},
		{src: "120.0", target: 120},
		{src: "-120.", target: -120},
		{src: []byte("125."), target: 125},
		{src: []byte("125"), target: 125},
		{src: true, target: 1},
		{src: false, target: 0},
	}

	for _, test := range tests {
		if v := ToFloat64(test.src); v != test.target {
			t.Errorf("target must be equal %v != %f", test.src, test.target)
		}
	}
}

func BenchmarkToFloat(b *testing.B) {
	values := []interface{}{120, int64(122), "123", "120.0", "120.", []byte("125."), true, false}
	for n := 0; n < b.N; n++ {
		_ = ToFloat(values[n%len(values)])
	}
}
