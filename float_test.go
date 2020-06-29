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
		if ToFloat64(test.src) != test.target {
			t.Errorf("target must be equal %v != %f", test.target, test.target)
		}
	}
}
