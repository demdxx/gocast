package gocast

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCast(t *testing.T) {
	var (
		tests = []any{
			int(100), int8(100), int16(100), int32(100), int64(100),
			uint(100), uint8(100), uint16(100), uint32(100), uint64(100),
			uintptr(100), float32(100.), float64(100.),
			"100", []byte("100"), &[]int32{100}[0],
		}
	)

	assert.Equal(t, 1, Cast[int](true))
	assert.Equal(t, 0, Cast[int](false))
	assert.Equal(t, true, Cast[bool](100))

	// Int
	for _, v := range tests {
		assert.Equal(t, 100, Cast[int](v))
		assert.Equal(t, 100, CastRecursive[int](v))
	}

	// Float
	for _, v := range tests {
		assert.Equal(t, float64(100), Cast[float64](v))
		assert.Equal(t, float64(100), CastRecursive[float64](v))
	}

	// String
	for _, v := range tests {
		assert.Equal(t, "100", strings.TrimSuffix(Cast[string](v), ".0"))
		assert.Equal(t, "100", strings.TrimSuffix(CastRecursive[string](v), ".0"))
	}
}
