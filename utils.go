package gocast

import (
  "reflect"
)

func reflectTarget(r reflect.Value) reflect.Value {
  for reflect.Ptr == r.Kind() {
    r = r.Elem()
  }
  return r
}
