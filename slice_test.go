package gocast

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSliceStruct struct {
	ID   int
	Text string
}

type testSliceCastStruct struct {
	Text string
}

func (it *testSliceCastStruct) CastSet(ctx context.Context, v any) error {
	it.Text = Str(v)
	return nil
}

func TestToSlice(t *testing.T) {
	type customVal string
	tests := []struct {
		src any
		trg any
		err error
		cfn func(v any) (any, error)
	}{
		{
			src: []int{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) (any, error) { return AnySlice[string](v), nil },
		},
		{
			src: []int8{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) (any, error) { return AnySlice[string](v), nil },
		},
		{
			src: []int16{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) (any, error) { return AnySlice[string](v), nil },
		},
		{
			src: []int32{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) (any, error) { return AnySlice[string](v), nil },
		},
		{
			src: []int64{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) (any, error) { return AnySlice[string](v), nil },
		},
		{
			src: []uint{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) (any, error) { return AnySlice[string](v), nil },
		},
		{
			src: []uint8{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) (any, error) { return AnySlice[string](v), nil },
		},
		{
			src: []uint16{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) (any, error) { return AnySlice[string](v), nil },
		},
		{
			src: []uint32{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) (any, error) { return AnySlice[string](v), nil },
		},
		{
			src: []uint64{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) (any, error) { return AnySlice[string](v), nil },
		},
		{
			src: []customVal{"1", "2", "3", "4"},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) (any, error) { return AnySlice[string](v), nil },
		},
		{
			src: []any{"1", "2.5", 6, 1.2},
			trg: []int{1, 2, 6, 1},
			cfn: func(v any) (any, error) { return AnySlice[int](v), nil },
		},
		{
			src: []float64{1, 2.5, 6, 1.2},
			trg: []int{1, 2, 6, 1},
			cfn: func(v any) (any, error) { return AnySlice[int](v), nil },
		},
		{
			src: []float32{1, 2.5, 6, 1.2},
			trg: []int{1, 2, 6, 1},
			cfn: func(v any) (any, error) { return AnySlice[int](v), nil },
		},
		{
			src: []bool{true, false},
			trg: []int{1, 0},
			cfn: func(v any) (any, error) { return AnySlice[int](v), nil },
		},
		{
			src: []any{"1", "2.5", 6, 1.2},
			trg: []any{"1", "2.5", 6, 1.2},
			cfn: func(v any) (any, error) { return AnySlice[any](v), nil },
		},
		{
			src: []int{1, 2, 3, 4},
			trg: []any{1, 2, 3, 4},
			cfn: func(v any) (any, error) { return AnySlice[any](v), nil },
		},
		{
			src: []int64{1, 2, 3, 4},
			trg: []any{int64(1), int64(2), int64(3), int64(4)},
			cfn: func(v any) (any, error) { return AnySlice[any](v), nil },
		},
		{
			src: []any{"1", "2.5", 6, 1.2, "999.5"},
			trg: []float64{1, 2.5, 6, 1.2, 999.5},
			cfn: func(v any) (any, error) { return AnySlice[float64](v), nil },
		},
		{
			src: []string{"1", "2.5", "6", "1.2", "999.5"},
			trg: []float64{1, 2.5, 6, 1.2, 999.5},
			cfn: func(v any) (any, error) { return AnySlice[float64](v), nil },
		},
		{
			src: []any{"1", "2.5", 6, 1.2, "999.5", true},
			trg: []int{1, 2, 6, 1, 999, 1},
			cfn: func(v any) (any, error) { arr := []int{}; err := TryToAnySlice(&arr, v); return arr, err },
		},
		{
			src: []any{"1", "2.5", 6, 1.2, "999.5", true},
			trg: []int{1, 2, 6, 1, 999, 1},
			cfn: func(v any) (any, error) { return AnySlice[int](v), nil },
		},
		{
			src: []any{"1", "2.5", 6, 1.2, "999.5", true},
			trg: []int{1, 2, 6, 1, 999, 1},
			cfn: func(v any) (any, error) { return Slice[int](v.([]any)), nil },
		},
		{
			src: []map[any]any{{"ID": 1, "Text": "text1"}, {"ID": 2, "Text": "text1"}},
			trg: []testSliceStruct{{ID: 1, Text: "text1"}, {ID: 2, Text: "text1"}},
			cfn: func(v any) (any, error) { return Slice[testSliceStruct](v.([]map[any]any)), nil },
		},
		{
			src: []any{"text1", "text2"},
			trg: []testSliceCastStruct{{Text: "text1"}, {Text: "text2"}},
			cfn: func(v any) (any, error) { return Slice[testSliceCastStruct](v.([]any)), nil },
		},
		{
			src: []any{"text1", "text2"},
			trg: []any{"text1", "text2"},
			cfn: func(v any) (any, error) { return Slice[any](v.([]any)), nil },
		},
		{
			src: nil,
			err: ErrInvalidParams,
			cfn: func(v any) (any, error) { return TryAnySlice[int](v) },
		},
		{
			src: 1,
			err: ErrInvalidParams,
			cfn: func(v any) (any, error) { return TryAnySlice[int](v) },
		},
	}
	for _, test := range tests {
		res, err := test.cfn(test.src)
		if test.err != nil {
			assert.ErrorContains(t, err, test.err.Error())
		} else if assert.NoError(t, err) {
			assert.ElementsMatch(t, test.trg, res)
		}
	}
}

func TestIsSlice(t *testing.T) {
	assert.True(t, IsSlice([]int{}))
	assert.True(t, IsSlice([]bool{}))
	assert.True(t, IsSlice([]testSliceCastStruct{}))
	assert.True(t, IsSlice(([]any)(nil)))
	assert.False(t, IsSlice("not a slice"))
}
