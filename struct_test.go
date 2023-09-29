package gocast

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type customInt int

func (c *customInt) CastSet(ctx context.Context, v any) error {
	switch val := v.(type) {
	case customInt:
		*c = val
	default:
		*c = customInt(Number[int](v))
	}
	return nil
}

type testStructParent struct {
	ParentName string `field:"parent_name"`
}

type testStruct struct {
	testStructParent

	Name        string      `field:"name"`
	Value       int64       `field:"value"`
	Count       customInt   `field:"count"`
	Counts      []customInt `field:"counts"`
	AnyTarget   any         `field:"anytarget"`
	NilVal      any         `field:"nilval,omitempty"`
	ignore      bool        `field:"ignore"`
	CreatedAt   time.Time   `field:"created_at"`
	UpdatedAt   time.Time   `field:"updated_at"`
	PublishedAt time.Time   `field:"published_at"`
}

var testStructPreparedValue = map[string]any{
	"parent_name":  "parent name",
	"name":         "test",
	"value":        "1900",
	"count":        112.2,
	"counts":       []any{"1.1", 2.1},
	"created_at":   "2020/10/10",
	"updated_at":   time.Now().Unix(),
	"published_at": time.Now(),
	"anytarget":    "hi",
	"nilval":       nil,
}

func testPreparedStruct(t *testing.T, it *testStruct) {
	assert.Equal(t, "test", it.Name)
	assert.Equal(t, int64(1900), it.Value)
	assert.Equal(t, customInt(112), it.Count)
	assert.ElementsMatch(t, []customInt{1, 2}, it.Counts)
	assert.Equal(t, 2020, it.CreatedAt.Year())
	assert.Equal(t, time.Now().Year(), it.UpdatedAt.Year())
	assert.Equal(t, time.Now().Year(), it.PublishedAt.Year())
	assert.Equal(t, "hi", it.AnyTarget)
	assert.Equal(t, false, it.ignore)
	assert.Nil(t, it.NilVal)
}

func TestStructGetSetFieldValue(t *testing.T) {
	st := &testStruct{}
	ctx := context.TODO()

	assert.NoError(t, SetStructFieldValue(ctx, &st, "Name", "TestName"), "set Name field value")
	assert.NoError(t, SetStructFieldValue(ctx, &st, "Value", int64(127)), "set Value field value")
	assert.NoError(t, SetStructFieldValue(ctx, &st, "Count", 3.1), "set Count field value")
	assert.NoError(t, SetStructFieldValue(ctx, &st, "CreatedAt", "2020/10/10"), "set CreatedAt field value")
	assert.NoError(t, SetStructFieldValue(ctx, &st, "UpdatedAt", time.Now().Unix()), "set UpdatedAt field value")
	assert.NoError(t, SetStructFieldValue(ctx, &st, "PublishedAt", time.Now()), "set PublishedAt field value")
	assert.NoError(t, SetStructFieldValue(ctx, &st, "AnyTarget", "hi"), "set AnyTarget field value")
	assert.NoError(t, SetStructFieldValue(ctx, &st, "ParentName", "parent name"), "set ParentName field value")
	assert.Error(t, SetStructFieldValue(ctx, &st, "UndefinedField", int64(127)), "set UndefinedField field value must be error")

	name, err := StructFieldValue(st, "Name")
	assert.NoError(t, err, "get Name value")
	assert.Equal(t, "TestName", name)

	value, err := StructFieldValue(st, "Value")
	assert.NoError(t, err, "get Value value")
	assert.Equal(t, int64(127), value)

	count, err := StructFieldValue(st, "Count")
	assert.NoError(t, err, "get Count value")
	assert.Equal(t, customInt(3), count)

	createdAt, err := StructFieldValue(st, "CreatedAt")
	assert.NoError(t, err, "get CreatedAt value")
	assert.Equal(t, 2020, createdAt.(time.Time).Year())

	updatedAt, err := StructFieldValue(st, "UpdatedAt")
	assert.NoError(t, err, "get UpdatedAt value")
	assert.Equal(t, time.Now().Year(), updatedAt.(time.Time).Year())

	publishedAt, err := StructFieldValue(st, "PublishedAt")
	assert.NoError(t, err, "get PublishedAt value")
	assert.Equal(t, time.Now().Year(), publishedAt.(time.Time).Year())

	anyTarget, err := StructFieldValue(st, "AnyTarget")
	assert.NoError(t, err, "get AnyTarget value")
	assert.Equal(t, "hi", anyTarget)

	parentName, err := StructFieldValue(st, "ParentName")
	assert.NoError(t, err, "get ParentName value")
	assert.Equal(t, "parent name", parentName)

	_, err = StructFieldValue(st, "UndefinedField")
	assert.Error(t, err, "get UndefinedField value must be error")
}

