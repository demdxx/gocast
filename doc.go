// Package gocast provides runtime type conversion, deep copying, and
// struct/map manipulation using reflection and generics.
//
// # Overview
//
// gocast is designed for scenarios where data arrives as [any] (JSON decoding,
// configuration files, ORM rows, dynamic APIs) and must be coerced into
// concrete Go types without boilerplate.
//
// # Type Conversion
//
// The primary entry points for type conversion are:
//
//   - [Cast] / [TryCast] — convert any value to a target generic type T.
//   - [Number] / [TryNumber] — fast conversion to any numeric type.
//   - [Str] / [TryStr] — convert any value to string.
//   - [Bool] — convert any value to bool.
//
// By convention:
//   - Functions prefixed with Try return (T, error).
//   - Functions without the prefix return T (zero value on error).
//   - Functions suffixed with Context accept a [context.Context] for custom
//     [CastSetter] hooks.
//
// # Deep Copy
//
//   - [TryCopy] / [Copy] — deep copy any value; circular references are handled
//     automatically via a visited-pointer map.
//   - [TryCopyWithOptions] — deep copy with [CopyOptions] (max depth,
//     unexported-field skipping, circular-reference ignoring).
//   - [CopySlice] / [CopyMap] — type-safe helpers for slices and maps.
//
// # Struct and Map Mapping
//
//   - [TryCopyStruct] / [Struct] — populate a struct from a map or another struct
//     using struct-tag–driven field name resolution (json, field, sql, …).
//   - [ToMap] / [Map] / [TryMap] — convert a struct or map into a map type.
//   - [StructWalk] — recursively visit all fields of a struct.
//   - [SetStructFieldValue] / [StructFieldValue] — get or set individual struct
//     fields by name using reflection.
//
// # Custom Types
//
// Types can participate in the conversion pipeline by implementing [CastSetter]:
//
//	type Money int64
//
//	func (m *Money) CastSet(ctx context.Context, v any) error {
//	    *m = Money(gocast.Float64(v) * 1_000_000)
//	    return nil
//	}
//
// # Deprecated APIs
//
// The following identifiers are deprecated and will be removed in v3:
//
//   - [Any] — use the top-level conversion functions directly.
//   - [ErrCopyCircularReference] — never returned; circular refs are handled transparently.
//   - [MustCopy] — identical to [Copy]; use [Copy] instead.
//   - [Float] — alias for [Float64]; use [Float64] instead.
//   - [ToUint64ByReflect] — use [ReflectToUint64] instead.
//   - [ShallowCopy] — identity function; assign the value directly.
package gocast
