// Copyright (c) 2014 Dmitry Ponomarev <demdxx@gmail.com>
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
	"context"
	"reflect"
	"strings"
	"time"
)

// TryCopyStruct convert any input type into the target structure
func TryCopyStruct(dst, src any, tags ...string) (err error) {
	return TryCopyStructContext(context.Background(), dst, src, tags...)
}

// TryCopyStructContext convert any input type into the target structure
func TryCopyStructContext(ctx context.Context, dst, src any, tags ...string) (err error) {
	if dst == nil || src == nil {
		return ErrInvalidParams
	}

	if sintf, ok := dst.(CastSetter); ok {
		if sintf.CastSet(ctx, src) == nil {
			return nil
		}
	}

	// Set time value in case of time.Time type as destination target
	switch dst.(type) {
	case time.Time, *time.Time:
		return setFieldTimeValue(reflect.ValueOf(dst), src)
	}

	var (
		destVal        = reflectTarget(reflect.ValueOf(dst))
		destType       = destVal.Type()
		destFieldTypes = ReflectStructFields(destType)
		srcVal         = reflectTarget(reflect.ValueOf(src))
		names          []string
		v              any
	)

	// Check source type is map or struct, otherwise return error unsupported type
	if srcVal.Kind() != reflect.Map && srcVal.Kind() != reflect.Struct {
		return wrapError(ErrUnsupportedSourceType, destType.Name())
	}

	// Iterate over destination fields and set values from source
	for _, ft := range destFieldTypes {
		field := destVal.FieldByName(ft.Name)
		if !field.CanSet() {
			continue
		}

		// Get passable field names
		if names = fieldNames(ft, tags...); len(names) < 1 {
			continue
		}

		// Get value from map or struct
		if srcVal.Kind() == reflect.Map {
			v = reflectMapValueByStringKeys(srcVal, names)
		} else {
			v, _ = ReflectStructFieldValue(srcVal, names...)
		}

		// Set field value
		if v == nil {
			if err = setFieldValueReflect(ctx, field, reflect.Zero(field.Type())); err != nil {
				break
			}
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			err = TryCopyStructContext(ctx, field.Addr().Interface(), v, tags...)
		default:
			var vl any
			if vl, err = TryToTypeContext(ctx, v, field.Type(), tags...); err == nil {
				val := reflect.ValueOf(vl)
				if val.Kind() == reflect.Ptr && val.Kind() != field.Kind() {
					val = val.Elem()
				}
				err = setFieldValueReflect(ctx, field, val)
			} else if setter, _ := field.Interface().(CastSetter); setter != nil {
				err = setter.CastSet(ctx, v)
			} else if field.CanAddr() {
				if setter, _ := field.Addr().Interface().(CastSetter); setter != nil {
					err = setter.CastSet(ctx, v)
				}
			}
		}

		if err != nil {
			break
		}
	}

	return err
}

// Struct convert any input type into the target structure
func Struct[R any](src any, tags ...string) (R, error) {
	return StructContext[R](context.Background(), src, tags...)
}

// StructContext convert any input type into the target structure
func StructContext[R any](ctx context.Context, src any, tags ...string) (R, error) {
	var res R
	err := TryCopyStructContext(ctx, &res, src, tags...)
	return res, err
}

// ToStruct convert any input type into the target structure
//
// Deprecated: Use TryCopyStruct instead
func ToStruct(dst, src any, tags ...string) error {
	return TryCopyStruct(dst, src, tags...)
}

// StructFieldNames returns the field names from the structure
func StructFieldNames(st any, tag string) []string {
	s := reflectTarget(reflect.ValueOf(st))
	return ReflectStructFieldNames(s.Type(), tag)
}

// ReflectStructFieldNames returns the field names from the structure
func ReflectStructFieldNames(t reflect.Type, tag string) []string {
	if t.Kind() != reflect.Struct {
		return nil
	}
	fields := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tfield := t.Field(i)
		if tfield.Anonymous && tfield.Type.Kind() == reflect.Struct {
			fields = append(fields, ReflectStructFieldNames(tfield.Type, tag)...)
		} else {
			fname, _ := fieldName(tfield, tag)
			if fname != "" && fname != "-" {
				fields = append(fields, fname)
			}
		}
	}
	return fields
}

