package gocast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogic(t *testing.T) {
	t.Run("Or", func(t *testing.T) {
		assert.Equal(t, 1, Or(1, 2, 3))
		assert.Equal(t, 2, Or[any](nil, 2, 3))
		assert.Equal(t, 3, Or[any](nil, nil, 3))
		assert.Equal(t, "a", Or[any](nil, nil, "a"))
		assert.Equal(t, "a", Or[any](nil, "a", nil))
		assert.Equal(t, "a", Or[any]("a", nil, nil))
		assert.Equal(t, "a", Or[any]("a", "b", nil))
		assert.Equal(t, "a", Or[any]("a", nil, "b"))
		assert.Equal(t, "a", Or[any](nil, "a", "b"))

		assert.Equal(t, 2, Or[any](0, 2, 3))
		assert.Equal(t, 3, Or[any](0, 0, 3))
		assert.Equal(t, "a", Or[any]("", 0, "a"))
		assert.Equal(t, "a", Or[any]("", "a", 0))
		assert.Equal(t, "a", Or[any]("a", 0, nil))

		assert.Equal(t, 2, Or(0, 2, 3))
		assert.Equal(t, "a", Or("", "a", "c"))
		assert.Equal(t, "a", Or("", "a", "c"))
		assert.Equal(t, "a", Or("a", "", "c"))
	})

	t.Run("IfThen", func(t *testing.T) {
		assert.Equal(t, 1, IfThen(true, 1, 2))
		assert.Equal(t, 2, IfThen(false, 1, 2))
	})

	t.Run("PtrAsValue", func(t *testing.T) {
		assert.Equal(t, 1, PtrAsValue(&[]int{1}[0], 2))
		assert.Equal(t, 2, PtrAsValue((*int)(nil), 2))
	})
}
