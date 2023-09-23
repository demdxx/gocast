package gocast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	assert.Equal(t, 0, Len([]any{}))
	assert.Equal(t, 2, Len([]any{1, "2"}))
	assert.Equal(t, 3, Len([]string{"1", "2", "3"}))
	assert.Equal(t, 3, Len("123"))
	assert.Equal(t, 4, Len([]int{1, 2, 3, 4}))
	assert.Equal(t, 4, Len([]float64{1, 2, 3, 4}))
	assert.Equal(t, 4, Len(&[]float64{1, 2, 3, 4}))
	assert.Equal(t, 3, Len([]byte("123")))
	assert.Equal(t, 0, Len([]int(nil)))

	assert.Equal(t, 0, Len((map[string]any)(nil)))
	assert.Equal(t, 1, Len(map[string]any{"a": 1}))
	assert.Equal(t, 1, Len(map[string]string{"a": "1"}))
	assert.Equal(t, 1, Len(map[any]any{"a": "1"}))
	assert.Equal(t, 1, Len(&map[any]any{"a": "1"}))

	assert.Equal(t, -1, Len(1))
	assert.Equal(t, -1, Len(struct{}{}))
	assert.Equal(t, -1, Len[any](nil))
}
