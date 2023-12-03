package gocast

import (
	"context"
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

	var target5 = map[string]float64{}
	err = TryMapCopy(target5, struct{ Foo float64 }{Foo: 1.0}, true)
	assert.NoError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(map[string]float64{"Foo": 1.0}, target5))
	err = TryMapCopy(target5, map[string]any{"Foo": 1.0}, true)
	assert.NoError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(map[string]float64{"Foo": 1.0}, target5))

	target6 := MapContext[string, float32](context.TODO(), struct{ Foo float32 }{Foo: 1.0})
	assert.Equal(t, true, reflect.DeepEqual(map[string]float32{"Foo": 1.0}, target6))

	target7 := MapRecursiveContext[string, map[string]any](context.TODO(), struct{ Foo struct{ Bar string } }{Foo: struct{ Bar string }{Bar: "boo"}})
	assert.Equal(t, true, reflect.DeepEqual(map[string]map[string]any{"Foo": {"Bar": "boo"}}, target7))

	// Nil check for map values
	var nilMap = map[string]any{"default": nil, "sub1": []any{nil}, "sub2": map[string]any{"n1": nil, "n2": nil, "n3": nil}}
	var target8 = map[string]any{}
	err = ToMap(&target8, nilMap, true)
	assert.NoError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(nilMap, target8))
}

func TestIsMap(t *testing.T) {
	tests := []struct {
		src any
		trg bool
	}{
		{src: map[string]string{"1": "2"}, trg: true},
		{src: map[string]int{"1": 2}, trg: true},
		{src: map[string]float64{"1": 2.0}, trg: true},
		{src: map[string]float32{"1": 2.0}, trg: true},
		{src: map[any]uint{"1": 2}, trg: true},
		{src: map[any]uint8{"1": 2}, trg: true},
		{src: map[any]uint16{"1": 2}, trg: true},
		{src: 1, trg: false},
		{src: "1", trg: false},
		{src: []string{"1"}, trg: false},
		{src: []int{1}, trg: false},
		{src: []float64{1.0}, trg: false},
		{src: []float32{1.0}, trg: false},
		{src: []uint{1}, trg: false},
		{src: []uint8{1}, trg: false},
	}

	for _, test := range tests {
		assert.Equal(t, test.trg, IsMap(test.src))
	}
}
