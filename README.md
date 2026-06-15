# GoCast

[![GoDoc](https://pkg.go.dev/badge/github.com/demdxx/gocast/v2.svg)](https://pkg.go.dev/github.com/demdxx/gocast/v2)
[![Build Status](https://github.com/demdxx/gocast/workflows/Tests/badge.svg)](https://github.com/demdxx/gocast/actions?workflow=Tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/demdxx/gocast/v2)](https://goreportcard.com/report/github.com/demdxx/gocast/v2)
[![Coverage Status](https://coveralls.io/repos/github/demdxx/gocast/badge.svg?branch=master)](https://coveralls.io/github/demdxx/gocast?branch=master)

GoCast is a Go library for runtime type conversion, deep copying, and struct/map
manipulation using reflection and generics. It is designed for scenarios where
data arrives as `any` — from JSON decoding, configuration files, ORM rows, or
dynamic APIs — and must be coerced into concrete Go types without boilerplate.

**Requires Go 1.21+**

## Features

- **Universal type casting** — `Cast[T]`, `TryCast[T]`, `Number[T]`, `Str`, `Bool` and their `Context` / `Try` variants
- **Deep copy** — circular-reference–safe `TryCopy`, `Copy`; advanced `TryCopyWithOptions` with depth and field-filter controls
- **Struct mapping** — populate a struct from a map or another struct using struct-tag–driven field resolution (`json`, `field`, `sql`, …)
- **Map conversion** — convert structs and maps into typed `map[K]V` with optional recursive field processing
- **Struct walking** — `StructWalk` visits every field recursively with full path tracking
- **Custom types** — implement `CastSetter` for full control over how a type is populated

## Installation

```sh
go get github.com/demdxx/gocast/v2
```

```go
import gocast "github.com/demdxx/gocast/v2"
```

## Quick Start

```go
// String conversion
gocast.Str(8)                    // "8"
gocast.Str(8.31)                 // "8.31"
gocast.Str([]byte("hello"))      // "hello"
gocast.Str(nil)                  // ""

// Numeric conversion (fast — no reflection for common types)
gocast.Number[int](8.31)         // 8
gocast.Number[int]("8")          // 8
gocast.Number[int](true)         // 1
gocast.Number[float32]("2.12")   // 2.12

// Generic cast (handles struct ↔ map conversions too)
var v any = 8
gocast.Cast[int](v)              // 8
gocast.Cast[float64]("2.")       // 2.0

val, err := gocast.TryCast[int]("123.2") // 123, nil

// Struct → map
res := gocast.Map[string, any](struct{ ID int64 }{ID: 1})
// map[string]any{"ID": int64(1)}

// Deep copy
copied, err := gocast.TryCopy(map[string]any{"key": []int{1, 2, 3}})
```

## Type Conversion

Functions follow a consistent naming convention:

| Pattern | Behaviour |
|---------|-----------|
| `TryX(v)` | Returns `(T, error)` |
| `X(v)` | Returns `T`, zero value on error |
| `XContext(ctx, v)` | Same as above, passes context to `CastSetter` hooks |

```go
// All return (T, error)
val, err := gocast.TryCast[int](src)
val, err := gocast.TryCastContext[int](ctx, src, "json")
val, err := gocast.TryNumber[float64](src)
val, err := gocast.TryStr(src)

// All return T (zero on error)
gocast.Cast[int](src)
gocast.Number[float64](src)
gocast.Str(src)
gocast.Bool(src)
gocast.Int(src)
gocast.Int64(src)
gocast.Float64(src)
```

## Deep Copy

`TryCopy` handles circular references automatically via a visited-pointer map.

```go
// Returns (T, error)
copied, err := gocast.TryCopy(original)

// Panics on error — use when the input type is known to be safe
copied := gocast.Copy(original)

// Type-safe helpers for slices and maps
copiedSlice := gocast.CopySlice([]int{1, 2, 3})
copiedMap   := gocast.CopyMap(map[string]int{"a": 1})

// Copy any interface value with type preservation
copied, err := gocast.TryAnyCopy(src)
copied       = gocast.AnyCopy(src)
```

### Copy with options

```go
opts := gocast.CopyOptions{
    MaxDepth:               5,     // zero means unlimited
    IgnoreUnexportedFields: true,  // skip unexported struct fields
    IgnoreCircularRefs:     false, // false = preserve cycle, true = break cycle
}
copied, err := gocast.TryCopyWithOptions(original, opts)
```

### Circular reference example

```go
type Node struct {
    Value int
    Next  *Node
}
node1 := &Node{Value: 1}
node2 := &Node{Value: 2}
node1.Next = node2
node2.Next = node1 // cycle

// TryCopy preserves the cycle in the copy
copied, err := gocast.TryCopy(node1)
// copied.Next.Next == copied  ✓
```

## Struct Mapping

```go
type User struct {
    ID    uint64 `json:"id"`
    Email string `json:"email"`
}

// Populate from a map (tag-driven key resolution)
var u User
err := gocast.TryCopyStruct(&u, map[string]any{"id": 19, "email": "user@example.com"}, "json")

// Or use the generic helper
u, err := gocast.Struct[User](map[string]any{"id": 19, "email": "user@example.com"}, "json")
```

### Individual field access

```go
// Get
id, err := gocast.StructFieldValue(u, "ID")

// Set
err := gocast.SetStructFieldValue(context.Background(), &u, "ID", uint64(19))
err  = gocast.SetStructFieldValue(context.Background(), &u, "Email", "user@example.com")
```

## Map Conversion

```go
type Product struct {
    ID    int64  `json:"id"`
    Title string `json:"title"`
}

src := Product{ID: 1, Title: "Gopher Plush"}

// Struct → map[string]any using json tags
m := gocast.Map[string, any](src, "json")
// map[string]any{"id": int64(1), "title": "Gopher Plush"}

// With recursive nested struct conversion
m, err := gocast.TryMapRecursive[string, any](src, "json")

// Map → map with key/value type conversion
m2 := gocast.Map[string, string](map[int]int{1: 2})
// map[string]string{"1": "2"}
```

## Struct Walking

`StructWalk` visits every exported field recursively. It is useful for populating
structs from external sources (env vars, INI files, config maps).

```go
err := gocast.StructWalk(ctx, &cfg, func(
    ctx    context.Context,
    obj    gocast.StructWalkObject,
    field  gocast.StructWalkField,
    path   []string,
) error {
    // path is the list of parent field names leading to this field
    key := strings.Join(append(path, field.Name()), ".")
    if v, ok := source[key]; ok {
        return field.SetValue(ctx, v)
    }
    return nil
})

// Use a struct tag to name path segments
err = gocast.StructWalk(ctx, &cfg, walker, gocast.WalkWithPathTag("json"))

// Or a custom function
err = gocast.StructWalk(ctx, &cfg, walker, gocast.WalkWithPathExtractor(
    func(ctx context.Context, obj gocast.StructWalkObject, field gocast.StructWalkField, path []string) string {
        return strings.ToLower(field.Name())
    },
))
```

Return `gocast.ErrWalkSkip` from the walker to skip nested fields of a struct.
Return `gocast.ErrWalkStop` to stop the walk entirely (converted to `nil` by `StructWalk`).

## Custom Types

Implement `CastSetter` to control how a type is populated during struct mapping
or slice/map conversion:

```go
type Money int64

func (m *Money) CastSet(ctx context.Context, v any) error {
    switch val := v.(type) {
    case Money:
        *m = val
    default:
        *m = Money(gocast.Float64(v) * 1_000_000)
    }
    return nil
}

type Order struct {
    ID    int64
    Total Money
}

var order Order
err := gocast.TryCopyStruct(&order, map[string]any{
    "ID":    1,
    "Total": "12000.00",
})
// order.Total == 12_000_000_000
```

## Error Handling

```go
// Try* variants — always return an error
copied, err := gocast.TryCopy(v)
val,    err := gocast.TryCast[int](v)

// Non-Try variants — return zero value on error (never panic)
val  := gocast.Cast[int](v)
val  := gocast.Number[int](v)

// Copy / CopySlice / CopyMap — panic on unsupported type (e.g. chan, func)
copied := gocast.Copy(v)
```

Errors are wrapped and can be inspected with `errors.Is`:

```go
_, err := gocast.TryCast[int](make(chan int))
errors.Is(err, gocast.ErrCopyUnsupportedType) // true
```

## Benchmarks

```
goos: darwin / goarch: arm64 / cpu: Apple M2 Ultra

BenchmarkBool-24                  20453542     58.87 ns/op      0 B/op   0 allocs/op
BenchmarkToFloat-24               10951923    107.0  ns/op      0 B/op   0 allocs/op
BenchmarkToInt-24                  9870794    121.1  ns/op      0 B/op   0 allocs/op
BenchmarkToString-24                836929   1622    ns/op      5 B/op   0 allocs/op

BenchmarkCopy/simple_int-24      100000000     10.03 ns/op      8 B/op   1 allocs/op
BenchmarkCopy/simple_string-24    28639788     41.20 ns/op     32 B/op   2 allocs/op
BenchmarkCopy/simple_struct-24    16650103     71.49 ns/op     48 B/op   2 allocs/op
BenchmarkCopy/complex_struct-24    1796786    663.3  ns/op    632 B/op  18 allocs/op
BenchmarkCopy/slice-24             5762702    205.5  ns/op    152 B/op   4 allocs/op
BenchmarkCopy/map-24               1966965    607.1  ns/op    504 B/op  23 allocs/op

BenchmarkIsEmpty-24               37383031     31.23 ns/op      0 B/op   0 allocs/op
BenchmarkParseTime-24               374346   3130    ns/op    700 B/op  17 allocs/op
```

Run benchmarks locally:

```sh
make bench
```

## API Reference

### Type Conversion

```go
func Cast[R, T any](v T, tags ...string) R
func TryCast[R, T any](v T, tags ...string) (R, error)
func CastContext[R, T any](ctx context.Context, v T, tags ...string) R
func TryCastContext[R, T any](ctx context.Context, v T, tags ...string) (R, error)

func CastRecursive[R, T any](v T, tags ...string) R
func TryCastRecursive[R, T any](v T, tags ...string) (R, error)

func Number[R Numeric](v any) R
func TryNumber[R Numeric](v any) (R, error)

func Str(v any) string           // string conversion
func Bool(v any) bool            // bool conversion
func Int(v any) int              // integer helpers
func Int8(v any) int8
func Int16(v any) int16
func Int32(v any) int32
func Int64(v any) int64
func Uint(v any) uint
func Uint8(v any) uint8
func Uint16(v any) uint16
func Uint32(v any) uint32
func Uint64(v any) uint64
func Float64(v any) float64
func Float32(v any) float32
```

### Deep Copy

```go
func TryCopy[T any](src T) (T, error)
func Copy[T any](src T) T

func TryCopyWithOptions[T any](src T, opts CopyOptions) (T, error)

func CopySlice[T any](src []T) []T
func CopyMap[K comparable, V any](src map[K]V) map[K]V

func TryAnyCopy(src any) (any, error)
func AnyCopy(src any) any
```

### Struct & Map

```go
func TryCopyStruct(dst, src any, tags ...string) error
func Struct[R any](src any, tags ...string) (R, error)

func ToMap(dst, src any, recursive bool, tags ...string) error
func Map[K comparable, V any](src any, tags ...string) map[K]V
func TryMap[K comparable, V any](src any, tags ...string) (map[K]V, error)
func MapRecursive[K comparable, V any](src any, tags ...string) map[K]V
func TryMapRecursive[K comparable, V any](src any, tags ...string) (map[K]V, error)
func ToMapFrom(src any, recursive bool, tags ...string) (map[any]any, error)

func StructFieldValue(st any, names ...string) (any, error)
func SetStructFieldValue(ctx context.Context, st any, name string, value any) error
func StructFieldNames(st any, tag string) []string
func StructFieldTags(st any, tag string) map[string]string

func StructWalk(ctx context.Context, v any, walker func(...) error, options ...WalkOption) error
func WalkWithPathTag(tagName string) WalkOption
func WalkWithPathExtractor(fn func(...) string) WalkOption
```

### Utilities

```go
func IsEmpty(v any) bool
func IsNil(v any) bool
func IsStr(v any) bool
func IsSlice(v any) bool
func IsMap(v any) bool
func IsStruct(v any) bool
func ParseTime(v string) (time.Time, error)
func Or[T comparable](vals ...T) T
func IfThen[T any](cond bool, a, b T) T
func Ptr[T any](v T) *T
func PtrAsValue[T any](v *T) T
```

### Errors

```go
var ErrInvalidParams                 = errors.New("invalid params")
var ErrUnsupportedType               = errors.New("unsupported destination type")
var ErrUnsupportedSourceType         = errors.New("unsupported source type")
var ErrUnsettableValue               = errors.New("can't set value")
var ErrStructFieldNameUndefined      = errors.New("struct field name undefined")
var ErrStructFieldValueCantBeChanged = errors.New("struct field value cant be changed")
var ErrCopyUnsupportedType           = errors.New("copy: unsupported type")
var ErrCopyInvalidValue              = errors.New("copy: invalid value")
var ErrWalkSkip                      = errors.New("skip field walk")
var ErrWalkStop                      = errors.New("stop field walk")
```

### Deprecated

The following identifiers are deprecated and will be removed in v3:

| Identifier | Replacement |
|------------|-------------|
| `MustCopy[T]` | `Copy[T]` |
| `Float(v)` | `Float64(v)` |
| `ToUint64ByReflect(v)` | `ReflectToUint64(v)` |
| `ShallowCopy[T]` | direct assignment |
| `Any` struct | top-level conversion functions |
| `ErrCopyCircularReference` | _(never returned; remove usage)_ |

## Examples

See the [examples/](examples/) directory for runnable usage examples and the
[test files](.) for comprehensive coverage of all API paths.

## License

GoCast is released under the MIT License. See the [LICENSE](LICENSE) file for details.
