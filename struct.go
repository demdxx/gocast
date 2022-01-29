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
	"database/sql"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Structure method errors
var (
	ErrStructFieldNameUndefined      = errors.New("struct field name undefined")
	ErrStructFieldValueCantBeChanged = errors.New("struct field value cant be changed")
)

// ToStruct convert any input type into the target structure
func ToStruct(dst, src interface{}, tag string) (err error) {
	if dst == nil || src == nil {
		return errInvalidParams
	}

	if sintf, ok := dst.(sql.Scanner); ok {
		if sintf.Scan(src) == nil {
			return nil
		}
	}

	switch dst.(type) {
	case time.Time, *time.Time:
		switch v := src.(type) {
		case time.Time:
			s := reflectTarget(reflect.ValueOf(dst))
			s.Set(reflect.ValueOf(v))
		case *time.Time:
			s := reflectTarget(reflect.ValueOf(dst))
			s.Set(reflect.ValueOf(*v))
		case string:
			var tm time.Time
			if tm, err = ParseTime(v); err == nil {
				s := reflectTarget(reflect.ValueOf(dst))
				s.Set(reflect.ValueOf(tm))
			}
		case int64:
			s := reflectTarget(reflect.ValueOf(dst))
			s.Set(reflect.ValueOf(time.Unix(v, 0)))
		default:
			err = errUnsupportedType
		}
	default:

		s := reflectTarget(reflect.ValueOf(dst))
		t := s.Type()

		switch src.(type) {
		case map[interface{}]interface{}, map[string]interface{}, map[string]string:
			for i := 0; i < s.NumField(); i++ {
				f := s.Field(i)
				if f.CanSet() {
					// Get passable field names
					names := fieldNames(t.Field(i), tag)
					if len(names) < 1 {
						continue
					}

					// Get value from map
					v := mapValueByStringKeys(src, names)

					// Set field value
					if v == nil {
						f.Set(reflect.Zero(f.Type()))
					} else {
						switch f.Kind() {
						case reflect.Struct:
							if err = ToStruct(f.Addr().Interface(), v, tag); err != nil {
								return err
							}
						default:
							var vl interface{}
							if vl, err = ToType(reflect.ValueOf(v), f.Type(), tag); err == nil {
								val := reflect.ValueOf(vl)
								if val.Kind() == reflect.Ptr && val.Kind() != f.Kind() {
									val = val.Elem()
								}
								if val.Kind() == f.Kind() {
									f.Set(val)
								} else {
									err = errUnsupportedType
									break
								}
							} else {
								return err
							}
						} // end switch
					} // end else
				}
			}
		default:
			err = errUnsupportedType
		}
	}
	return err
}

// StructFields returns the field names from the structure
func StructFields(st interface{}, tag string) []string {
	fields := []string{}

	s := reflectTarget(reflect.ValueOf(st))
	t := s.Type()

	for i := 0; i < s.NumField(); i++ {
		fname, _ := fieldName(t.Field(i), tag)
		if len(fname) > 0 && fname != "-" {
			fields = append(fields, fname)
		}
	}
	return fields
}

// StructFieldTags returns Map with key->tag matching
func StructFieldTags(st interface{}, tag string) map[string]string {
	fields := map[string]string{}
	keys, values := StructFieldTagsUnsorted(st, tag)

	for i, k := range keys {
		fields[k] = values[i]
	}
	return fields
}

// StructFieldTagsUnsorted returns field names and tag targets separately
func StructFieldTagsUnsorted(st interface{}, tag string) ([]string, []string) {
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
func StructFieldValue(st interface{}, name string) (interface{}, error) {
	s := reflectTarget(reflect.ValueOf(st))
	t := s.Type()
	if _, ok := t.FieldByName(name); ok {
		return s.FieldByName(name).Interface(), nil
	}
	return nil, errors.Wrap(ErrStructFieldNameUndefined, name)
}

// SetStructFieldValue puts value into the struct field
func SetStructFieldValue(st interface{}, name string, value interface{}) error {
	s := reflectTarget(reflect.ValueOf(st))
	t := s.Type()
	if _, ok := t.FieldByName(name); ok {
		field := s.FieldByName(name)
		if !field.CanSet() {
			return errors.Wrap(ErrStructFieldValueCantBeChanged, name)
		}
		field.Set(reflect.ValueOf(value))
		return nil
	}
	return errors.Wrap(ErrStructFieldNameUndefined, name)
}

///////////////////////////////////////////////////////////////////////////////
/// MARK: Helpers
///////////////////////////////////////////////////////////////////////////////

var fieldNameArr = []string{"field", "schema", "sql", "json", "xml", "yaml"}

func fieldNames(f reflect.StructField, tag string) []string {
	names := fieldTagArr(f, tag)
	switch names[0] {
	case "", "-":
		return []string{f.Name}
	default:
	}
	return []string{names[0], f.Name}
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
	return
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
