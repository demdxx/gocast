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

// TryMapCopy converts source into destination or return error
func TryMapCopy[K comparable, V any](dst map[K]V, src any, recursive bool, tags ...string) error {
	if dst == nil || src == nil {
		return ErrInvalidParams
	}
	var (
		srcVal  = reflectTarget(reflect.ValueOf(src))
		srcType = srcVal.Type()
	)
	switch srcType.Kind() {
	case reflect.Map:
		for _, k := range srcVal.MapKeys() {
			field := srcVal.MapIndex(k)
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
	case reflect.Struct:
		for i := 0; i < srcVal.NumField(); i++ {
			name, omitempty := fieldNameFromTags(srcType.Field(i), tags...)
			if len(name) > 0 {
				key, err := TryCast[K](name)
				if err != nil {
					return err
				}
				field := srcVal.Field(i)
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
		return wrapError(ErrUnsupportedSourceType, srcType.String())
	}
	return nil
}

// ToMap cast your Source into the Destination type
// tag defines the tags name in the structure to map the keys
func ToMap(dst, src any, recursive bool, tags ...string) error {
	if dst == nil || src == nil {
		return ErrInvalidParams
	}

	var (
		err      error
		destVal  = reflectTarget(reflect.ValueOf(dst))
		destType = destVal.Type()
		srcVal   = reflectTarget(reflect.ValueOf(src))
		srcType  = srcVal.Type()
	)

	if dst = destVal.Interface(); dst == nil {
		dst = reflect.MakeMap(destType).Interface()
		destVal = reflect.ValueOf(dst)
	}

	switch dest := dst.(type) {
	case map[any]any:
		switch srcType.Kind() {
		case reflect.Map:
			for _, k := range srcVal.MapKeys() {
				field := srcVal.MapIndex(k)
				if recursive {
					dest[k.Interface()], err = mapDestValue(field.Interface(), destType, recursive, tags...)
					if err != nil {
						return err
					}
				} else {
					dest[k.Interface()] = field.Interface()
				}
			}
		case reflect.Struct:
			for i := 0; i < srcVal.NumField(); i++ {
				name, omitempty := fieldNameFromTags(srcType.Field(i), tags...)
				if len(name) > 0 {
					field := srcVal.Field(i)
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
			err = wrapError(ErrUnsupportedSourceType, srcType.String())
		}
	case map[string]any:
		switch srcType.Kind() {
		case reflect.Map:
			for _, k := range srcVal.MapKeys() {
				field := srcVal.MapIndex(k)
				if recursive {
					dest[Str(k.Interface())], err = mapDestValue(field.Interface(), destType, recursive, tags...)
					if err != nil {
						return err
					}
				} else {
					dest[Str(k.Interface())] = field.Interface()
				}
			}
		case reflect.Struct:
			for i := 0; i < srcVal.NumField(); i++ {
				name, omitempty := fieldNameFromTags(srcType.Field(i), tags...)
				if len(name) > 0 {
					field := srcVal.Field(i)
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
			err = wrapError(ErrUnsupportedSourceType, srcType.String())
		}
	case map[string]string:
		switch srcType.Kind() {
		case reflect.Map:
			for _, k := range srcVal.MapKeys() {
				dest[Str(k.Interface())] = Str(srcVal.MapIndex(k).Interface())
			}
		case reflect.Struct:
			for i := 0; i < srcVal.NumField(); i++ {
				name, omitempty := fieldNameFromTags(srcType.Field(i), tags...)
				if len(name) > 0 {
					fl := getValue(srcVal.Field(i).Interface())
					if !omitempty || !IsEmpty(fl) {
						dest[name] = Str(fl)
					}
				} // end if
			}
		default:
			err = wrapError(ErrUnsupportedSourceType, srcType.String())
		}
	default:
		switch destType.Kind() {
		case reflect.Map, reflect.Struct:
			keyType := destType.Key()
			elemType := destType.Elem()
			switch srcType.Kind() {
			case reflect.Map:
				for _, k := range srcVal.MapKeys() {
					keyVal, err := ReflectTryToType(k, keyType, false)
					if err != nil {
						return wrapError(err, Str(k.Interface()))
					}
					mapVal := reflectTarget(srcVal.MapIndex(k))
					val, err := ReflectTryToType(mapVal, elemType, false)
					if err != nil {
						return wrapError(err, "`"+Str(k.Interface())+"` value")
					}
					destVal.SetMapIndex(reflect.ValueOf(keyVal), reflect.ValueOf(val))
				}
			case reflect.Struct:
				for i := 0; i < srcVal.NumField(); i++ {
					name, omitempty := fieldNameFromTags(srcType.Field(i), tags...)
					if len(name) > 0 {
						flVal := reflectTarget(srcVal.Field(i))
						fl := getValue(flVal.Interface())
						if !omitempty || !IsEmpty(fl) {
							keyVal, err := TryToType(name, keyType)
							if err != nil {
								return wrapError(err, name)
							}
							val, err := ReflectTryToType(flVal, elemType, false)
							if err != nil {
								return wrapError(err, "`"+name+"` value")
							}
							destVal.SetMapIndex(reflect.ValueOf(keyVal), reflect.ValueOf(val))
						}
					} // end if
				}
			default:
				err = wrapError(ErrUnsupportedSourceType, srcType.String())
			}
		default:
			err = wrapError(ErrUnsupportedType, destType.String())
		}
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

func reflectMapValueByStringKeys(src reflect.Value, keys []string) any {
	mKeys := src.MapKeys()
	for _, key := range keys {
		for _, mKey := range mKeys {
			if Str(mKey.Interface()) == key {
				return src.MapIndex(mKey).Interface()
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
