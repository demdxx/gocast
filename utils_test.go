package gocast

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	var tests = []struct {
		src    any
		target bool
	}{
		{src: 120, target: false},
		{src: int8(120), target: false},
		{src: int16(120), target: false},
		{src: int32(120), target: false},
		{src: int64(120), target: false},
		{src: uint(121), target: false},
		{src: uint8(121), target: false},
		{src: uint16(121), target: false},
		{src: uint32(121), target: false},
		{src: uint64(121), target: false},
		{src: uintptr(121), target: false},
		{src: float32(122.), target: false},
		{src: float64(122.), target: false},
		{src: "123", target: false},
		{src: "", target: true},
		{src: nil, target: true},
		{src: func() *struct{ s string } { return nil }(), target: true},
		{src: any(nil), target: true},
		{src: any(func() *struct{ s string } { return nil }()), target: true},
		{src: []byte("125."), target: false},
		{src: []byte(""), target: true},
		{src: true, target: false},
		{src: false, target: true},
		{src: []any{}, target: true},
		{src: []any{1}, target: false},
		{src: []int{1}, target: false},
		{src: []int8{1}, target: false},
		{src: []int16{1}, target: false},
		{src: []int32{1}, target: false},
		{src: []int64{1}, target: false},
		{src: []uint{1}, target: false},
		{src: []uint8{1}, target: false},
		{src: []uint16{1}, target: false},
		{src: []uint32{1}, target: false},
		{src: []uint64{1}, target: false},
		{src: []uintptr{1}, target: false},
		{src: []float32{1}, target: false},
		{src: []float64{1}, target: false},
		{src: []bool{}, target: true},
		{src: []string{}, target: true},
		// uintptr
		{src: uintptr(0), target: true},
		// complex
		{src: complex64(0), target: true},
		{src: complex128(0), target: true},
		{src: complex(1, 0), target: false},
		// named/derived types
		{src: myInt(0), target: true},
		{src: myInt(5), target: false},
		{src: myString(""), target: true},
		{src: myString("x"), target: false},
		// struct zero-value
		{src: struct{ A, B int }{}, target: true},
		{src: struct{ A, B int }{A: 1}, target: false},
		// time.Time zero
		{src: time.Time{}, target: true},
		{src: time.Now(), target: false},
		// nil func
		{src: (func())(nil), target: true},
		// nil chan
		{src: (chan int)(nil), target: true},
		// IsEmptyChecker interface
		{src: customEmpty(true), target: true},
		{src: customEmpty(false), target: false},
	}

	t.Run("IsEmpty", func(t *testing.T) {
		for _, test := range tests {
			if v := IsEmpty(test.src); v != test.target {
				t.Errorf("target must be equal %v != %v", test.src, test.target)
			}
		}
	})
	t.Run("IsEmptyByReflection", func(t *testing.T) {
		for _, test := range tests {
			if v := IsEmptyByReflection(reflect.ValueOf(test.src)); v != test.target {
				t.Errorf("target must be equal %v != %v", test.src, test.target)
			}
		}
	})
}

func BenchmarkIsEmpty(b *testing.B) {
	values := []any{120, int64(122), "123", "120.0", "120.", []byte("125."), true, false}
	for n := 0; n < b.N; n++ {
		_ = IsEmpty(values[n%len(values)])
	}
}

func TestIsNil(t *testing.T) {
	tests := []struct {
		src    any
		target bool
	}{
		{src: nil, target: true},
		{src: 120, target: false},
		{src: int8(120), target: false},
		{src: int16(120), target: false},
		{src: int32(120), target: false},
		{src: int64(120), target: false},
		{src: uint(121), target: false},
		{src: uint8(121), target: false},
		{src: uint16(121), target: false},
		{src: uint32(121), target: false},
		{src: uint64(121), target: false},
		{src: uintptr(121), target: false},
		{src: any(nil), target: true},
		// nil nilable kinds
		{src: ([]int)(nil), target: true},
		{src: []int{1, 2}, target: false},
		{src: (map[string]int)(nil), target: true},
		{src: map[string]int{"a": 1}, target: false},
		{src: (chan int)(nil), target: true},
		{src: make(chan int), target: false},
		{src: (func())(nil), target: true},
		{src: func() {}, target: false},
		// nil pointer
		{src: (*int)(nil), target: true},
		{src: func() *int { v := 1; return &v }(), target: false},
	}

	for _, test := range tests {
		assert.Equal(t, test.target, IsNil(test.src), "IsNil failed for %T", test.src)
	}
}

// --- helper types for IsEmpty tests ---

type myInt int
type myString string

// customEmpty implements IsEmptyChecker.
type customEmpty bool

func (c customEmpty) IsEmpty() bool { return bool(c) }
