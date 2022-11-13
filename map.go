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

	"github.com/pkg/errors"
)

// TryMapCopy converts source into destination or return error
func TryMapCopy[K comparable, V any](dst map[K]V, src any, recursive bool, tags ...string) error {
	if dst == nil || src == nil {
		return errInvalidParams
	}
	var (
		s = reflectTarget(reflect.ValueOf(src))
		t = s.Type()
	)
	switch {
	case reflect.Map == t.Kind():
		for _, k := range s.MapKeys() {
			field := s.MapIndex(k)
			key, err := TryCast[K](k.Interface())
			if err == nil {
				if recursive {
					dst[key], err = TryCastRecursive[V](field.Interface(), tags...)
				} else {
					dst[key], err = TryCast[V](field.Interface(), tags...)
				}
			}
			if err != nil {
				return err
			}
		}
	case reflect.Struct == t.Kind():
		for i := 0; i < s.NumField(); i++ {
			name, omitempty := fieldNameFromTags(t.Field(i), tags...)
			if len(name) > 0 {
				key, err := TryCast[K](name)
				if err != nil {
					return err
				}
				field := s.Field(i)
				fl := getValue(field.Interface())
				if !omitempty || !IsEmpty(fl) {
					if recursive {
						dst[key], err = TryCastRecursive[V](fl, tags...)
					} else {
						dst[key], err = TryCast[V](fl, tags...)
					}
					if err != nil {
						return err
					}
				} // end if !omitempty || !IsEmpty(fl)
			}
		}
	default:
		return errors.Wrap(errUnsupportedSourceType, t.String())
	}
	return nil
}

// ToMap cast your Source into the Destination type
// tag defines the tags name in the structure to map the keys
func ToMap(dst, src any, recursive bool, tags ...string) error {
	if dst == nil || src == nil {
		return errInvalidParams
	}

	var (
		err      error
		destType = reflect.TypeOf(dst)
		s        = reflectTarget(reflect.ValueOf(src))
		t        = s.Type()
	)
	dst = reflectTarget(reflect.ValueOf(dst)).Interface()

	switch dest := dst.(type) {
	case map[any]any:
		switch {
		case reflect.Map == t.Kind():
			for _, k := range s.MapKeys() {
				field := s.MapIndex(k)
				if recursive {
					dest[k.Interface()], err = mapDestValue(field.Interface(), destType, recursive, tags...)
					if err != nil {
						return err
					}
				} else {
					dest[k.Interface()] = field.Interface()
				}
			}
		case reflect.Struct == t.Kind():
			for i := 0; i < s.NumField(); i++ {
				name, omitempty := fieldNameFromTags(t.Field(i), tags...)
				if len(name) > 0 {
					field := s.Field(i)
					fl := getValue(field.Interface())
					if !omitempty || !IsEmpty(fl) {
						if recursive {
							dest[name], err = mapDestValue(fl, destType, recursive, tags...)
							if err != nil {
								return err
							}
						} else {
							dest[name] = fl
						}
					} // end if !omitempty || !IsEmpty(fl)
				}
			}
		default:
			err = errUnsupportedSourceType
		}
	case map[string]any:
		switch {
		case reflect.Map == t.Kind():
			for _, k := range s.MapKeys() {
				field := s.MapIndex(k)
				if recursive {
					dest[ToString(k.Interface())], err = mapDestValue(field.Interface(), destType, recursive, tags...)
					if err != nil {
						return err
					}
				} else {
					dest[ToString(k.Interface())] = field.Interface()
				}
			}
		case reflect.Struct == t.Kind():
			for i := 0; i < s.NumField(); i++ {
				name, omitempty := fieldNameFromTags(t.Field(i), tags...)
				if len(name) > 0 {
					field := s.Field(i)
					fl := getValue(field.Interface())
					if !omitempty || !IsEmpty(fl) {
						if recursive {
							dest[name], err = mapDestValue(fl, destType, recursive, tags...)
							if err != nil {
								return err
							}
						} else {
							dest[name] = fl
						}
					} // end if !omitempty || !IsEmpty(fl)
				}
			}
		default:
			err = errUnsupportedSourceType
		}
	case map[string]string:
		switch {
		case reflect.Map == t.Kind():
			for _, k := range s.MapKeys() {
				dest[ToString(k.Interface())] = ToString(s.MapIndex(k).Interface())
			}
		case reflect.Struct == t.Kind():
			for i := 0; i < s.NumField(); i++ {
				name, omitempty := fieldNameFromTags(t.Field(i), tags...)
				if len(name) > 0 {
					fl := getValue(s.Field(i).Interface())
					if !omitempty || !IsEmpty(fl) {
						dest[name] = ToString(fl)
					}
				} // end if
			}
		default:
			err = errUnsupportedSourceType
		}
	default:
		err = errUnsupportedType
	}
	return err
}