func TestStructCast(t *testing.T) {
	res, err := Struct[testStruct](testStructPreparedValue, `field`)
	assert.NoError(t, err)
	testPreparedStruct(t, &res)
}

func TestStructCastNested(t *testing.T) {
	testStructPrepared, err := Struct[testStruct](testStructPreparedValue, `field`)
	assert.NoError(t, err)
	testPreparedStruct(t, &testStructPrepared)

	type testStruct2 struct {
		Sub      testStruct             `field:"sub"`
		SubMap   map[string]*testStruct `field:"submap"`
		SubSlice []*testStruct          `field:"subslice"`
	}

	data := map[string]any{
		"sub": testStructPreparedValue,
		"submap": map[string]any{
			"a": testStructPreparedValue,
			"b": testStructPreparedValue,
			"c": &testStructPrepared,
		},
		"subslice": []any{
			testStructPreparedValue,
			testStructPreparedValue,
			&testStructPrepared,
		},
	}

	t.Run("cast", func(t *testing.T) {
		res, err := Struct[testStruct2](data, "field")
		assert.NoError(t, err)
		testPreparedStruct(t, &res.Sub)
		assert.Equal(t, 3, len(res.SubMap))
		for k, val := range res.SubMap {
			t.Run("SubMap["+k+"]", func(t *testing.T) {
				testPreparedStruct(t, val)
			})
		}
		assert.Equal(t, 3, len(res.SubSlice))
		for i, val := range res.SubSlice {
			t.Run("SubSlice["+Str(i)+"]", func(t *testing.T) {
				testPreparedStruct(t, val)
			})
		}
	})

	t.Run("error", func(t *testing.T) {
		_, err := Struct[testStruct2](1, "field")
		assert.ErrorIs(t, err, ErrUnsupportedSourceType)
	})
}

func TestStructFieldNames(t *testing.T) {
	fields := StructFieldNames(testStruct{}, "-")
	assert.ElementsMatch(t,
		[]string{"Name", "Value", "Count", "Counts", "CreatedAt", "UpdatedAt",
			"PublishedAt", "AnyTarget", "NilVal", "ignore", "ParentName"}, fields)
	fields = StructFieldNames(testStruct{}, "")
	assert.ElementsMatch(t,
		[]string{"name", "value", "count", "counts", "created_at", "updated_at",
			"published_at", "anytarget", "nilval", "ignore", "parent_name"}, fields)
}

func TestStructFieldTags(t *testing.T) {
	fields := StructFieldTags(testStruct{}, "field")
	assert.Equal(t,
		map[string]string{
			"ParentName":  "parent_name",
			"AnyTarget":   "anytarget",
			"Count":       "count",
			"Counts":      "counts",
			"Name":        "name",
			"NilVal":      "nilval",
			"CreatedAt":   "created_at",
			"UpdatedAt":   "updated_at",
			"PublishedAt": "published_at",
			"Value":       "value",
			"ignore":      "ignore"}, fields)
}

func TestIsStruct(t *testing.T) {
	assert.True(t, IsStruct(testStruct{}))
	assert.False(t, IsStruct(1))
}

func BenchmarkGetSetFieldValue(b *testing.B) {
	st := &struct{ Name string }{}
	ctx := context.TODO()

	b.Run("set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := SetStructFieldValue(ctx, st, "Name", "value")
			if err != nil {
				b.Error("set field error", err.Error())
			}
		}
	})

	b.Run("get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := StructFieldValue(st, "Name")
			if err != nil {
				b.Error("get field error", err.Error())
			}
		}
	})
}
