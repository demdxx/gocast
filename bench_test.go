package gocast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func bench1(v any) int {
	switch iv := v.(type) {
	case int:
		return int(iv)
	case int8:
		return int(iv)
	case int16:
		return int(iv)
	case int32:
		return int(iv)
	}
	return 0
}

func bench2[T Numeric](v T) int {
	switch iv := any(v).(type) {
	case int:
		return int(iv)
	case int8:
		return int(iv)
	case int16:
		return int(iv)
	case int32:
		return int(iv)
	}
	return 0
}

func BenchmarkApproachTest(b *testing.B) {
	var vals = []struct {
		source int32
		target int
	}{
		{source: 1, target: 1},
		{source: 10, target: 10},
		{source: 100, target: 100},
	}
	b.ReportAllocs()
	b.Run("bench1", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				for _, v := range vals {
					assert.Equal(b, v.target, bench1(v.source))
				}
			}
		})
	})
	b.Run("bench2", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				for _, v := range vals {
					assert.Equal(b, v.target, bench2(v.source))
				}
			}
		})
	})
}