// TryMapFrom source creates new map to convert
func TryMapFrom[K comparable, V any](src any, recursive bool, tags ...string) (map[K]V, error) {
	dst := make(map[K]V)
	err := TryMapCopy(dst, src, recursive, tags...)
	return dst, err
}

// TryMapRecursive creates new map to convert from soruce type with recursive field processing
func TryMapRecursive[K comparable, V any](src any, tags ...string) (map[K]V, error) {
	return TryMapFrom[K, V](src, true, tags...)
}

// TryMap creates new map to convert from soruce type
func TryMap[K comparable, V any](src any, tags ...string) (map[K]V, error) {
	return TryMapFrom[K, V](src, false, tags...)
}

// MapRecursive creates map from source or returns nil
func MapRecursive[K comparable, V any](src any, tags ...string) map[K]V {
	m, _ := TryMapRecursive[K, V](src, tags...)
	return m
}

// Map creates map from source or returns nil
func Map[K comparable, V any](src any, tags ...string) map[K]V {
	m, _ := TryMap[K, V](src, tags...)
	return m
}

// ToMapFrom any Map/Object type
func ToMapFrom(src any, recursive bool, tags ...string) (map[any]any, error) {
	dst := make(map[any]any)
	err := ToMap(dst, src, recursive, tags...)
	return dst, err
}

// ToSiMap converts input Map/Object type into the map[string]any
//
// Deprecated: Use TryMapFrom[string, any](...) instead
func ToSiMap(src any, recursive bool, tags ...string) (map[string]any, error) {
	return TryMapFrom[string, any](src, recursive, tags...)
}

// ToStringMap converts input Map/Object type into the map[string]string
//
// Deprecated: Use TryMapFrom[string, string](...) instead
func ToStringMap(src any, recursive bool, tags ...string) (map[string]string, error) {
	return TryMapFrom[string, string](src, recursive, tags...)
}

///////////////////////////////////////////////////////////////////////////////
/// MARK: Helpers
///////////////////////////////////////////////////////////////////////////////

func mapValueByStringKeys(src any, keys []string) any {
	switch m := src.(type) {
	case map[any]any:
		for k, v := range m {
			skey := Str(k)
			for _, ks := range keys {
				if skey == ks {
					return v
				}
			}
		}
	case map[string]any:
		for k, v := range m {
			for _, ks := range keys {
				if k == ks {
					return v
				}
			}
		}
	case map[string]string:
		for k, v := range m {
			for _, ks := range keys {
				if k == ks {
					var i any = v
					return i
				}
			}
		}
	}
	return nil
}

func mapDestValue(fl any, destType reflect.Type, recursive bool, tags ...string) (any, error) {
	field := reflect.ValueOf(fl)
	switch field.Kind() {
	case reflect.Slice:
		if field.Len() > 0 {
			switch field.Index(0).Kind() {
			case reflect.Map, reflect.Struct:
				list := make([]any, field.Len())
				for i := 0; i < field.Len(); i++ {
					var v any = reflect.New(destType)
					if err := ToMap(v, field.Index(i), recursive, tags...); err != nil {
						return nil, err
					}
					list = append(list, v)
				}
				return list, nil
			}
		}
	case reflect.Map, reflect.Struct:
		var v any = reflect.New(destType)
		if err := ToMap(v, fl, recursive, tags...); err != nil {
			return nil, err
		}
		return v, nil
	}
	return fl, nil
}
