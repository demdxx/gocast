package gocast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAny verifies that the deprecated Any wrapper correctly delegates to the
// top-level conversion functions.
func TestAny(t *testing.T) {
	a := &Any{v: "42"}

	assert.Equal(t, "42", a.Str())
	assert.Equal(t, 42, a.Int())
	assert.Equal(t, int64(42), a.Int64())
	assert.Equal(t, uint(42), a.Uint())
	assert.Equal(t, uint64(42), a.Uint64())
	assert.Equal(t, float32(42), a.Float32())
	assert.Equal(t, float64(42), a.Float64())
	assert.Equal(t, "42", a.Any())

	// Bool: only "1", "t", "T", "true" (case-insensitive) return true
	assert.False(t, a.Bool())
	assert.True(t, (&Any{v: "true"}).Bool())
	assert.True(t, (&Any{v: "1"}).Bool())

	t.Run("slice", func(t *testing.T) {
		a2 := &Any{v: []any{1, 2, 3}}
		assert.True(t, a2.IsSlice())
		assert.Equal(t, 3, len(a2.Slice()))
	})

	t.Run("map", func(t *testing.T) {
		a3 := &Any{v: map[string]any{"k": "v"}}
		assert.True(t, a3.IsMap())
		assert.False(t, a3.IsSlice())
	})

	t.Run("nil checks", func(t *testing.T) {
		a4 := &Any{v: nil}
		assert.True(t, a4.IsNil())
		assert.True(t, a4.IsEmpty())
	})

	t.Run("non-nil not empty", func(t *testing.T) {
		a5 := &Any{v: "hello"}
		assert.False(t, a5.IsNil())
		assert.False(t, a5.IsEmpty())
	})
}
