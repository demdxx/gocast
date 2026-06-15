package gocast

import (
	"context"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStructWalkObjectParentAndStruct(t *testing.T) {
	ctx := context.Background()

	type Inner struct{ V int }
	type Outer struct{ Inner Inner }

	src := Outer{Inner: Inner{V: 42}}

	var (
		outerObj StructWalkObject
		innerObj StructWalkObject
	)

	err := StructWalk(ctx, &src, func(ctx context.Context, curObj StructWalkObject, field StructWalkField, path []string) error {
		if len(path) == 0 {
			outerObj = curObj
		} else {
			innerObj = curObj
		}
		return nil
	})
	assert.NoError(t, err)

	// outerObj is the root — it has no parent
	if assert.NotNil(t, outerObj) {
		assert.Nil(t, outerObj.Parent())
		assert.NotNil(t, outerObj.Struct())
	}

	// innerObj is the nested struct — its parent is outerObj
	if assert.NotNil(t, innerObj) {
		assert.Equal(t, outerObj, innerObj.Parent())
		assert.NotNil(t, innerObj.Struct())
	}
}

func TestStructWalkErrWalkSkip(t *testing.T) {
	ctx := context.Background()

	type S struct {
		Skip  string
		Other string
	}

	visited := map[string]bool{}
	src := S{Skip: "s", Other: "o"}

	err := StructWalk(ctx, &src, func(ctx context.Context, curObj StructWalkObject, field StructWalkField, path []string) error {
		visited[field.Name()] = true
		if field.Name() == "Skip" {
			return ErrWalkSkip
		}
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, visited["Skip"])
	assert.True(t, visited["Other"])
}

func TestStructWalkErrWalkStop(t *testing.T) {
	ctx := context.Background()

	type S struct {
		A string
		B string
		C string
	}

	visited := []string{}
	src := S{A: "a", B: "b", C: "c"}

	err := StructWalk(ctx, &src, func(ctx context.Context, curObj StructWalkObject, field StructWalkField, path []string) error {
		visited = append(visited, field.Name())
		if field.Name() == "B" {
			return ErrWalkStop
		}
		return nil
	})
	assert.NoError(t, err) // ErrWalkStop is converted to nil by StructWalk
	assert.Contains(t, visited, "A")
	assert.Contains(t, visited, "B")
	assert.NotContains(t, visited, "C")
}

func TestWalkWithPathExtractor(t *testing.T) {
	ctx := context.Background()

	type Item struct {
		Name  string `custom:"name_field"`
		Value int    `custom:"value_field"`
	}

	src := Item{Name: "gocast", Value: 42}
	collected := map[string]any{}

	extractor := func(ctx context.Context, curObj StructWalkObject, field StructWalkField, path []string) string {
		tag := field.Tag("custom")
		if tag == "" {
			return field.Name()
		}
		return tag
	}

	err := StructWalk(ctx, &src, func(ctx context.Context, curObj StructWalkObject, field StructWalkField, path []string) error {
		name := extractor(ctx, curObj, field, path)
		collected[name] = field.Value()
		return nil
	}, WalkWithPathExtractor(extractor))

	assert.NoError(t, err)
	assert.Equal(t, "gocast", collected["name_field"])
	assert.Equal(t, 42, collected["value_field"])
}

func TestStructWalk(t *testing.T) {
	ctx := context.Background()
	emptyWalker := func(ctx context.Context, curObj StructWalkObject, field StructWalkField, path []string) error {
		return nil
	}
	t.Run("noStruct", func(t *testing.T) {
		err := StructWalk(ctx, 1, emptyWalker)
		assert.ErrorIs(t, err, ErrUnsupportedSourceType)

		err = StructWalk(ctx, map[string]any{}, emptyWalker)
		assert.ErrorIs(t, err, ErrUnsupportedSourceType)

		err = StructWalk(ctx, struct{}{}, emptyWalker)
		assert.NoError(t, err)
	})

	t.Run("init.env", func(t *testing.T) {
		_ = os.Setenv("TEST_V1", "test")
		_ = os.Setenv("TEST_V2", "1")

		var testStruct = struct {
			V1 string `env:"TEST_V1"`
			V2 int    `env:"TEST_V2"`
		}{}

		err := StructWalk(ctx, &testStruct, func(ctx context.Context, curObj StructWalkObject, field StructWalkField, path []string) error {
			assert.True(t, field.IsEmpty())
			err := field.SetValue(ctx, os.Getenv(field.Tag("env")))
			if assert.NoError(t, err, `set value for field "%s"`, field.Name()) {
				if field.Name() == "V1" {
					assert.Equal(t, "TEST_V1", field.Tag("env"))
					assert.Equal(t, "test", field.Value())
				}
				if field.Name() == "V2" {
					assert.Equal(t, "TEST_V2", field.Tag("env"))
					assert.Equal(t, 1, field.Value())
				}
			}
			return err
		})
		assert.NoError(t, err)
	})

	t.Run("init.nested", func(t *testing.T) {
		type (
			N2 struct {
				Text string `field:"text"`
			}
			N1 struct {
				V1 string `field:"v1"`
				V2 int    `field:"v2"`
				N2 N2     `field:"n2"`
			}
			nestedStruct struct {
				T  time.Time `field:"t"`
				V1 string    `field:"v1"`
				V2 int       `field:"v2"`
				N1 N1        `field:"n1"`
			}
		)
		source := map[string]any{
			"t":  "2021-01-01T00:00:00Z",
			"v1": "test",
			"v2": 1,
			"n1": map[string]any{
				"v1": "test",
				"v2": "1",
				"n2": map[string]any{
					"text": "test",
				},
			},
		}
		testStruct := nestedStruct{}
		targetStruct := nestedStruct{
			T:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			V1: "test",
			V2: 1,
			N1: N1{
				V1: "test",
				V2: 1,
				N2: N2{Text: "test"},
			},
		}

		err := StructWalk(ctx, &testStruct, func(ctx context.Context, curObj StructWalkObject, field StructWalkField, path []string) error {
			if field.RefValue().Kind() == reflect.Struct {
				switch field.Value().(type) {
				case time.Time:
				default:
					return nil
				}
			}
			data := source
			for _, p := range path {
				data = data[p].(map[string]any)
			}

			if field.RefValue().Kind() != reflect.Struct {
				assert.True(t, field.IsEmpty())
			}
			err := field.SetValue(ctx, data[field.Tag("field")])
			assert.NoError(t, err, `set value for field "%s.%s"`, strings.Join(path, "."), field.Name())
			return err
		}, WalkWithPathTag("field"))

		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(testStruct, targetStruct), "compare struct: %#v", testStruct)
	})
}
