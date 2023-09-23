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

func TestIsNumeric(t *testing.T) {
	trueTests := []string{
		"101", "+.1", "101.0", "2.35E-01", "2.35e-01", "-2e-01", "+2e-01", "1e10", "1e+10",
		"0x0aB1", "0b0101", "0B0101", "0o0123", "0777",
	}
	falseTests := []string{
		"1x10", "-+2e-01", "2Ee-01", "3..14", "3.1.4", "-", ".", "E", "e", "x",
		"0x0aB1x", "0b01b", "0o0123o", "0777x", "0x", "0b", "0o", "0x0aB1x", "0b01b", "0o0123o", "0777x",
		"0x01YT", "0b0101012", "0o01238", "0778",
	}

	for _, v := range trueTests {
		assert.True(t, IsNumericStr(v), v)
	}

	for _, v := range falseTests {
		assert.False(t, IsNumericStr(v), v)
	}

	assert.True(t, IsNumericOnlyStr("0123456789"))
	assert.False(t, IsNumericOnlyStr("0123456789abcdefABCDEF"))
	assert.False(t, IsNumericOnlyStr("0.1"))
}
