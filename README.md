# GoCast

[![GoDoc](https://godoc.org/github.com/demdxx/gocast?status.svg)](https://godoc.org/github.com/demdxx/gocast)
[![Build Status](https://github.com/demdxx/gocast/workflows/Tests/badge.svg)](https://github.com/demdxx/gocast/actions?workflow=Tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/demdxx/gocast)](https://goreportcard.com/report/github.com/demdxx/gocast)
[![Coverage Status](https://coveralls.io/repos/github/demdxx/gocast/badge.svg?branch=master)](https://coveralls.io/github/demdxx/gocast?branch=master)

## Introduction

GoCast is a powerful Go library that allows you to easily convert between different basic types in a consistent and efficient way. Whether you need to convert strings, numbers, or other basic types, GoCast has got you covered.

## Features

- **Universal Type Casting**: GoCast provides a set of methods for universal type casting, making it easy to convert data between various types.
- **Struct Field Manipulation**: GoCast allows you to set and retrieve values of struct fields dynamically, making it a valuable tool for working with complex data structures.
- **Custom Type Support**: You can define custom types and implement your own conversion logic, giving you full control over how data is cast.

## Installation

To use GoCast in your Go project, simply import it:

```go
import "github.com/demdxx/gocast/v2"
```

## Usage Example

Here are some examples of how you can use GoCast:

```go
// Example string casting:
gocast.Str("strasstr")           // "strasstr"
gocast.Str(8)                    // "8"
gocast.Str(8.31)                 // "8.31"
gocast.Str([]byte("one time"))   // "one time"
gocast.Str(nil)                  // ""

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

res := gocast.Map[string, any](struct{ID int64}{ID: 1}) // map[string]any{"ID": 1}

gocast.Copy(map[string]any{"ID": 1}) // map[string]any{"ID": 1}
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

## Struct Field Manipulation

GoCast also allows you to work with struct fields dynamically:

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

## Custom Type Support

You can define and use custom types with GoCast:

```go
// Define custom type
type Money int64

func (m *Money) CastSet(ctx context.Context, v any) error {
  switch val := v.(type) {
  case Money:
    *m = val
  default:
    *m = Money(gocast.Float64(v) * 1000000)
  }
  return nil
}

// Use custom type in structs
type Car struct {
  ID int64
  Price Money
}

var car Car

// Mapping values into struct
gocast.TryCopyStruct(&car, map[string]any{"ID":1, "Price": "12000.00"})
```

## Benchmarks

Here are some benchmark results for GoCast:

```sh
> go test -benchmem -v -race -bench=.

goos: darwin
goarch: amd64
pkg: github.com/demdxx/gocast/v2
BenchmarkApproachTest
BenchmarkApproachTest/bench1
BenchmarkApproachTest/bench1-24           348097         3064 ns/op        0 B/op          0 allocs/op
BenchmarkApproachTest/bench2
BenchmarkApproachTest/bench2-24           394160         3005 ns/op        0 B/op          0 allocs/op
BenchmarkBool
BenchmarkBool-24                        20453542           58.87 ns/op             0 B/op          0 allocs/op
BenchmarkToBoolByReflect
BenchmarkToBoolByReflect-24             17354990           70.62 ns/op             0 B/op          0 allocs/op
BenchmarkToFloat
BenchmarkToFloat-24                     10951923          107.0 ns/op              0 B/op          0 allocs/op
BenchmarkToInt
BenchmarkToInt-24                        9870794          121.1 ns/op              0 B/op          0 allocs/op
BenchmarkToUint
BenchmarkToUint-24                       9729873          121.3 ns/op              0 B/op          0 allocs/op
BenchmarkToStringByReflect
BenchmarkToStringByReflect-24             922710         1601 ns/op        5 B/op          0 allocs/op
BenchmarkToString
BenchmarkToString-24                      836929         1622 ns/op        5 B/op          0 allocs/op
BenchmarkGetSetFieldValue
BenchmarkGetSetFieldValue/set
BenchmarkGetSetFieldValue/set-24         1000000         1021 ns/op       64 B/op          4 allocs/op
BenchmarkGetSetFieldValue/get
BenchmarkGetSetFieldValue/get-24         1869465          643.8 ns/op             48 B/op          3 allocs/op
BenchmarkParseTime
BenchmarkParseTime-24                     374346         3130 ns/op      700 B/op         17 allocs/op
BenchmarkIsEmpty
BenchmarkIsEmpty-24                     37383031           31.23 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/demdxx/gocast/v2     17.982s
```

## License

GoCast is released under the MIT License. See the [LICENSE](LICENSE) file for details.
