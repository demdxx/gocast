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

	var target1 = map[any]any{}
	err := ToMap(&target1, struct{ Foo string }{Foo: "boo"}, true)
	assert.NoError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(map[any]any{"Foo": "boo"}, target1))
	err = ToMap(&target1, map[string]any{"Foo": "boo"}, true)
	assert.NoError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(map[any]any{"Foo": "boo"}, target1))

	var target2 = map[string]any{}
	err = ToMap(&target2, struct{ Foo string }{Foo: "boo"}, true)
	assert.NoError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(map[string]any{"Foo": "boo"}, target2))
	err = ToMap(&target2, map[string]any{"Foo": "boo"}, true)
	assert.NoError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(map[string]any{"Foo": "boo"}, target2))

	var target3 = map[string]string{}
	err = ToMap(&target3, struct{ Foo string }{Foo: "boo"}, true)
	assert.NoError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(map[string]string{"Foo": "boo"}, target3))
	err = ToMap(&target3, map[string]any{"Foo": "boo"}, true)
	assert.NoError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(map[string]string{"Foo": "boo"}, target3))

	var target4 = map[string]int{}
	err = ToMap(&target4, struct{ Foo int }{Foo: 1}, true)
	assert.NoError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(map[string]int{"Foo": 1}, target4))
	err = ToMap(&target4, map[string]any{"Foo": 1}, true)
	assert.NoError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(map[string]int{"Foo": 1}, target4))
}
