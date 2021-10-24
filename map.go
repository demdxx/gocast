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

// ToMapFrom any Map/Object type
func ToMapFrom(src interface{}, tag string, recursive bool) (map[interface{}]interface{}, error) {
	dst := make(map[interface{}]interface{})
	err := ToMap(dst, src, tag, recursive)
	return dst, err
}

// ToMap cast your Source into the Destination type
// tag defines the tags name in the structure to map the keys
func ToMap(dst, src interface{}, tag string, recursive bool) error {
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
	case map[interface{}]interface{}:
		switch {
		case reflect.Map == t.Kind():
			for _, k := range s.MapKeys() {
				field := s.MapIndex(k)
				if recursive {
					dest[k.Interface()] = mapDestValue(field.Interface(), destType, tag, recursive)
				} else {
					dest[k.Interface()] = field.Interface()
				}
			}
		case reflect.Struct == t.Kind():
			for i := 0; i < s.NumField(); i++ {
				field := s.Field(i)
				if field.CanSet() {
					name, omitempty := fieldName(t.Field(i), tag)
					if len(name) > 0 {
						field := s.Field(i)
						fl := getValue(field.Interface())
						if !omitempty || !IsEmpty(fl) {
							if recursive {
								dest[name] = mapDestValue(fl, destType, tag, recursive)
							} else {
								dest[name] = fl
							}
						} // end if !omitempty || !IsEmpty(fl)
					}
				} // end if
			}
		default:
			err = errUnsupportedSourceType
		}
	case map[string]interface{}:
		switch {
		case reflect.Map == t.Kind():
			for _, k := range s.MapKeys() {
				field := s.MapIndex(k)
				if recursive {
					dest[ToString(k.Interface())] = mapDestValue(field.Interface(), destType, tag, recursive)
				} else {
					dest[ToString(k.Interface())] = field.Interface()
				}
			}
		case reflect.Struct == t.Kind():
			for i := 0; i < s.NumField(); i++ {
				field := s.Field(i)
				if field.CanSet() {
					name, omitempty := fieldName(t.Field(i), tag)
					if len(name) > 0 {
						fl := getValue(field.Interface())
						if !omitempty || !IsEmpty(fl) {
							if recursive {
								dest[name] = mapDestValue(fl, destType, tag, recursive)
							} else {
								dest[name] = fl
							}
						} // end if !omitempty || !IsEmpty(fl)
					}
				} // end if
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
				field := s.Field(i)
				if field.CanSet() {
					name, omitempty := fieldName(t.Field(i), tag)
					if len(name) > 0 {
						fl := getValue(s.Field(i).Interface())
						if !omitempty || !IsEmpty(fl) {
							dest[name] = ToString(fl)
						}
					} // end if
				}
			}
		default:
			err = errUnsupportedSourceType
		}
	default:
		err = errUnsupportedType
	}
	return err
}

// ToSiMap converts input Map/Object type into the map[string]interface{}
func ToSiMap(src interface{}, tag string, recursive bool) (map[string]interface{}, error) {
	dst := make(map[string]interface{})
	err := ToMap(dst, src, tag, recursive)
	return dst, err
}

// ToStringMap converts input Map/Object type into the map[string]string
func ToStringMap(src interface{}, tag string, recursive bool) (map[string]string, error) {
	dst := make(map[string]string)
	err := ToMap(dst, src, tag, recursive)
	return dst, err
}

///////////////////////////////////////////////////////////////////////////////
/// MARK: Helpers
///////////////////////////////////////////////////////////////////////////////

func mapValueByStringKeys(src interface{}, keys []string) interface{} {
	switch m := src.(type) {
	case map[interface{}]interface{}:
		for k, v := range m {
			skey := ToString(k)
			for _, ks := range keys {
				if skey == ks {
					return v
				}
			}
		}
	case map[string]interface{}:
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
					var i interface{} = v
					return i
				}
			}
		}
	}
	return nil
}

func mapDestValue(fl interface{}, destType reflect.Type, tag string, recursive bool) interface{} {
	field := reflect.ValueOf(fl)
	switch field.Kind() {
	case reflect.Slice:
		if field.Len() > 0 {
			switch field.Index(0).Kind() {
			case reflect.Map, reflect.Struct:
				list := make([]interface{}, field.Len())
				for i := 0; i < field.Len(); i++ {
					var v interface{} = reflect.New(destType)
					if ToMap(v, field.Index(i), tag, recursive) == nil {
						list = append(list, v)
					}
				}
				return list
			}
		}
	case reflect.Map, reflect.Struct:
		var v interface{} = reflect.New(destType)
		if ToMap(v, fl, tag, recursive) == nil {
			return v
		}
	}
	return fl
}
