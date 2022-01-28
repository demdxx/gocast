GoCast
======
[![GoDoc](https://godoc.org/github.com/demdxx/gocast?status.svg)](https://godoc.org/github.com/demdxx/gocast)
[![Build Status](https://github.com/demdxx/gocast/workflows/run%20tests/badge.svg)](https://github.com/demdxx/gocast/actions?workflow=run%20tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/demdxx/gocast)](https://goreportcard.com/report/github.com/demdxx/gocast)
[![Coverage Status](https://coveralls.io/repos/github/demdxx/gocast/badge.svg)](https://coveralls.io/github/demdxx/gocast)

Easily convert basic types into any other basic types.

## What is GoCast?

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

gocast.ToFloat32("2.12")         // 2.12
gocast.ToFloat64("2.")           // 2.0
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

### Structures

```go
type User struct {
  ID    uint64
  Email string
}

var user User

// Set structure values
err := gocast.SetStructFieldValue(&user, "ID", uint64(19))
err := gocast.SetStructFieldValue(&user, "Email", "iamawesome@mail.com")

id, err := gocast.StructFieldValue(user, "ID")
email, err := gocast.StructFieldValue(user, "Email")
fmt.Printf("User: %d - %s", id, email)
// > User: 19 - iamawesome@mail.com
```

## Benchmarks

```sh
> go test -benchmem -v -race -bench=.

BenchmarkToBool
BenchmarkToBool-8               34045201                44.23 ns/op            0 B/op          0 allocs/op
BenchmarkToBoolByReflect
BenchmarkToBoolByReflect-8      19063656                62.12 ns/op            0 B/op          0 allocs/op
BenchmarkToFloat
BenchmarkToFloat-8              17534598                69.61 ns/op            2 B/op          0 allocs/op
BenchmarkToInt
BenchmarkToInt-8                17316328                68.66 ns/op            2 B/op          0 allocs/op
BenchmarkToUint
BenchmarkToUint-8               17812291                66.71 ns/op            2 B/op          0 allocs/op
BenchmarkToStringByReflect
BenchmarkToStringByReflect-8      485696              2108 ns/op               6 B/op          0 allocs/op
BenchmarkToString
BenchmarkToString-8               532069              1950 ns/op               6 B/op          0 allocs/op
BenchmarkParseTime
BenchmarkParseTime-8              695787              1671 ns/op             464 B/op          5 allocs/op
BenchmarkIsEmpty
BenchmarkIsEmpty-8              24910284                46.36 ns/op            0 B/op          0 allocs/op
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

