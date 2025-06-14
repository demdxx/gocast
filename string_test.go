package gocast

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var stringTypecastTests = []struct {
	value  any
	target string
}{
	{value: 1, target: "1"},
	{value: int8(1), target: "1"},
	{value: int16(1), target: "1"},
	{value: int32(1), target: "1"},
	{value: int64(1), target: "1"},
	{value: uint(1), target: "1"},
	{value: uint8(1), target: "1"},
	{value: uint16(1), target: "1"},
	{value: uint32(1), target: "1"},
	{value: uint64(1), target: "1"},
	{value: 1.1, target: "1.1"},
	{value: float32(1.5), target: "1.5"},
	{value: true, target: "true"},
	{value: false, target: "false"},
	{value: []byte(`byte`), target: "byte"},
	{value: `str`, target: "str"},
	{value: nil, target: ""},
}

func TestToStringByReflect(t *testing.T) {
	for _, test := range stringTypecastTests {
		assert.Equal(t, ReflectStr(reflect.ValueOf(test.value)), test.target)
	}
}

func TestToString(t *testing.T) {
	for _, test := range stringTypecastTests {
		assert.Equal(t, Str(test.value), test.target)
	}
}

func TestIsStr(t *testing.T) {
	tests := []struct {
		value  any
		target bool
	}{
		{value: 1, target: false},
		{value: nil, target: false},
		{value: int8(1), target: false},
		{value: int16(1), target: false},
		{value: []byte("notstr"), target: false},
		{value: []int8{1, 2, 3, 4, 5}, target: false},
		{value: []any{'1', '2', '3'}, target: false},
		{value: "str", target: true},
	}
	for _, test := range tests {
		assert.Equal(t, IsStr(test.value), test.target)
	}
}

func BenchmarkToStringByReflect(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := rand.Intn(len(stringTypecastTests))
			v := reflect.ValueOf(stringTypecastTests[i].value)
			_ = ReflectStr(v)
		}
	})
}

func BenchmarkToString(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := rand.Intn(len(stringTypecastTests))
			_ = Str(stringTypecastTests[i].value)
		}
	})
}
