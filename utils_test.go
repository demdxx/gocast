package gocast

import "testing"

func TestIsEmpty(t *testing.T) {
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
	}

	for _, test := range tests {
		if v := IsEmpty(test.src); v != test.target {
			t.Errorf("target must be equal %v != %v", test.src, test.target)
		}
	}
}

func BenchmarkIsEmpty(b *testing.B) {
	values := []any{120, int64(122), "123", "120.0", "120.", []byte("125."), true, false}
	for n := 0; n < b.N; n++ {
		_ = IsEmpty(values[n%len(values)])
	}
}
