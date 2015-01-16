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
)

func To(v, to interface{}) interface{} {
  if nil == v || nil == to {
    return nil
  }
  return ToKind(v, reflect.ValueOf(to).Kind())
}

func ToKind(v interface{}, to reflect.Kind) interface{} {
  if nil == v {
    return nil
  }

  switch to {
  case reflect.String:
    return ToString(v)
  case reflect.Bool:
    return ToBool(v)
  case reflect.Int:
    return ToInt(v)
  case reflect.Int8:
    return (int8)(ToInt(v))
  case reflect.Int16:
    return (int16)(ToInt(v))
  case reflect.Int32:
    return ToInt32(v)
  case reflect.Int64:
    return ToInt64(v)
  case reflect.Uint:
    return ToUint(v)
  case reflect.Uint8:
    return (uint8)(ToUint(v))
  case reflect.Uint16:
    return (uint16)(ToUint(v))
  case reflect.Uint32:
    return ToUint32(v)
  case reflect.Uint64:
    return ToUint64(v)
  case reflect.Float32:
    return ToFloat32(v)
  case reflect.Float64:
    return ToFloat64(v)
  }
  return nil
}
