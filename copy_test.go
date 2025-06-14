package gocast

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testCopy[T any](t *testing.T, src T, expectErr error) T {
	dst, err := TryCopy(src)
	if expectErr != nil {
		assert.ErrorContains(t, err, expectErr.Error(), "error message mismatch for %T", src)
	} else if assert.NoError(t, err, "unexpected error for %T", src) {
		if IsNil(dst) && !IsNil(src) {
			t.Errorf("expected non-nil destination for %T, got nil", src)
		}
		if !assert.True(t, IsNil(dst) == IsNil(src), "expected destination nil state to match source for %T", src) {
			return dst
		}
		switch v := any(src).(type) {
		case nil:
			assert.Nil(t, dst, "expected nil destination for nil source")
		case []int, []string:
			assert.ElementsMatch(t, v, dst, "expected destination to match source for slice of int")
		case map[string]int:
			assert.True(t, reflect.DeepEqual(v, dst), "expected destination to match source for map of string to int")
		default:
			t1 := reflect.TypeOf(v)
			t2 := reflect.TypeOf(dst)
			if assert.True(t, t1 == t2, "expected destination type to match source type for %T", v) {
				if t1.Kind() == reflect.Struct {
					assert.True(t, reflect.DeepEqual(v, dst), "expected destination to match source for struct %T", v)
				} else {
					assert.Equal(t, v, dst, "expected destination type to match source type for %T", v)
				}
			}
		}
	}
	return dst
}

func TestCopy(t *testing.T) {
	t.Run("nil value", func(t *testing.T) { testCopy[any](t, nil, nil) })
	t.Run("int value", func(t *testing.T) { testCopy(t, 42, nil) })
	t.Run("string value", func(t *testing.T) { testCopy(t, "hello", nil) })
	t.Run("slice of int", func(t *testing.T) { testCopy(t, []int{1, 2, 3}, nil) })
	t.Run("map of string to int", func(t *testing.T) { testCopy(t, map[string]int{"a": 1, "b": 2}, nil) })
	t.Run("struct value", func(t *testing.T) { testCopy(t, struct{ Name string }{Name: "test"}, nil) })
	t.Run("unsupported type", func(t *testing.T) {
		testCopy(t, make(chan int), errors.New("unsupported type: chan int"))
	})
	t.Run("deeply nested struct", func(t *testing.T) {
		type Inner struct {
			Value int
			S     string
			Sl    []int
		}
		type Outer struct {
			Inner Inner
		}
		src := Outer{Inner: Inner{
			Value: 42,
			S:     "nested",
			Sl:    []int{1, 2, 3},
		}}
		testCopy(t, src, nil)
	})
	t.Run("nil pointer", func(t *testing.T) {
		var src *struct{ Name string }
		dst := testCopy(t, src, nil)
		assert.Nil(t, dst, "expected nil destination for nil pointer source")
		assert.True(t, reflect.TypeOf(dst) == reflect.TypeOf(src),
			"expected destination type to match source type for nil pointer")
	})
	t.Run("non-nil pointer", func(t *testing.T) {
		src := &struct{ Name string }{Name: "pointer"}
		dst := testCopy(t, src, nil)
		dst.Name = "modified"
		assert.NotEqual(t, src.Name, dst.Name, "expected pointer copy to be independent")
	})
}
