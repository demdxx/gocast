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
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// CastSetter interface from some type into the specific value
type CastSetter interface {
	CastSet(v any) error
}

// TryCopyStruct convert any input type into the target structure
func TryCopyStruct(dst, src any, tags ...string) (err error) {
	if dst == nil || src == nil {
		return errInvalidParams
	}

	if sintf, ok := dst.(CastSetter); ok {
		if sintf.CastSet(src) == nil {
			return nil
		}
	}

	switch dst.(type) {
	case time.Time, *time.Time:
		switch v := src.(type) {
		case time.Time:
			s := reflectTarget(reflect.ValueOf(dst))
			err = setFieldValue(s, v)
		case *time.Time:
			s := reflectTarget(reflect.ValueOf(dst))
			err = setFieldValue(s, *v)
		case string:
			var tm time.Time
			if tm, err = ParseTime(v); err == nil {
				s := reflectTarget(reflect.ValueOf(dst))
				err = setFieldValue(s, tm)
			}
		case int64:
			s := reflectTarget(reflect.ValueOf(dst))
			err = setFieldValue(s, time.Unix(v, 0))
		case uint64:
			s := reflectTarget(reflect.ValueOf(dst))
			err = setFieldValue(s, time.Unix(int64(v), 0))
		default:
			err = errUnsupportedType
		}
	default:
		s := reflectTarget(reflect.ValueOf(dst))
		t := s.Type()

		switch src.(type) {
		case map[any]any, map[string]any, map[string]string:
			for i := 0; i < s.NumField(); i++ {
				f := s.Field(i)
				if f.CanSet() {
					// Get passable field names
					names := fieldNames(t.Field(i), tags...)
					if len(names) < 1 {
						continue
					}

					// Get value from map
					v := mapValueByStringKeys(src, names)

					// Set field value
					if v == nil {
						err = setFieldValueReflect(f, reflect.Zero(f.Type()))
					} else {
						switch f.Kind() {
						case reflect.Struct:
							if err = TryCopyStruct(f.Addr().Interface(), v, tags...); err != nil {
								return err
							}
						default:
							var vl any
							if vl, err = TryToType(v, f.Type(), tags...); err == nil {
								val := reflect.ValueOf(vl)
								if val.Kind() == reflect.Ptr && val.Kind() != f.Kind() {
									val = val.Elem()
								}
								if val.Kind() == f.Kind() || f.Kind() == reflect.Interface {
									err = setFieldValueReflect(f, val)
								} else {
									err = errUnsupportedType
								}
							} else if setter, _ := f.Interface().(CastSetter); setter != nil {
								err = setter.CastSet(v)
							} else if f.CanAddr() {
								if setter, _ := f.Addr().Interface().(CastSetter); setter != nil {
									err = setter.CastSet(v)
								}
							}
						} // end switch
					} // end else
					if err != nil {
						break
					}
				}
			}
		default:
			err = errUnsupportedType
		}
	}
	return err
}

// Struct convert any input type into the target structure
func Struct[R any](src any, tags ...string) (R, error) {
	var res R
	err := TryCopyStruct(&res, src, tags...)
	return res, err
}

// ToStruct convert any input type into the target structure
//
// Deprecated: Use TryCopyStruct instead
func ToStruct(dst, src any, tags ...string) error {
	return TryCopyStruct(dst, src, tags...)
}

// StructFields returns the field names from the structure
func StructFields(st any, tag string) []string {
	s := reflectTarget(reflect.ValueOf(st))
	t := s.Type()

	fields := make([]string, 0, s.NumField())
	for i := 0; i < s.NumField(); i++ {
		fname, _ := fieldName(t.Field(i), tag)
		if fname != "" && fname != "-" {
			fields = append(fields, fname)
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
	keys := []string{}
	values := []string{}

	s := reflectTarget(reflect.ValueOf(st))
	t := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := t.Field(i)
		tag := fieldTag(f, tag)
		if len(tag) > 0 && tag != "-" {
			keys = append(keys, f.Name)
			values = append(values, tag)
		}
	}
	return keys, values
}

// StructFieldValue returns the value of the struct field
func StructFieldValue(st any, name string) (any, error) {
	s := reflectTarget(reflect.ValueOf(st))
	t := s.Type()
	if _, ok := t.FieldByName(name); ok {
		return s.FieldByName(name).Interface(), nil
	}
	return nil, errors.Wrap(ErrStructFieldNameUndefined, name)
}

// SetStructFieldValue puts value into the struct field
func SetStructFieldValue(st any, name string, value any) (err error) {
	s := reflectTarget(reflect.ValueOf(st))
	t := s.Type()
	if _, ok := t.FieldByName(name); ok {
		field := s.FieldByName(name)
		if !field.CanSet() {
			return errors.Wrap(ErrStructFieldValueCantBeChanged, name)
		}
		switch field.Interface().(type) {
		case time.Time, *time.Time:
			switch v := value.(type) {
			case time.Time:
				s := reflectTarget(field)
				err = setFieldValue(s, v)
			case *time.Time:
				s := reflectTarget(field)
				err = setFieldValue(s, *v)
			case string:
				var tm time.Time
				if tm, err = ParseTime(v); err == nil {
					s := reflectTarget(field)
					err = setFieldValue(s, tm)
				}
			case int64:
				s := reflectTarget(field)
				err = setFieldValue(s, time.Unix(v, 0))
			case uint64:
				s := reflectTarget(field)
				err = setFieldValue(s, time.Unix(int64(v), 0))
			default:
				err = errUnsupportedType
			}
		default:
			err = setFieldValue(field, value)
		}
		return err
	}
	return errors.Wrap(ErrStructFieldNameUndefined, name)
}

///////////////////////////////////////////////////////////////////////////////
/// MARK: Helpers
///////////////////////////////////////////////////////////////////////////////

var fieldNameArr = []string{"field", "schema", "sql", "json", "xml", "yaml"}

func setFieldValue(field reflect.Value, val any) error {
	if setter, _ := field.Interface().(CastSetter); setter != nil {
		return setter.CastSet(val)
	} else if field.CanAddr() {
		if setter, _ := field.Addr().Interface().(CastSetter); setter != nil {
			return setter.CastSet(val)
		} else if vl := reflect.ValueOf(val); field.Kind() == vl.Kind() || field.Kind() == reflect.Interface {
			field.Set(vl)
		} else {
			return errors.Wrap(errUnsupportedType, field.Type().String())
		}
	} else if vl := reflect.ValueOf(val); field.Kind() == vl.Kind() {
		field.Set(vl)
	} else {
		return errors.Wrap(errUnsupportedType, field.Type().String())
	}
	return nil
}

func setFieldValueReflect(field, val reflect.Value) error {
	if setter, _ := field.Interface().(CastSetter); setter != nil {
		return setter.CastSet(val.Interface())
	} else if field.CanAddr() {
		if setter, _ := field.Addr().Interface().(CastSetter); setter != nil {
			return setter.CastSet(val.Interface())
		} else if field.Kind() == val.Kind() || field.Kind() == reflect.Interface {
			field.Set(val)
		} else {
			return errors.Wrap(errUnsupportedType, field.Type().String())
		}
	} else if field.Kind() == val.Kind() {
		field.Set(val)
	} else {
		return errors.Wrap(errUnsupportedType, field.Type().String())
	}
	return nil
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
