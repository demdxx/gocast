package gocast

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	var tests = []struct {
		src    string
		target time.Time
	}{
		{src: "2021/10/24", target: time.Date(2021, 10, 24, 0, 0, 0, 0, time.UTC)},
		{src: "2021-10-24", target: time.Date(2021, 10, 24, 0, 0, 0, 0, time.UTC)},
		{src: "Sun, 24 Oct 2021 00:00:00 UTC", target: time.Date(2021, 10, 24, 0, 0, 0, 0, time.UTC)},
		{src: "Sun Oct 24 00:00:00 UTC 2021", target: time.Date(2021, 10, 24, 0, 0, 0, 0, time.UTC)},
		{src: "Sun, 24 Oct 2021 00:00:00 UTC", target: time.Date(2021, 10, 24, 0, 0, 0, 0, time.UTC)},
	}

	for _, test := range tests {
		v, err := ParseTime(test.src)
		if err != nil {
			t.Error(err)
		}
		if v != test.target {
			t.Errorf("target must be equal %v != %v", v, test.target)
		}
	}
}

func BenchmarkParseTime(b *testing.B) {
	values := []string{
		"2021/10/24",
		"2021-10-24",
		"Sun, 24 Oct 2021 00:00:00 UTC",
		"Sun Oct 24 00:00:00 UTC 2021",
		"Sun, 24 Oct 2021 00:00:00 UTC"}
	for n := 0; n < b.N; n++ {
		_, _ = ParseTime(values[n%len(values)])
	}
}
