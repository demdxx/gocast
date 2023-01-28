package gocast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumber(t *testing.T) {
	assert.Equal(t, int(101), Number[int](int(101)))
	assert.Equal(t, int(101), Number[int](int8(101)))
	assert.Equal(t, int(101), Number[int](int16(101)))
	assert.Equal(t, int(101), Number[int](int32(101)))
	assert.Equal(t, int(101), Number[int](int64(101)))
	assert.Equal(t, int(101), Number[int](uint(101)))
	assert.Equal(t, int(101), Number[int](uint8(101)))
	assert.Equal(t, int(101), Number[int](uint16(101)))
	assert.Equal(t, int(101), Number[int](uint32(101)))
	assert.Equal(t, int(101), Number[int](uint64(101)))
	assert.Equal(t, int(101), Number[int](uintptr(101)))
	assert.Equal(t, int(101), Number[int](float64(101.)))
	assert.Equal(t, int(101), Number[int](float32(101.)))
	assert.Equal(t, int(101), Number[int]("101"))
	assert.Equal(t, int(101), Number[int]("101.0"))
	assert.Equal(t, int(1), Number[int](true))
	assert.Equal(t, int(0), Number[int](false))
	assert.Equal(t, int(0), Number[int](nil))
	assert.Equal(t, 2.35e-01, Number[float64]("2.35E-01"))
	assert.Equal(t, 2.35e-01, Number[float64]("2.35e-01"))
	assert.Equal(t, 2e-01, Number[float64]("2e-01"))

	v, err := TryNumber[int](struct{ S int }{})
	assert.Equal(t, int(0), v)
	assert.Error(t, err)
}
