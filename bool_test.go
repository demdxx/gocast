package gocast

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	tests := []struct {
		name     string
		src      any
		expected bool
	}{
		// nil
		{"nil", nil, false},

		// string cases
		{"string_1", "1", true},
		{"string_T", "T", true},
		{"string_t", "t", true},
		{"string_true", "true", true},
		{"string_TRUE", "TRUE", true},
		{"string_True", "True", true},
		{"string_empty", "", false},
		{"string_false", "false", false},
		{"string_0", "0", false},
		{"string_f", "f", false},
		{"string_F", "F", false},
		{"string_other", "other", false},

		// []byte cases
		{"bytes_1", []byte("1"), true},
		{"bytes_T", []byte("T"), true},
		{"bytes_t", []byte("t"), true},
		{"bytes_true", []byte("true"), true},
		{"bytes_TRUE", []byte("TRUE"), true},
		{"bytes_True", []byte("True"), true},
		{"bytes_empty", []byte(""), false},
		{"bytes_false", []byte("false"), false},
		{"bytes_0", []byte("0"), false},
		{"bytes_f", []byte("f"), false},
		{"bytes_F", []byte("F"), false},
		{"bytes_other", []byte("other"), false},

		// bool cases
		{"bool_true", true, true},
		{"bool_false", false, false},

		// int cases
		{"int_positive", 120, true},
		{"int_negative", -120, true},
		{"int_zero", 0, false},

		// int8 cases
		{"int8_positive", int8(120), true},
		{"int8_negative", int8(-120), true},
		{"int8_zero", int8(0), false},

		// int16 cases
		{"int16_positive", int16(120), true},
		{"int16_negative", int16(-120), true},
		{"int16_zero", int16(0), false},

		// int32 cases
		{"int32_positive", int32(120), true},
		{"int32_negative", int32(-120), true},
		{"int32_zero", int32(0), false},

		// int64 cases
		{"int64_positive", int64(120), true},
		{"int64_negative", int64(-120), true},
		{"int64_zero", int64(0), false},

		// uint cases
		{"uint_positive", uint(120), true},
		{"uint_zero", uint(0), false},

		// uint8 cases
		{"uint8_positive", uint8(120), true},
		{"uint8_zero", uint8(0), false},

		// uint16 cases
		{"uint16_positive", uint16(120), true},
		{"uint16_zero", uint16(0), false},

		// uint32 cases
		{"uint32_positive", uint32(120), true},
		{"uint32_zero", uint32(0), false},

		// uint64 cases
		{"uint64_positive", uint64(120), true},
		{"uint64_zero", uint64(0), false},

		// uintptr cases
		{"uintptr_positive", uintptr(120), true},
		{"uintptr_zero", uintptr(0), false},

		// float32 cases
		{"float32_positive", float32(1.5), true},
		{"float32_negative", float32(-1.5), true},
		{"float32_zero", float32(0.0), false},

		// float64 cases
		{"float64_positive", float64(1.5), true},
		{"float64_negative", float64(-1.5), true},
		{"float64_zero", float64(0.0), false},

		// slice cases
		{"slice_int_nonempty", []int{1, 2, 3}, true},
		{"slice_int_empty", []int{}, false},
		{"slice_int8_nonempty", []int8{1, 2}, true},
		{"slice_int8_empty", []int8{}, false},
		{"slice_int16_nonempty", []int16{1, 2}, true},
		{"slice_int16_empty", []int16{}, false},
		{"slice_int32_nonempty", []int32{1, 2}, true},
		{"slice_int32_empty", []int32{}, false},
		{"slice_int64_nonempty", []int64{1, 2}, true},
		{"slice_int64_empty", []int64{}, false},
		{"slice_uint_nonempty", []uint{1, 2}, true},
		{"slice_uint_empty", []uint{}, false},
		{"slice_uint16_nonempty", []uint16{1, 2}, true},
		{"slice_uint16_empty", []uint16{}, false},
		{"slice_uint32_nonempty", []uint32{1, 2}, true},
		{"slice_uint32_empty", []uint32{}, false},
		{"slice_uint64_nonempty", []uint64{1, 2}, true},
		{"slice_uint64_empty", []uint64{}, false},
		{"slice_float32_nonempty", []float32{1.0, 2.0}, true},
		{"slice_float32_empty", []float32{}, false},
		{"slice_float64_nonempty", []float64{1.0, 2.0}, true},
		{"slice_float64_empty", []float64{}, false},
		{"slice_bool_nonempty", []bool{true, false}, true},
		{"slice_bool_empty", []bool{}, false},
		{"slice_any_nonempty", []any{1, "test"}, true},
		{"slice_any_empty", []any{}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Bool(test.src)
			assert.Equal(t, test.expected, result, "Bool(%v) should return %v", test.src, test.expected)
		})
	}
}

