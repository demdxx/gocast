package gocast

import "testing"

func TestToInt(t *testing.T) {
	var tests = []struct {
		src    interface{}
		target int
	}{
		{src: 120, target: 120},
		{src: int64(121), target: 121},
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
		if ToInt(test.src) != test.target {
			t.Errorf("target must be equal %v != %d", test.target, test.target)
		}
	}
}

func TestToUInt(t *testing.T) {
	var tests = []struct {
		src    interface{}
		target uint
	}{
		{src: 120, target: 120},
		{src: int64(121), target: 121},
		{src: 122., target: 122},
		{src: "123", target: 123},
		{src: "120.0", target: 120},
		{src: "120.", target: 120},
		{src: []byte("125."), target: 125},
		{src: []byte("125"), target: 125},
		{src: true, target: 1},
		{src: false, target: 0},
	}

	for _, test := range tests {
		if ToUint(test.src) != test.target {
			t.Errorf("target must be equal %v != %d", test.target, test.target)
		}
	}
}
