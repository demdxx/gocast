GoCast
======
[![GoDoc](https://godoc.org/github.com/demdxx/gocast?status.svg)](https://godoc.org/github.com/demdxx/gocast)
[![Build Status](https://api.travis-ci.org/demdxx/gocast.svg?branch=master)](https://travis-ci.org/demdxx/gocast)
[![Go Report Card](https://goreportcard.com/badge/github.com/demdxx/gocast)](https://goreportcard.com/report/github.com/demdxx/gocast)

Easily convert basic types into any other basic types.

## What is Cast?

Cast is a library to convert between different GO types in a consistent and easy way.

The library provides a set of methods with the universal casting of types.
All casting methods starts from `To{TargetType}`.

## Usage example

```go
import "github.com/demdxx/gocast"

// Example ToString:
gocast.ToString("strasstr")         // "strasstr"
gocast.ToString(8)                  // "8"
gocast.ToString(8.31)               // "8.31"
gocast.ToString([]byte("one time")) // "one time"
gocast.ToString(nil)                // ""

var foo interface{} = "one more time"
gocast.ToString(foo)                // "one more time"

// Example ToInt:
gocast.ToInt(8)                  // 8
gocast.ToInt(8.31)               // 8
gocast.ToInt("8")                // 8
gocast.ToInt(true)               // 1
gocast.ToInt(false)              // 0

var eight interface{} = 8
gocast.ToInt(eight)              // 8
gocast.ToInt(nil)                // 0
```

```go
func sumAll(vals ...interface{}) int {
  var result int = 0
  for _, v := range vals {
    result += gocast.ToInt(v)
  }
  return result
}
```

License
=======

The MIT License (MIT)

Copyright (c) 2014 Dmitry Ponomarev <demdxx@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