func TestReflectToBool(t *testing.T) {
	tests := []struct {
		name     string
		src      any
		expected bool
	}{
		// Invalid reflect.Value
		{"invalid", nil, false}, // будет создан невалидный reflect.Value

		// string cases
		{"string_1", "1", true},
		{"string_true", "true", true},
		{"string_t", "t", true},
		{"string_empty", "", false},
		{"string_false", "false", false},
		{"string_other", "other", false},

		// []byte cases (через slice reflection)
		{"bytes_1", []byte("1"), true},
		{"bytes_true", []byte("true"), true},
		{"bytes_t", []byte("t"), true},
		{"bytes_empty", []byte(""), false},
		{"bytes_false", []byte("false"), false},

		// slice cases
		{"slice_nonempty", []int{1, 2, 3}, true},
		{"slice_empty", []int{}, false},
		{"slice_string_nonempty", []string{"a", "b"}, true},
		{"slice_string_empty", []string{}, false},

		// array cases
		{"array_nonempty", [3]int{1, 2, 3}, true},
		{"array_empty", [0]int{}, false},

		// map cases
		{"map_nonempty", map[string]int{"a": 1}, true},
		{"map_empty", map[string]int{}, false},

		// bool cases
		{"bool_true", true, true},
		{"bool_false", false, false},

		// int cases
		{"int_positive", 120, true},
		{"int_negative", -120, true},
		{"int_zero", 0, false},

		// int8 cases
		{"int8_positive", int8(120), true},
		{"int8_zero", int8(0), false},

		// int16 cases
		{"int16_positive", int16(120), true},
		{"int16_zero", int16(0), false},

		// int32 cases
		{"int32_positive", int32(120), true},
		{"int32_zero", int32(0), false},

		// int64 cases
		{"int64_positive", int64(120), true},
		{"int64_zero", int64(0), false},

		// uint cases
		{"uint_positive", uint(120), true},
		{"uint_zero", uint(0), false},

		// uint8 cases
		{"uint8_positive", uint8(120), true},
		{"uint8_zero", uint8(0), false},

		// uint16 cases
		{"uint16_positive", uint16(120), true},
		{"uint16_zero", uint16(0), false},

		// uint32 cases
		{"uint32_positive", uint32(120), true},
		{"uint32_zero", uint32(0), false},

		// uint64 cases
		{"uint64_positive", uint64(120), true},
		{"uint64_zero", uint64(0), false},

		// uintptr cases
		{"uintptr_positive", uintptr(120), true},
		{"uintptr_zero", uintptr(0), false},

		// float32 cases
		{"float32_positive", float32(1.5), true},
		{"float32_zero", float32(0.0), false},

		// float64 cases
		{"float64_positive", float64(1.5), true},
		{"float64_zero", float64(0.0), false},

		// other types (should return false)
		{"complex64", complex64(1 + 2i), false},
		{"complex128", complex128(1 + 2i), false},
		{"chan", make(chan int), false},
		{"func", func() {}, false},
		{"ptr", new(int), false},
		{"struct", struct{ X int }{X: 1}, false},
		{"interface", interface{}(nil), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var v reflect.Value
			if test.name == "invalid" {
				v = reflect.Value{} // создаем невалидный reflect.Value
			} else {
				v = reflect.ValueOf(test.src)
			}

			result := ReflectToBool(v)
			assert.Equal(t, test.expected, result, "ReflectToBool(%v) should return %v", test.src, test.expected)
		})
	}
}

