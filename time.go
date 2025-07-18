// Copyright (c) 2014 – 2016 Dmitry Ponomarev <demdxx@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package gocast

import (
	"time"
)

var timeFormats = []string{
	time.RFC1123Z,
	time.RFC3339Nano,
	time.UnixDate,
	time.RubyDate,
	time.RFC1123,
	time.RFC3339,
	time.RFC822,
	time.RFC850,
	time.RFC822Z,
	"2006-01-02",
	"2006-01-02 15:04:05",
	"2006/01/02",
	"2006/01/02 15:04:05",
}

// ParseTime from string
func ParseTime(tm string, tmFmt ...string) (t time.Time, err error) {
	if len(tmFmt) == 0 {
		tmFmt = timeFormats
	}
	for _, f := range tmFmt {
		if t, err = time.Parse(f, tm); err == nil {
			break
		}
	}
	return t, err
}
