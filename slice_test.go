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
		cfn func(v any) any
	}{
		{
			src: []int{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) any { return ToStringSlice(v) },
		},
		{
			src: []int8{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) any { return ToStringSlice(v) },
		},
		{
			src: []int16{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) any { return ToStringSlice(v) },
		},
		{
			src: []int32{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) any { return ToStringSlice(v) },
		},
		{
			src: []int64{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) any { return ToStringSlice(v) },
		},
		{
			src: []uint{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) any { return ToStringSlice(v) },
		},
		{
			src: []uint8{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) any { return ToStringSlice(v) },
		},
		{
			src: []uint16{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) any { return ToStringSlice(v) },
		},
		{
			src: []uint32{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) any { return ToStringSlice(v) },
		},
		{
			src: []uint64{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) any { return ToStringSlice(v) },
		},
		{
			src: []customVal{"1", "2", "3", "4"},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v any) any { return ToStringSlice(v) },
		},
		{
			src: []any{"1", "2.5", 6, 1.2},
			trg: []int{1, 2, 6, 1},
			cfn: func(v any) any { return ToIntSlice(v) },
		},
		{
			src: []float64{1, 2.5, 6, 1.2},
			trg: []int{1, 2, 6, 1},
			cfn: func(v any) any { return ToIntSlice(v) },
		},
		{
			src: []float32{1, 2.5, 6, 1.2},
			trg: []int{1, 2, 6, 1},
			cfn: func(v any) any { return ToIntSlice(v) },
		},
		{
			src: []bool{true, false},
			trg: []int{1, 0},
			cfn: func(v any) any { return ToIntSlice(v) },
		},
		{
			src: []any{"1", "2.5", 6, 1.2},
			trg: []any{"1", "2.5", 6, 1.2},
			cfn: func(v any) any { return ToInterfaceSlice(v) },
		},
		{
			src: []int{1, 2, 3, 4},
			trg: []any{1, 2, 3, 4},
			cfn: func(v any) any { return ToInterfaceSlice(v) },
		},
		{
			src: []int64{1, 2, 3, 4},
			trg: []any{int64(1), int64(2), int64(3), int64(4)},
			cfn: func(v any) any { return ToInterfaceSlice(v) },
		},
		{
			src: []any{"1", "2.5", 6, 1.2, "999.5"},
			trg: []float64{1, 2.5, 6, 1.2, 999.5},
			cfn: func(v any) any { return ToFloat64Slice(v) },
		},
		{
			src: []string{"1", "2.5", "6", "1.2", "999.5"},
			trg: []float64{1, 2.5, 6, 1.2, 999.5},
			cfn: func(v any) any { return ToFloat64Slice(v) },
		},
		{
			src: []any{"1", "2.5", 6, 1.2, "999.5", true},
			trg: []int{1, 2, 6, 1, 999, 1},
			cfn: func(v any) any {
				arr := []int{}
				_ = ToSlice(&arr, v, "")
				return arr
			},
		},
		{
			src: []any{"1", "2.5", 6, 1.2, "999.5", true},
			trg: []int{1, 2, 6, 1, 999, 1},
			cfn: func(v any) any { return Slice[int](v.([]any)) },
		},
		{
			src: []map[any]any{{"ID": 1, "Text": "text1"}, {"ID": 2, "Text": "text1"}},
			trg: []testSliceStruct{{ID: 1, Text: "text1"}, {ID: 2, Text: "text1"}},
			cfn: func(v any) any { return Slice[testSliceStruct](v.([]map[any]any)) },
		},
		{
			src: []any{"text1", "text2"},
			trg: []testSliceCastStruct{{Text: "text1"}, {Text: "text2"}},
			cfn: func(v any) any { return Slice[testSliceCastStruct](v.([]any)) },
		},
	}
	for _, test := range tests {
		res := test.cfn(test.src)
		assert.ElementsMatch(t, test.trg, res)
	}
}
