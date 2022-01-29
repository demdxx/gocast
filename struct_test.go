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
	assert.Error(t, SetStructFieldValue(&st, "UndefinedField", int64(127)), "set UndefinedField field value must be error")

	name, err := StructFieldValue(st, "Name")
	assert.NoError(t, err, "get Name value")
	assert.Equal(t, "TestName", name)

	value, err := StructFieldValue(st, "Value")
	assert.NoError(t, err, "get Value value")
	assert.Equal(t, int64(127), value)

	_, err = StructFieldValue(st, "UndefinedField")
	assert.Error(t, err, "get UndefinedField value must be error")
}

func BenchmarkGetSetFieldValue(b *testing.B) {
	st := &struct{ Name string }{}

	b.Run("set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := SetStructFieldValue(st, "Name", "value")
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
