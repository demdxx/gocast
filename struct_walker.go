package gocast

import (
	"context"
	"errors"
	"reflect"
)

var (
	ErrWalkSkip = errors.New("skip field walk")
	ErrWalkStop = errors.New("stop field walk")
)

type StructWalkObject interface {
	Parent() StructWalkObject
	RefValue() reflect.Value
	Struct() any
}

type structWalkObject struct {
	parent StructWalkObject
	strct  reflect.Value
}

func (obj *structWalkObject) Parent() StructWalkObject { return obj.parent }
func (obj *structWalkObject) RefValue() reflect.Value  { return obj.strct }
func (obj *structWalkObject) Struct() any              { return obj.strct.Interface() }

// StructWalkField is the type of the field visited by StructWalk
type StructWalkField interface {
	Name() string
	Tag(name string) string
	IsEmpty() bool
	RefValue() reflect.Value
	Value() any
	SetValue(ctx context.Context, v any) error
}

type structWalkField struct {
	name      string
	fieldVal  reflect.Value
	fieldType reflect.StructField
}

func (fl *structWalkField) Name() string {
	return fl.name
}

func (fl *structWalkField) Tag(name string) string {
	return fl.fieldType.Tag.Get(name)
}

func (fl *structWalkField) IsEmpty() bool {
	return IsEmpty(fl.Value())
}

func (fl *structWalkField) RefValue() reflect.Value {
	return fl.fieldVal
}

func (fl *structWalkField) Value() any {
	return fl.fieldVal.Interface()
}

func (fl *structWalkField) SetValue(ctx context.Context, v any) error {
	if !fl.fieldVal.CanSet() {
		return wrapError(ErrStructFieldValueCantBeChanged, fl.name)
	}
	return setFieldValueNoCastSetter(ctx, fl.fieldVal, v, true)
}

// StructWalk walks the struct recursively
func StructWalk(ctx context.Context, v any, walker structWalkerFunc, options ...WalkOption) error {
	structVal := reflectTarget(reflect.ValueOf(v))
	if structVal.Kind() != reflect.Struct {
		return wrapError(ErrUnsupportedSourceType, structVal.Type().Name())
	}
	opt := StructWalkOptions{}
	for _, o := range options {
		o(&opt)
	}
	err := _structWalk(ctx, &structWalkObject{strct: structVal}, walker, &opt)
	return IfThen(errors.Is(err, ErrWalkStop), nil, err)
}

func _structWalk(ctx context.Context, v StructWalkObject, walker structWalkerFunc, opt *StructWalkOptions, path ...string) error {
	var (
		err           error
		structVal     = v.RefValue()
		structValType = structVal.Type()
	)
	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i)
		fieldType := structValType.Field(i)
		if !fieldType.IsExported() {
			continue
		}
		fieldWrapper := structWalkField{
			name:      fieldType.Name,
			fieldVal:  field,
			fieldType: fieldType,
		}
		if err = walker(ctx, v, &fieldWrapper, path); err == nil {
			if stTrg := reflectTarget(field); stTrg.Kind() == reflect.Struct {
				newV := structWalkObject{parent: v, strct: stTrg}
				pathName := opt.PathName(ctx, v, &fieldWrapper, path)
				err = _structWalk(ctx, &newV, walker, opt, append(path, pathName)...)
			}
		}
		if err != nil && !errors.Is(err, ErrWalkSkip) {
			return err
		}
	}
	return nil
}
