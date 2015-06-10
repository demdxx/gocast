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
)

func ToStruct(dst, src interface{}, tag string) (err error) {
  if nil == dst || nil == src {
    err = errInvalidParams
  } else {
    s := reflectTarget(reflect.ValueOf(dst))
    t := s.Type()

    switch src.(type) {
    case map[interface{}]interface{}, map[string]interface{}, map[string]string:
      for i := 0; i < s.NumField(); i++ {
        f := s.Field(i)
        v := mapValueByStringKeys(src, fieldNames(t.Field(i), tag))
        f.Set(reflect.ValueOf(To(v, f.Kind())))
      }
      break
    default:
      err = errUnsupportedType
    }
  }
  return
}

func StructFields(st interface{}, tag string) []string {
  fields := []string{}

  s := reflectTarget(reflect.ValueOf(st))
  t := s.Type()

  for i := 0; i < s.NumField(); i++ {
    fname, _ := fieldName(t.Field(i), tag)
    if len(fname) > 0 && "-" != fname {
      fields = append(fields, fname)
    }
  }
  return fields
}

func StructFieldTags(st interface{}, tag string) map[string]string {
  fields := map[string]string{}
  keys, values := StructFieldTagsUnsorted(st, tag)

  for i, k := range keys {
    fields[k] = values[i]
  }
  return fields
}

func StructFieldTagsUnsorted(st interface{}, tag string) ([]string, []string) {
  keys := []string{}
  values := []string{}

  s := reflectTarget(reflect.ValueOf(st))
  t := s.Type()

  for i := 0; i < s.NumField(); i++ {
    f := t.Field(i)
    tag := fieldTag(f, tag)
    if len(tag) > 0 && "-" != tag {
      keys = append(keys, f.Name)
      values = append(values, tag)
    }
  }
  return keys, values
}

///////////////////////////////////////////////////////////////////////////////
/// MARK: Helpers
///////////////////////////////////////////////////////////////////////////////

var fieldNameArr = []string{"field", "schema", "sql", "json", "xml", "yaml"}

func fieldName(f reflect.StructField, tag string) (name string, omitempty bool) {
  names := fieldNames(f, tag)
  name = names[0]
  if len(names) > 0 {
    if "omitempty" == names[len(names)-1] {
      omitempty = true
    }
  }
  return
}

func fieldNames(f reflect.StructField, tag string) []string {
  return strings.Split(fieldTag(f, tag), ",")
}

func fieldTag(f reflect.StructField, tag string) string {
  if "-" != tag {
    var fields string
    var tags []string

    if len(tag) > 0 {
      tags = strings.Split(tag, ",")
    } else {
      tags = fieldNameArr
    }

    for _, k := range tags {
      fields = f.Tag.Get(k)
      if len(fields) > 0 {
        break
      }
    }

    if len(fields) > 0 {
      if "-" == fields {
        return ""
      }
      return fields
    }
  }
  return f.Name
}
