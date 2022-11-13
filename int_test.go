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
		if ToInt(test.src) != test.target {
			t.Errorf("target must be equal %v != %d", test.target, test.target)
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
		_ = ToInt(values[n%len(values)])
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
		if ToUint(test.src) != test.target {
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
		_ = ToUint(values[n%len(values)])
	}
}
