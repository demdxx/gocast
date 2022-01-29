package gocast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSlice(t *testing.T) {
	type customVal string
	tests := []struct {
		src interface{}
		trg interface{}
		cfn func(v interface{}) interface{}
	}{
		{
			src: []int{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v interface{}) interface{} { return ToStringSlice(v) },
		},
		{
			src: []int32{1, 2, 3, 4},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v interface{}) interface{} { return ToStringSlice(v) },
		},
		{
			src: []customVal{"1", "2", "3", "4"},
			trg: []string{"1", "2", "3", "4"},
			cfn: func(v interface{}) interface{} { return ToStringSlice(v) },
		},
		{
			src: []interface{}{"1", "2.5", 6, 1.2},
			trg: []int{1, 2, 6, 1},
			cfn: func(v interface{}) interface{} { return ToIntSlice(v) },
		},
		{
			src: []float64{1, 2.5, 6, 1.2},
			trg: []int{1, 2, 6, 1},
			cfn: func(v interface{}) interface{} { return ToIntSlice(v) },
		},
		{
			src: []float32{1, 2.5, 6, 1.2},
			trg: []int{1, 2, 6, 1},
			cfn: func(v interface{}) interface{} { return ToIntSlice(v) },
		},
		{
			src: []bool{true, false},
			trg: []int{1, 0},
			cfn: func(v interface{}) interface{} { return ToIntSlice(v) },
		},
		{
			src: []interface{}{"1", "2.5", 6, 1.2},
			trg: []interface{}{"1", "2.5", 6, 1.2},
			cfn: func(v interface{}) interface{} { return ToInterfaceSlice(v) },
		},
		{
			src: []int{1, 2, 3, 4},
			trg: []interface{}{1, 2, 3, 4},
			cfn: func(v interface{}) interface{} { return ToInterfaceSlice(v) },
		},
		{
			src: []int64{1, 2, 3, 4},
			trg: []interface{}{int64(1), int64(2), int64(3), int64(4)},
			cfn: func(v interface{}) interface{} { return ToInterfaceSlice(v) },
		},
		{
			src: []interface{}{"1", "2.5", 6, 1.2, "999.5"},
			trg: []float64{1, 2.5, 6, 1.2, 999.5},
			cfn: func(v interface{}) interface{} { return ToFloat64Slice(v) },
		},
		{
			src: []string{"1", "2.5", "6", "1.2", "999.5"},
			trg: []float64{1, 2.5, 6, 1.2, 999.5},
			cfn: func(v interface{}) interface{} { return ToFloat64Slice(v) },
		},
		{
			src: []interface{}{"1", "2.5", 6, 1.2, "999.5", true},
			trg: []int{1, 2, 6, 1, 999, 1},
			cfn: func(v interface{}) interface{} {
				arr := []int{}
				ToSlice(&arr, v, "")
				return arr
			},
		},
	}
	for _, test := range tests {
		res := test.cfn(test.src)
		assert.ElementsMatch(t, test.trg, res)
	}
}
