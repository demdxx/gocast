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
    s := reflect.ValueOf(dst).Elem()
    t := s.Type()

    switch src.(type) {
    case map[interface{}]interface{}:
      for i := 0; i < s.NumField(); i++ {
        f := s.Field(i)
        v := mapValueByStringKeys(src, fieldNames(t.Field(i), tag))
        f.Set(reflect.ValueOf(v))
      }
      break
    case map[string]interface{}:
      for i := 0; i < s.NumField(); i++ {
        dst.(map[string]interface{})[fieldName(t.Field(i), tag)] = s.Field(i).Interface()
      }
      break
    case map[string]string:
      for i := 0; i < s.NumField(); i++ {
        dst.(map[string]string)[fieldName(t.Field(i), tag)] = ToString(s.Field(i).Interface())
      }
      break
    default:
      err = errUnsupportedType
    }
  }
  return
}

///////////////////////////////////////////////////////////////////////////////
/// MARK: Helpers
///////////////////////////////////////////////////////////////////////////////

var fieldNameArr = []string{"field", "schema", "sql", "json", "xml", "yaml"}

func fieldName(f reflect.StructField, tag string) string {
  return fieldNames(f, tag)[0]
}

func fieldNames(f reflect.StructField, tag string) []string {
  return strings.Split(fieldTag(f, tag), ",")
}

func fieldTag(f reflect.StructField, tag string) string {
  if "-" != tag {
    var fields string
    if len(tag) > 0 {
      fields = f.Tag.Get(tag)
    } else {
      for _, k := range fieldNameArr {
        fields = f.Tag.Get(k)
        if len(fields) > 0 {
          break
        }
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
