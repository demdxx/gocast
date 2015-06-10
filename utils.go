package gocast

import (
  "database/sql/driver"
  "reflect"
)

func reflectTarget(r reflect.Value) reflect.Value {
  for reflect.Ptr == r.Kind() {
    r = r.Elem()
  }
  return r
}

func isEmpty(v interface{}) bool {
  if nil == v {
    return true
  }
  t := reflect.ValueOf(v)
  switch t.Kind() {
  case reflect.Slice, reflect.Map, reflect.String:
    return t.Len() < 1
  case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
    return 0 == t.Int()
  case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
    return 0 == t.Uint()
  case reflect.Float32, reflect.Float64:
    return 0 == t.Float()
  }
  return false
}

func getValue(v interface{}) interface{} {
  if nil != v {
    if vl, ok := v.(driver.Valuer); ok {
      v, _ = vl.Value()
    }
  }
  return v
}
