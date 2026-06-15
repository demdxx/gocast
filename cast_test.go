package gocast

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCast(t *testing.T) {
	type testStructType struct {
		S string `json:"s"`
		A int    `json:"a"`
		B bool   `json:"b"`
	}
	var (
		tests = []any{
			int(100), int8(100), int16(100), int32(100), int64(100),
			uint(100), uint8(100), uint16(100), uint32(100), uint64(100),
			uintptr(100), float32(100.), float64(100.),
			"100", []byte("100"), &[]int32{100}[0],
		}
		testStruct = testStructType{S: "test", A: 1, B: true}
	)

	t.Run("simple", func(t *testing.T) {
		assert.Equal(t, 1, Cast[int](true))
		assert.Equal(t, int8(0), Cast[int8](false))
		assert.Equal(t, int16(100), Cast[int16]("100"))
		assert.Equal(t, int32(100), Cast[int32]("100"))
		assert.Equal(t, int64(100), Cast[int64]("100"))
		assert.Equal(t, uint(1), Cast[uint](true))
		assert.Equal(t, uint8(0), Cast[uint8](false))
		assert.Equal(t, uint16(100), Cast[uint16]("100"))
		assert.Equal(t, uint32(100), Cast[uint32]("100"))
		assert.Equal(t, uint64(100), Cast[uint64]("100"))
		assert.Equal(t, uint64(100), Cast[uint64](&[]int32{100}[0]))
		assert.Equal(t, true, Cast[bool](100))
		assert.Equal(t, testStructType{}, ToType(testStructType{}, nil))
	})

	t.Run("int", func(t *testing.T) {
		for _, v := range tests {
			assert.Equal(t, 100, Cast[int](v))
			assert.Equal(t, 100, CastRecursive[int](v))
		}
	})

	t.Run("float", func(t *testing.T) {
		for _, v := range tests {
			assert.Equal(t, float32(100), Cast[float32](v))
			assert.Equal(t, float32(100), CastRecursive[float32](v))
			assert.Equal(t, float64(100), Cast[float64](v))
			assert.Equal(t, float64(100), CastRecursive[float64](v))
		}
	})

	t.Run("string", func(t *testing.T) {
		for _, v := range tests {
			assert.Equal(t, "100", strings.TrimSuffix(Cast[string](v), ".0"))
			assert.Equal(t, "100", strings.TrimSuffix(CastRecursive[string](v), ".0"))
		}
	})

	t.Run("struct", func(t *testing.T) {
		mp := Cast[map[string]any](&testStruct, "json")
		if assert.NotNil(t, mp) {
			newStruct, err := TryCast[testStructType](mp, "json")
			assert.NoError(t, err)
			assert.Equal(t, testStruct, newStruct)
		}
	})
}

func TestTryTo(t *testing.T) {
	t.Run("basic conversion", func(t *testing.T) {
		var target int
		v, err := TryTo(42.5, target)
		assert.NoError(t, err)
		assert.Equal(t, 42, v)
	})

	t.Run("nil input returns nil", func(t *testing.T) {
		v, err := TryTo(nil, 0)
		assert.NoError(t, err)
		assert.Nil(t, v)
	})

	t.Run("To swallows error", func(t *testing.T) {
		v := To("hello", 0)
		assert.Equal(t, 0, v)
	})

	t.Run("context variant", func(t *testing.T) {
		v, err := TryToContext(context.Background(), "99", 0)
		assert.NoError(t, err)
		assert.Equal(t, 99, v)
	})
}

func TestReflectTryToType(t *testing.T) {
	t.Run("string to int via reflect", func(t *testing.T) {
		v, err := ReflectTryToType(reflect.ValueOf("7"), reflect.TypeOf(0), false)
		assert.NoError(t, err)
		assert.Equal(t, 7, v)
	})

	t.Run("ReflectToType returns nil for slice target on error", func(t *testing.T) {
		// A chan cannot be converted to []int — expect nil result
		v := ReflectToType(reflect.ValueOf(make(chan int)), reflect.TypeOf([]int{}))
		assert.Nil(t, v)
	})

	t.Run("ReflectToTypeContext returns nil for slice target on error", func(t *testing.T) {
		v := ReflectToTypeContext(context.Background(), reflect.ValueOf(make(chan int)), reflect.TypeOf([]int{}))
		assert.Nil(t, v)
	})
}

func TestTryCastValue(t *testing.T) {
	t.Run("non-recursive int", func(t *testing.T) {
		v, err := TryCastValue[int]("55", false)
		assert.NoError(t, err)
		assert.Equal(t, 55, v)
	})

	t.Run("context variant", func(t *testing.T) {
		v, err := TryCastValueContext[string](context.Background(), 123, false)
		assert.NoError(t, err)
		assert.Equal(t, "123", v)
	})
}

func TestTryCastRecursive(t *testing.T) {
	type Inner struct{ X int }
	src := map[string]any{"X": 7}

	t.Run("struct from map recursive", func(t *testing.T) {
		v, err := TryCastRecursive[Inner](src, "json")
		assert.NoError(t, err)
		assert.Equal(t, 7, v.X)
	})

	t.Run("context variant", func(t *testing.T) {
		v, err := TryCastRecursiveContext[Inner](context.Background(), src)
		assert.NoError(t, err)
		assert.Equal(t, 7, v.X)
	})
}
