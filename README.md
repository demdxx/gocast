# GoCast

[![GoDoc](https://godoc.org/github.com/demdxx/gocast?status.svg)](https://godoc.org/github.com/demdxx/gocast)
[![Build Status](https://github.com/demdxx/gocast/workflows/Tests/badge.svg)](https://github.com/demdxx/gocast/actions?workflow=Tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/demdxx/gocast)](https://goreportcard.com/report/github.com/demdxx/gocast)
[![Coverage Status](https://coveralls.io/repos/github/demdxx/gocast/badge.svg?branch=master)](https://coveralls.io/github/demdxx/gocast?branch=master)

## Introduction

GoCast is a powerful Go library that allows you to easily convert between different basic types in a consistent and efficient way. Whether you need to convert strings, numbers, or perform deep copying of complex data structures, GoCast has got you covered.

### ✨ Latest Improvements

- **🚀 Enhanced Deep Copy**: Advanced deep copying with circular reference detection and custom options
- **🎯 Type-Safe Functions**: Specialized `CopySlice[T]()` and `CopyMap[K,V]()` functions with better performance
- **⚙️ Flexible Options**: `CopyOptions` for controlling copy behavior (depth limits, field filtering)
- **🔧 Better Architecture**: Modular design with improved error handling and edge case support
- **📈 Performance**: Optimized paths for different data types with minimal memory allocations

## Features

- **Universal Type Casting**: GoCast provides a set of methods for universal type casting, making it easy to convert data between various types.
- **Deep Copy Operations**: Advanced deep copying with support for circular references, custom options, and type-safe specialized functions.
- **Struct Field Manipulation**: GoCast allows you to set and retrieve values of struct fields dynamically, making it a valuable tool for working with complex data structures.
- **Custom Type Support**: You can define custom types and implement your own conversion logic, giving you full control over how data is cast.
- **High Performance**: Optimized for speed with specialized paths for different data types and minimal allocations.

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

## Deep Copy Operations

GoCast provides powerful deep copying capabilities with support for complex data structures:

```go
// Basic deep copy
original := map[string]any{"data": []int{1, 2, 3}}
copied, err := gocast.TryCopy(original)

// Type-safe specialized functions
originalSlice := []int{1, 2, 3, 4, 5}
copiedSlice := gocast.CopySlice(originalSlice)

originalMap := map[string]int{"a": 1, "b": 2}
copiedMap := gocast.CopyMap(originalMap)

// Copy with panic on error
result := gocast.MustCopy(complexStruct)

// Copy with custom options
opts := gocast.CopyOptions{
    MaxDepth: 5,                     // Limit recursion depth
    IgnoreUnexportedFields: true,    // Skip unexported fields
    IgnoreCircularRefs: false,       // Handle circular references
}
copied, err := gocast.TryCopyWithOptions(original, opts)

// Circular reference handling
type Node struct {
    Value int
    Next  *Node
}
node1 := &Node{Value: 1}
node2 := &Node{Value: 2}
node1.Next = node2
node2.Next = node1  // Circular reference

copiedNode, err := gocast.TryCopy(node1) // Handles circular refs automatically
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
goarch: arm64
pkg: github.com/demdxx/gocast/v2
cpu: Apple M2 Ultra

# Core functionality benchmarks
BenchmarkBool-24                        20453542           58.87 ns/op             0 B/op          0 allocs/op
BenchmarkToFloat-24                     10951923          107.0 ns/op              0 B/op          0 allocs/op
BenchmarkToInt-24                        9870794          121.1 ns/op              0 B/op          0 allocs/op
BenchmarkToString-24                      836929         1622 ns/op        5 B/op          0 allocs/op

# Deep copy benchmarks
BenchmarkCopy/simple_int-24             100000000          10.03 ns/op       8 B/op           1 allocs/op
BenchmarkCopy/simple_string-24           28639788          41.20 ns/op      32 B/op           2 allocs/op
BenchmarkCopy/simple_struct-24           16650103          71.49 ns/op      48 B/op           2 allocs/op
BenchmarkCopy/complex_struct-24           1796786         663.3 ns/op     632 B/op          18 allocs/op
BenchmarkCopy/slice-24                    5762702         205.5 ns/op     152 B/op           4 allocs/op
BenchmarkCopy/map-24                      1966965         607.1 ns/op     504 B/op          23 allocs/op

# Specialized copy functions (type-safe and faster)
BenchmarkSpecializedFunctions/CopySlice_specialized-24    9462444         124.9 ns/op     160 B/op          11 allocs/op
BenchmarkSpecializedFunctions/CopyMap_specialized-24      2794791         427.3 ns/op     456 B/op          17 allocs/op

# Other operations
BenchmarkGetSetFieldValue/set-24         1000000         1021 ns/op       64 B/op           4 allocs/op
BenchmarkGetSetFieldValue/get-24         1869465          643.8 ns/op      48 B/op           3 allocs/op
BenchmarkParseTime-24                     374346         3130 ns/op      700 B/op          17 allocs/op
BenchmarkIsEmpty-24                     37383031           31.23 ns/op       0 B/op           0 allocs/op
```

### Performance Highlights

- **Deep Copy**: Highly optimized with specialized functions for different types
- **Type-Safe Copies**: `CopySlice` and `CopyMap` provide better performance for specific types
- **Memory Efficient**: Minimal allocations for most operations
- **Circular References**: Automatic detection and handling with zero performance cost for non-circular data

## API Reference

### Core Copy Functions

```go
// Basic deep copy with error handling
func TryCopy[T any](src T) (T, error)

// Deep copy with panic on error
func Copy[T any](src T) T

// Panic version of TryCopy
func MustCopy[T any](src T) T

// Copy any type (interface{})
func TryAnyCopy(src any) (any, error)
func AnyCopy(src any) any
```

### Specialized Copy Functions

```go
// Type-safe slice copying (40% faster than TryCopy for slices)
func CopySlice[T any](src []T) []T

// Type-safe map copying (30% faster than TryCopy for maps)
func CopyMap[K comparable, V any](src map[K]V) map[K]V

// Shallow copy (for immutable data)
func ShallowCopy[T any](src T) T

// Interface copy with type preservation
func CopyInterface(src any) (any, error)
```

### Advanced Copy Options

```go
type CopyOptions struct {
    IgnoreUnexportedFields bool  // Skip unexported struct fields
    MaxDepth               int   // Limit recursion depth
    IgnoreCircularRefs     bool  // Ignore circular references instead of preserving them
}

// Copy with custom options
func TryCopyWithOptions[T any](src T, opts CopyOptions) (T, error)
```

### Usage Recommendations

- **For simple types**: Use `TryCopy` or `Copy` - they're optimized automatically
- **For slices**: Use `CopySlice[T]()` for better type safety and performance
- **For maps**: Use `CopyMap[K,V]()` for better type safety and performance  
- **For complex nested data**: Use `TryCopyWithOptions` with `MaxDepth` to control resource usage
- **For circular references**: `TryCopy` handles them automatically
- **For performance-critical code**: Consider if deep copy is needed - immutable types don't require copying

### Examples

See the [examples directory](examples/) for more detailed usage examples and the [test files](.) for comprehensive usage patterns.

### Documentation

Full API documentation is available at [GoDoc](https://godoc.org/github.com/demdxx/gocast/v2).

## License

GoCast is released under the MIT License. See the [LICENSE](LICENSE) file for details.
