package gocast

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointer(t *testing.T) {
	t.Run("PtrAsValue", func(t *testing.T) {
		assert.Equal(t, 1, PtrAsValue(&[]int{1}[0], 2))
		assert.Equal(t, 2, PtrAsValue((*int)(nil), 2))
	})

	t.Run("Ptr", func(t *testing.T) {
		v := 1
		assert.Equal(t, &v, Ptr(1))
		assert.NotNil(t, Ptr(1))
		assert.True(t, *Ptr(1) == 1)
		assert.True(t, reflect.TypeOf(Ptr(1)).Kind() == reflect.Pointer)
	})
}
