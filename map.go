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

func ToMapFrom(src interface{}, tag string) (map[interface{}]interface{}, error) {
  dst := make(map[interface{}]interface{})
  err := ToMap(dst, src, tag)
  return dst, err
}

func ToMap(dst, src interface{}, tag string) (err error) {
  if nil == dst || nil == src {
    err = errInvalidParams
  } else {
    dst = reflectTarget(reflect.ValueOf(dst)).Interface()
    s := reflectTarget(reflect.ValueOf(src))
    t := s.Type()

    switch dst.(type) {
    case map[interface{}]interface{}:
      if reflect.Map == t.Kind() {
        for _, k := range s.MapKeys() {
          dst.(map[interface{}]interface{})[k.Interface()] = s.MapIndex(k).Interface()
        }
      } else {
        for i := 0; i < s.NumField(); i++ {
          dst.(map[interface{}]interface{})[fieldName(t.Field(i), tag)] = s.Field(i).Interface()
        }
      }
      break
    case map[string]interface{}:
      if reflect.Map == t.Kind() {
        for _, k := range s.MapKeys() {
          dst.(map[string]interface{})[ToString(k.Interface())] = s.MapIndex(k).Interface()
        }
      } else {
        for i := 0; i < s.NumField(); i++ {
          dst.(map[string]interface{})[fieldName(t.Field(i), tag)] = s.Field(i).Interface()
        }
      }
      break
    case map[string]string:
      if reflect.Map == t.Kind() {
        for _, k := range s.MapKeys() {
          dst.(map[string]interface{})[ToString(k.Interface())] = ToString(s.MapIndex(k).Interface())
        }
      } else {
        for i := 0; i < s.NumField(); i++ {
          dst.(map[string]string)[fieldName(t.Field(i), tag)] = ToString(s.Field(i).Interface())
        }
      }
      break
    default:
      err = errUnsupportedType
    }
  }
  return
}

func ToSiMap(src interface{}, tag string) (map[string]interface{}, error) {
  dst := make(map[string]interface{})
  err := ToMap(dst, src, tag)
  return dst, err
}

func ToStringMap(src interface{}, tag string) (map[string]string, error) {
  dst := make(map[string]string)
  err := ToMap(dst, src, tag)
  return dst, err
}

///////////////////////////////////////////////////////////////////////////////
/// MARK: Helpers
///////////////////////////////////////////////////////////////////////////////

func mapValueByKeys(src interface{}, keys []interface{}) interface{} {
  if nil == src || nil == keys {
    return nil
  }
  sKeys := make([]string, len(keys))
  for i, v := range keys {
    sKeys[i] = ToString(v)
  }
  return mapValueByStringKeys(src, sKeys)
}

func mapValueByStringKeys(src interface{}, keys []string) interface{} {
  switch src.(type) {
  case map[interface{}]interface{}:
    for k, v := range src.(map[interface{}]interface{}) {
      skey := ToString(k)
      for _, ks := range keys {
        if skey == ks {
          return v
        }
      }
    }
    break
  case map[string]interface{}:
    for k, v := range src.(map[string]interface{}) {
      for _, ks := range keys {
        if k == ks {
          return v
        }
      }
    }
    break
  case map[string]string:
    for k, v := range src.(map[string]string) {
      for _, ks := range keys {
        if k == ks {
          var i interface{} = v
          return i
        }
      }
    }
    break
  }
  return nil
}