func TestReflectToBoolSpecialCases(t *testing.T) {
	// Тест для []byte через reflection (особый случай в коде)
	t.Run("bytes_reflection_special", func(t *testing.T) {
		tests := []struct {
			name     string
			bytes    []byte
			expected bool
		}{
			{"bytes_1", []byte("1"), true},
			{"bytes_true", []byte("true"), true},
			{"bytes_t", []byte("t"), true},
			{"bytes_empty", []byte(""), false},
			{"bytes_false", []byte("false"), false},
			{"bytes_other", []byte("other"), false},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				v := reflect.ValueOf(test.bytes)
				result := ReflectToBool(v)
				assert.Equal(t, test.expected, result)
			})
		}
	})

	// Тест для slice не-bytes типов
	t.Run("non_bytes_slices", func(t *testing.T) {
		tests := []struct {
			name     string
			slice    any
			expected bool
		}{
			{"slice_string_nonempty", []string{"a", "b"}, true},
			{"slice_string_empty", []string{}, false},
			{"slice_interface_nonempty", []interface{}{1, "a"}, true},
			{"slice_interface_empty", []interface{}{}, false},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				v := reflect.ValueOf(test.slice)
				result := ReflectToBool(v)
				assert.Equal(t, test.expected, result)
			})
		}
	})
}

func TestBoolEdgeCases(t *testing.T) {
	// Тестируем граничные случаи

	t.Run("complex_types_fallback_to_reflection", func(t *testing.T) {
		// Типы, которые не обрабатываются напрямую в Bool и попадают в ReflectToBool
		tests := []struct {
			name     string
			src      any
			expected bool
		}{
			{"struct", struct{ X int }{X: 1}, false},
			{"ptr_non_nil", new(int), false},
			{"chan", make(chan int), false},
			{"func", func() {}, false},
			{"complex64", complex64(1 + 2i), false},
			{"complex128", complex128(1 + 2i), false},
			{"map_string_int_nonempty", map[string]int{"a": 1}, true},
			{"map_string_int_empty", map[string]int{}, false},
			{"array", [3]int{1, 2, 3}, true},
			{"array_empty", [0]int{}, false},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				result := Bool(test.src)
				assert.Equal(t, test.expected, result)
			})
		}
	})

	t.Run("string_case_sensitivity", func(t *testing.T) {
		// Проверяем case-insensitive сравнение для "true"
		cases := []struct {
			str      string
			expected bool
		}{
			{"true", true},
			{"TRUE", true},
			{"True", true},
			{"TrUe", true},
			{"tRuE", true},
			{"false", false},
			{"FALSE", false},
			{"False", false},
		}

		for _, test := range cases {
			t.Run("string_"+test.str, func(t *testing.T) {
				result := Bool(test.str)
				assert.Equal(t, test.expected, result)
			})
		}
	})

	t.Run("bytes_case_sensitivity", func(t *testing.T) {
		// Проверяем case-insensitive сравнение для []byte("true")
		cases := []struct {
			bytes    []byte
			expected bool
		}{
			{[]byte("true"), true},
			{[]byte("TRUE"), true},
			{[]byte("True"), true},
			{[]byte("TrUe"), true},
			{[]byte("tRuE"), true},
			{[]byte("false"), false},
			{[]byte("FALSE"), false},
			{[]byte("False"), false},
		}

		for _, test := range cases {
			t.Run("bytes_"+string(test.bytes), func(t *testing.T) {
				result := Bool(test.bytes)
				assert.Equal(t, test.expected, result)
			})
		}
	})
}

func BenchmarkBool(b *testing.B) {
	values := []any{120, uint64(122), "f", "true", "", []byte("t"), true, false, 0.}
	for n := 0; n < b.N; n++ {
		_ = Bool(values[n%len(values)])
	}
}

func BenchmarkToBoolByReflect(b *testing.B) {
	var (
		baseValues = []any{120, uint64(122), "f", "true", "", []byte("t"), true, false, 0.}
		values     = []reflect.Value{}
	)
	for _, v := range baseValues {
		values = append(values, reflect.ValueOf(v))
	}
	for n := 0; n < b.N; n++ {
		_ = ReflectToBool(values[n%len(values)])
	}
}

