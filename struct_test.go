package gocast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructGetSetFieldValue(t *testing.T) {
	st := &struct {
		Name  string
		Value int64
	}{}

	assert.NoError(t, SetStructFieldValue(&st, "Name", "TestName"), "set Name field value")
	assert.NoError(t, SetStructFieldValue(&st, "Value", int64(127)), "set Value field value")

	name, err := StructFieldValue(st, "Name")
	assert.NoError(t, err, "get Name value")
	assert.Equal(t, "TestName", name)

	value, err := StructFieldValue(st, "Value")
	assert.NoError(t, err, "get Value value")
	assert.Equal(t, int64(127), value)
}
