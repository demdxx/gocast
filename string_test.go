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
	{value: []byte(`byte`), target: "byte"},
	{value: `str`, target: "str"},
}

func TestToStringByReflect(t *testing.T) {
	for _, test := range stringTypecastTests {
		assert.Equal(t, ReflectToString(reflect.ValueOf(test.value)), test.target)
	}
}

func TestToString(t *testing.T) {
	for _, test := range stringTypecastTests {
		assert.Equal(t, ToString(test.value), test.target)
	}
}

func BenchmarkToStringByReflect(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := rand.Intn(len(stringTypecastTests))
			v := reflect.ValueOf(stringTypecastTests[i].value)
			_ = ReflectToString(v)
		}
	})
}

func BenchmarkToString(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := rand.Intn(len(stringTypecastTests))
			_ = ToString(stringTypecastTests[i].value)
		}
	})
}
