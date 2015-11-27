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
  "fmt"
  "reflect"
)

func ToStringByReflect(v reflect.Value) string {
  if !v.IsValid() {
    return ""
  }
  switch v.Kind() {
  case reflect.String:
    return v.String()
  case reflect.Bool:
    if v.Bool() {
      return "true"
    } else {
      return "false"
    }
  case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
    return fmt.Sprintf("%d", v.Int())
  case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
    return fmt.Sprintf("%u", v.Uint())
  case reflect.Float32, reflect.Float64:
    return fmt.Sprintf("%f", v.Float())
  }
  return fmt.Sprintf("%v", v.Interface())
}

func ToString(v interface{}) string {
  if nil == v {
    return ""
  }
  switch s := v.(type) {
  case string:
    return s
  case []byte:
    return string(s)
  }
  return ToStringByReflect(reflectTarget(reflect.ValueOf(v)))
}
