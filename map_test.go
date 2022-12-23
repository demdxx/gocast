package gocast

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	assert.Equal(t, true, reflect.DeepEqual(map[string]string{"1": "2"},
		Map[string, string](map[int]int{1: 2})))
	assert.Equal(t, true, reflect.DeepEqual(map[string]string{"1": "2"},
		Map[string, string](map[any]any{1: 2})))
	assert.Equal(t, true, reflect.DeepEqual(map[float64]float32{1: 2},
		Map[float64, float32](map[int16]string{1: "2.0"})))
	assert.Equal(t, true, reflect.DeepEqual(map[string]string{"Foo": "boo"},
		MapRecursive[string, string](struct{ Foo string }{Foo: "boo"})))
}