func BenchmarkBoolTypes(b *testing.B) {
	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Bool("true")
		}
	})

	b.Run("int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Bool(42)
		}
	})

	b.Run("bool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Bool(true)
		}
	})

	b.Run("bytes", func(b *testing.B) {
		bytes := []byte("true")
		for i := 0; i < b.N; i++ {
			_ = Bool(bytes)
		}
	})

	b.Run("slice", func(b *testing.B) {
		slice := []int{1, 2, 3}
		for i := 0; i < b.N; i++ {
			_ = Bool(slice)
		}
	})

	b.Run("reflection_fallback", func(b *testing.B) {
		m := map[string]int{"a": 1}
		for i := 0; i < b.N; i++ {
			_ = Bool(m)
		}
	})
}

// Дополнительные тесты для 100% покрытия
func TestBoolAdditionalCoverage(t *testing.T) {
	t.Run("missing_slice_uint8", func(t *testing.T) {
		// В Go []uint8 и []byte это одно и то же, поэтому []uint8{1, 2} будет обрабатываться
		// как []byte в специальном случае ReflectToBool и возвращать false,
		// поскольку {1, 2} не равно "1", "true" или "t"
		tests := []struct {
			name     string
			src      any
			expected bool
		}{
			{"slice_uint8_nonempty", []uint8{1, 2}, false}, // Не равно "1", "true", "t"
			{"slice_uint8_empty", []uint8{}, false},        // Пустой
			{"slice_uint8_true", []uint8("true"), true},    // Равно "true"
			{"slice_uint8_1", []uint8("1"), true},          // Равно "1"
			{"slice_uint8_t", []uint8("t"), true},          // Равно "t"
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				result := Bool(test.src)
				assert.Equal(t, test.expected, result)
			})
		}
	})

	t.Run("special_reflect_cases", func(t *testing.T) {
		// Тестируем особые случаи для ReflectToBool
		tests := []struct {
			name     string
			create   func() reflect.Value
			expected bool
		}{
			{
				"invalid_value",
				func() reflect.Value { return reflect.Value{} },
				false,
			},
			{
				"nil_pointer",
				func() reflect.Value {
					var p *int
					return reflect.ValueOf(p)
				},
				false,
			},
			{
				"interface_with_nil",
				func() reflect.Value {
					var i interface{} = nil
					return reflect.ValueOf(i)
				},
				false,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				v := test.create()
				result := ReflectToBool(v)
				assert.Equal(t, test.expected, result)
			})
		}
	})

	t.Run("bytes_edge_cases", func(t *testing.T) {
		// Дополнительные тесты для []byte случаев
		tests := []struct {
			name     string
			bytes    []byte
			expected bool
		}{
			{"bytes_nil", nil, false},                     // nil []byte
			{"bytes_single_space", []byte(" "), false},    // непустой, но не равен истинным значениям
			{"bytes_mixed_case", []byte("TrUe"), true},    // mixed case - должно быть true через bytes.EqualFold
			{"bytes_case_sensitive_t", []byte("T"), true}, // заглавная T
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				result := Bool(test.bytes)
				assert.Equal(t, test.expected, result)
			})
		}
	})

	t.Run("reflect_string_edge_cases", func(t *testing.T) {
		// Тесты для строковых случаев через рефлексию
		tests := []struct {
			name     string
			str      string
			expected bool
		}{
			{"reflect_string_T", "T", false}, // В ReflectToBool нет "T", только "t"
			{"reflect_string_1", "1", true},
			{"reflect_string_true", "true", true},
			{"reflect_string_t", "t", true},
			{"reflect_string_false", "false", false},
			{"reflect_string_empty", "", false},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				v := reflect.ValueOf(test.str)
				result := ReflectToBool(v)
				assert.Equal(t, test.expected, result)
			})
		}
	})

	t.Run("special_numeric_cases", func(t *testing.T) {
		// Тесты для особых числовых случаев
		tests := []struct {
			name     string
			src      any
			expected bool
		}{
			{"float32_zero", float32(0.0), false},
			{"float64_zero", float64(0.0), false},
			{"float32_tiny_positive", float32(1e-38), true},
			{"float64_tiny_positive", float64(1e-308), true},
			{"int_min", int(-2147483648), true},
			{"int_max", int(2147483647), true},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				result := Bool(test.src)
				assert.Equal(t, test.expected, result)
			})
		}
	})
}
