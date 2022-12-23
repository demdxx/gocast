GoCast
======
[![GoDoc](https://godoc.org/github.com/demdxx/gocast?status.svg)](https://godoc.org/github.com/demdxx/gocast)
[![Build Status](https://github.com/demdxx/gocast/workflows/run%20tests/badge.svg)](https://github.com/demdxx/gocast/actions?workflow=run%20tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/demdxx/gocast)](https://goreportcard.com/report/github.com/demdxx/gocast)
[![Coverage Status](https://coveralls.io/repos/github/demdxx/gocast/badge.svg?branch=master)](https://coveralls.io/github/demdxx/gocast?branch=master)

Easily convert basic types into any other basic types.

## What is GoCast?

Cast is a library to convert between different GO types in a consistent and easy way.

The library provides a set of methods with the universal casting of types.

**The new version** still support old conversion functions like *ToString*, *ToInt*, *ToMap* and etc.
The new version uses insted just the name od the type *Str*, *Number[type]*, *Map[K,V]* and etc.
Also was added new functions which returns the error apearing during type casting,
all such functions starts from word *Try{Name}*. Like *TryStr*, *TryNumber[type]*, *TryMap[K,V]*, *TryCast[type]* and etc.

## Usage example

```go
import "github.com/demdxx/gocast/v2"

// Example string casting:
gocast.Str("strasstr")           // "strasstr"
gocast.Str(8)                    // "8"
gocast.Str(8.31)                 // "8.31"
gocast.Str([]byte("one time"))   // "one time"
gocast.Str(nil)                  // ""

var foo any = "one more time"
gocast.Str(foo)                  // "one more time"

// Example number casting:
gocast.Number[int](8)            // 8
gocast.Number[int](8.31)         // 8
gocast.Number[int]("8")          // 8
gocast.Number[int](true)         // 1
gocast.Number[int](false)        // 0

var eight any = 8
gocast.Number[int](eight)        // 8
gocast.Cast[int](eight)          // 8
gocast.Number[int](nil)          // 0

// Number converts only into numeric values (simpler and faster then Cast)
gocast.Number[float32]("2.12")   // 2.12

// Cast converts any type to any other type
gocast.Cast[float64]("2.")       // 2.0

val, err := gocast.TryCast[int]("123.2") // 123, <nil>
```

```go
func sumAll(vals ...any) int {
  var result int = 0
  for _, v := range vals {
    result += gocast.Number[int](v)
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

goos: darwin
goarch: amd64
pkg: github.com/demdxx/gocast/v2
cpu: VirtualApple @ 2.50GHz
BenchmarkApproachTest
BenchmarkApproachTest/bench1
BenchmarkApproachTest/bench1-8    307263              3611 ns/op          0 B/op           0 allocs/op
BenchmarkApproachTest/bench2
BenchmarkApproachTest/bench2-8    337455              3546 ns/op          0 B/op           0 allocs/op
BenchmarkBool
BenchmarkBool-8                 16757396                75.54 ns/op       0 B/op           0 allocs/op
BenchmarkToBoolByReflect
BenchmarkToBoolByReflect-8      10259976               121.2 ns/op        0 B/op           0 allocs/op
BenchmarkToFloat
BenchmarkToFloat-8               7107516               169.0 ns/op        2 B/op           0 allocs/op
BenchmarkToInt
BenchmarkToInt-8                 7212882               168.5 ns/op        2 B/op           0 allocs/op
BenchmarkToUint
BenchmarkToUint-8                7202331               166.8 ns/op        2 B/op           0 allocs/op
BenchmarkToStringByReflect
BenchmarkToStringByReflect-8      970101              1251 ns/op          6 B/op           0 allocs/op
BenchmarkToString
BenchmarkToString-8               995446              1166 ns/op          6 B/op           0 allocs/op
BenchmarkGetSetFieldValue
BenchmarkGetSetFieldValue/set
BenchmarkGetSetFieldValue/set-8                  1468659               821.3 ns/op        32 B/op          2 allocs/op
BenchmarkGetSetFieldValue/get
BenchmarkGetSetFieldValue/get-8                  1407585               851.8 ns/op        48 B/op          3 allocs/op
BenchmarkParseTime
BenchmarkParseTime-8                              504934              2258 ns/op         464 B/op          5 allocs/op
BenchmarkIsEmpty
BenchmarkIsEmpty-8                              26776375                54.08 ns/op        0 B/op          0 allocs/op
PASS
ok      github.com/demdxx/gocast/v2     18.977s
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

