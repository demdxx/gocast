package gocast

import (
	"reflect"
	"testing"
)

func TestToInt(t *testing.T) {
	var tests = []struct {
		src    any
		target int
	}{
		{src: 120, target: 120},
		{src: int8(121), target: 121},
		{src: int16(121), target: 121},
		{src: int32(121), target: 121},
		{src: int64(121), target: 121},
		{src: uint(121), target: 121},
		{src: uint8(121), target: 121},
		{src: uint16(121), target: 121},
		{src: uint32(121), target: 121},
		{src: uint64(121), target: 121},
		{src: uintptr(121), target: 121},
		{src: float64(122.), target: 122},
		{src: float32(122.), target: 122},
		{src: "123", target: 123},
		{src: "120.0", target: 120},
		{src: "-120.", target: -120},
		{src: []byte("125."), target: 125},
		{src: []byte("125"), target: 125},
		{src: true, target: 1},
		{src: false, target: 0},
	}

	for _, test := range tests {
		if Int(test.src) != test.target {
			t.Errorf("target must be equal %v != %d", test.target, test.target)
		}
		if Int(test.src) != test.target {
			t.Errorf("target must be equal %v != %d", test.target, test.target)
		}
		if Int8(test.src) != int8(test.target) {
			t.Errorf("target must be equal %v != %d", test.target, test.target)
		}
		if Int16(test.src) != int16(test.target) {
			t.Errorf("target must be equal %v != %d", test.target, test.target)
		}
		if Int32(test.src) != int32(test.target) {
			t.Errorf("target must be equal %v != %d", test.target, test.target)
		}
		if Int64(test.src) != int64(test.target) {
			t.Errorf("target must be equal %v != %d", test.target, test.target)
		}
		if test.target >= 0 {
			if Uint(test.src) != uint(test.target) {
				t.Errorf("target must be equal %v != %d", test.target, test.target)
			}
			if Uint8(test.src) != uint8(test.target) {
				t.Errorf("target must be equal %v != %d", test.target, test.target)
			}
			if Uint16(test.src) != uint16(test.target) {
				t.Errorf("target must be equal %v != %d", test.target, test.target)
			}
			if Uint32(test.src) != uint32(test.target) {
				t.Errorf("target must be equal %v != %d", test.target, test.target)
			}
			if Uint64(test.src) != uint64(test.target) {
				t.Errorf("target must be equal %v != %d", test.target, test.target)
			}
		}
	}
}

func TestToInt64ByReflect(t *testing.T) {
	var tests = []struct {
		src    any
		target int
	}{
		{src: 120, target: 120},
		{src: int8(121), target: 121},
		{src: int16(121), target: 121},
		{src: int32(121), target: 121},
		{src: int64(121), target: 121},
		{src: uint(121), target: 121},
		{src: uint8(121), target: 121},
		{src: uint16(121), target: 121},
		{src: uint32(121), target: 121},
		{src: uint64(121), target: 121},
		{src: uintptr(121), target: 121},
		{src: float64(122.), target: 122},
		{src: float32(122.), target: 122},
		{src: "123", target: 123},
		{src: "120.0", target: 120},
		{src: "-120.", target: -120},
		{src: []byte("125."), target: 125},
		{src: []byte("125"), target: 125},
		{src: true, target: 1},
		{src: false, target: 0},
	}

	for _, test := range tests {
		if int(ReflectToInt64(reflect.ValueOf(test.src))) != test.target {
			t.Errorf("target must be equal %v != %d", test.target, test.target)
		}
	}
}

func BenchmarkToInt(b *testing.B) {
	values := []any{120, uint64(122), "123", "120.0", "-120.", []byte("125."), true, false}
	for n := 0; n < b.N; n++ {
		_ = Int(values[n%len(values)])
	}
}

func TestToUInt(t *testing.T) {
	var tests = []struct {
		src    any
		target uint
	}{
		{src: 120, target: 120},
		{src: int8(121), target: 121},
		{src: int16(121), target: 121},
		{src: int32(121), target: 121},
		{src: int64(121), target: 121},
		{src: uint(121), target: 121},
		{src: uint8(121), target: 121},
		{src: uint16(121), target: 121},
		{src: uint32(121), target: 121},
		{src: uint64(121), target: 121},
		{src: uintptr(121), target: 121},
		{src: float64(122.), target: 122},
		{src: float32(122.), target: 122},
		{src: "123", target: 123},
		{src: "120.0", target: 120},
		{src: "120.", target: 120},
		{src: []byte("125."), target: 125},
		{src: []byte("125"), target: 125},
		{src: true, target: 1},
		{src: false, target: 0},
	}

	for _, test := range tests {
		if Uint(test.src) != test.target {
			t.Errorf("target must be equal %v != %d", test.target, test.target)
		}
	}
}

func TestToUint64ByReflect(t *testing.T) {
	var tests = []struct {
		src    any
		target uint
	}{
		{src: 120, target: 120},
		{src: int8(121), target: 121},
		{src: int16(121), target: 121},
		{src: int32(121), target: 121},
		{src: int64(121), target: 121},
		{src: uint(121), target: 121},
		{src: uint8(121), target: 121},
		{src: uint16(121), target: 121},
		{src: uint32(121), target: 121},
		{src: uint64(121), target: 121},
		{src: uintptr(121), target: 121},
		{src: float64(122.), target: 122},
		{src: float32(122.), target: 122},
		{src: "123", target: 123},
		{src: "120.0", target: 120},
		{src: "120.", target: 120},
		{src: []byte("125."), target: 125},
		{src: []byte("125"), target: 125},
		{src: true, target: 1},
		{src: false, target: 0},
	}

	for _, test := range tests {
		if uint(ToUint64ByReflect(reflect.ValueOf(test.src))) != test.target {
			t.Errorf("target must be equal %v != %d", test.target, test.target)
		}
	}
}

func BenchmarkToUint(b *testing.B) {
	values := []any{120, int64(122), "123", "120.0", "120.", []byte("125."), true, false}
	for n := 0; n < b.N; n++ {
		_ = Uint(values[n%len(values)])
	}
}