// StructFieldTags returns Map with key->tag matching
func StructFieldTags(st any, tag string) map[string]string {
	fields := map[string]string{}
	keys, values := StructFieldTagsUnsorted(st, tag)

	for i, k := range keys {
		fields[k] = values[i]
	}
	return fields
}

// StructFieldTagsUnsorted returns field names and tag targets separately
func StructFieldTagsUnsorted(st any, tag string) ([]string, []string) {
	return ReflectStructFieldTagsUnsorted(reflectTarget(reflect.ValueOf(st)).Type(), tag)
}

// ReflectStructFieldTagsUnsorted returns field names and tag targets separately
func ReflectStructFieldTagsUnsorted(t reflect.Type, tag string) ([]string, []string) {
	if t.Kind() != reflect.Struct {
		return nil, nil
	}

	var (
		keys = make([]string, 0, t.NumField())
		tags = make([]string, 0, t.NumField())
	)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			k, v := ReflectStructFieldTagsUnsorted(f.Type, tag)
			keys = append(keys, k...)
			tags = append(tags, v...)
		} else {
			tag := strings.TrimSuffix(fieldTag(f, tag), ",omitempty")
			if len(tag) > 0 && tag != "-" {
				keys = append(keys, f.Name)
				tags = append(tags, tag)
			}
		}
	}
	return keys, tags
}

// ReflectStructFields returns the field names from the structure
func ReflectStructFields(t reflect.Type) []reflect.StructField {
	if t.Kind() != reflect.Struct {
		return nil
	}
	fields := make([]reflect.StructField, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tfield := t.Field(i)
		if tfield.Anonymous && tfield.Type.Kind() == reflect.Struct {
			fields = append(fields, ReflectStructFields(tfield.Type)...)
		} else {
			fields = append(fields, tfield)
		}
	}
	return fields
}

// StructFieldValue returns the value of the struct field
func StructFieldValue(st any, names ...string) (any, error) {
	structVal := reflectTarget(reflect.ValueOf(st))
	return ReflectStructFieldValue(structVal, names...)
}

// ReflectStructFieldValue returns the value of the struct field
func ReflectStructFieldValue(st reflect.Value, names ...string) (any, error) {
	structVal := reflectTarget(st)
	structType := structVal.Type()
	for _, name := range names {
		if _, ok := structType.FieldByName(name); ok {
			return structVal.FieldByName(name).Interface(), nil
		}
	}
	return nil, wrapError(ErrStructFieldNameUndefined, strings.Join(names, ", "))
}

// SetStructFieldValue puts value into the struct field
func SetStructFieldValue(ctx context.Context, st any, name string, value any) (err error) {
	s := reflectTarget(reflect.ValueOf(st))
	t := s.Type()
	if _, ok := t.FieldByName(name); ok {
		field := s.FieldByName(name)
		if !field.CanSet() {
			return wrapError(ErrStructFieldValueCantBeChanged, name)
		}
		err = setFieldValue(ctx, field, value)
		return err
	}
	return wrapError(ErrStructFieldNameUndefined, name)
}

// IsStruct returns true if the value is a struct
func IsStruct(v any) bool {
	return v != nil && reflect.TypeOf(v).Kind() == reflect.Struct
}

///////////////////////////////////////////////////////////////////////////////
/// MARK: Helpers
///////////////////////////////////////////////////////////////////////////////

var fieldNameArr = []string{"field", "schema", "sql", "json", "xml", "yaml"}

func setFieldValue(ctx context.Context, field reflect.Value, value any) (err error) {
	switch field.Interface().(type) {
	case time.Time, *time.Time:
		err = setFieldTimeValue(field, value)
	default:
		if setter, _ := field.Interface().(CastSetter); setter != nil {
			return setter.CastSet(ctx, value)
		} else if field.CanAddr() {
			if setter, _ := field.Addr().Interface().(CastSetter); setter != nil {
				return setter.CastSet(ctx, value)
			} else if vl := reflect.ValueOf(value); field.Kind() == vl.Kind() || field.Kind() == reflect.Interface {
				field.Set(vl)
			} else {
				return wrapError(ErrUnsupportedType, field.Type().String())
			}
		} else if vl := reflect.ValueOf(value); field.Kind() == vl.Kind() || field.Kind() == reflect.Interface {
			field.Set(vl)
		} else {
			return wrapError(ErrUnsupportedType, field.Type().String())
		}
	}
	return err
}

func setFieldValueReflect(ctx context.Context, field, value reflect.Value) (err error) {
	switch field.Interface().(type) {
	case time.Time, *time.Time:
		err = setFieldTimeValue(field, value.Interface())
	default:
		if setter, _ := field.Interface().(CastSetter); setter != nil {
			return setter.CastSet(ctx, value.Interface())
		} else if field.CanAddr() {
			if setter, _ := field.Addr().Interface().(CastSetter); setter != nil {
				return setter.CastSet(ctx, value.Interface())
			} else if field.Kind() == value.Kind() || field.Kind() == reflect.Interface {
				field.Set(value)
			} else {
				return wrapError(ErrUnsupportedType, field.Type().String())
			}
		} else if field.Kind() == value.Kind() || field.Kind() == reflect.Interface {
			field.Set(value)
		} else {
			return wrapError(ErrUnsupportedType, field.Type().String())
		}
	}
	return err
}

func setFieldTimeValue(field reflect.Value, value any) (err error) {
	switch v := value.(type) {
	case nil:
		s := reflectTarget(field)
		s.Set(reflect.ValueOf(time.Time{}))
	case time.Time:
		s := reflectTarget(field)
		s.Set(reflect.ValueOf(v))
	case *time.Time:
		s := reflectTarget(field)
		s.Set(reflect.ValueOf(*v))
	case string:
		var tm time.Time
		if tm, err = ParseTime(v); err == nil {
			s := reflectTarget(field)
			s.Set(reflect.ValueOf(tm))
		}
	case int64:
		s := reflectTarget(field)
		s.Set(reflect.ValueOf(time.Unix(v, 0)))
	case uint64:
		s := reflectTarget(field)
		s.Set(reflect.ValueOf(time.Unix(int64(v), 0)))
	default:
		err = wrapError(ErrUnsupportedType, field.String())
	}
	return err
}

func fieldNames(f reflect.StructField, tags ...string) []string {
	if len(tags) > 0 {
		names := fieldTagArr(f, tags[0])
		switch names[0] {
		case "", "-":
			return []string{f.Name}
		default:
		}
		return []string{names[0], f.Name}
	}
	return []string{f.Name, f.Name}
}

func fieldNameFromTags(f reflect.StructField, tags ...string) (name string, omitempty bool) {
	if len(tags) == 0 {
		return f.Name, false
	}
	return fieldName(f, tags[0])
}

func fieldName(f reflect.StructField, tag string) (name string, omitempty bool) {
	names := fieldTagArr(f, tag)
	name = names[0]
	if len(names) > 1 && names[len(names)-1] == "omitempty" {
		omitempty = true
	}
	if name == "" {
		name = f.Name
	}
	return name, omitempty
}

func fieldTagArr(f reflect.StructField, tag string) []string {
	return strings.Split(fieldTag(f, tag), ",")
}

func fieldTag(f reflect.StructField, tag string) string {
	if tag == "-" {
		return f.Name
	}
	var (
		fields string
		tags   []string
	)
	if tag != "" {
		tags = strings.Split(tag, ",")
	} else {
		tags = fieldNameArr
	}
	for _, k := range tags {
		fields = f.Tag.Get(k)
		if fields != "" {
			break
		}
	}
	if fields != "" {
		if fields == "-" {
			return ""
		}
		return fields
	}
	return f.Name
}
