// Copyright 2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package corim

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var ErrExtensionNotFound = errors.New("extension not found")

type IExtensionsValue interface{}

type IEntityValidator interface {
	ValidEntity(*Entity) error
}

type Extensions struct {
	IExtensionsValue
}

func (o *Extensions) ValidEntity(entity *Entity) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IExtensionsValue.(IEntityValidator)
	if ok {
		if err := ev.ValidEntity(entity); err != nil {
			return err
		}
	}

	return nil
}

func (o *Extensions) HaveExtensions() bool {
	return o.IExtensionsValue != nil
}

func (o *Extensions) Get(name string) (any, error) {
	if o.IExtensionsValue == nil {
		return nil, fmt.Errorf("%w: %s", ErrExtensionNotFound, name)
	}

	extType := reflect.TypeOf(o.IExtensionsValue)
	extVal := reflect.ValueOf(o.IExtensionsValue)
	if extType.Kind() == reflect.Pointer {
		extType = extType.Elem()
		extVal = extVal.Elem()
	}

	var fieldName, fieldJSONTag, fieldCBORTag string
	for i := 0; i < extVal.NumField(); i++ {
		typeField := extType.Field(i)
		fieldName = typeField.Name

		tag, ok := typeField.Tag.Lookup("json")
		if ok {
			fieldJSONTag = strings.Split(tag, ",")[0]
		}

		tag, ok = typeField.Tag.Lookup("cbor")
		if ok {
			fieldCBORTag = strings.Split(tag, ",")[0]
		}

		if fieldName == name || fieldJSONTag == name || fieldCBORTag == name {
			return extVal.Field(i).Interface(), nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrExtensionNotFound, name)
}

func (o *Extensions) GetString(name string) (string, error) {
	v, err := o.Get(name)
	if err != nil {
		return "", err
	}

	switch t := v.(type) {
	case string:
		return t, nil
	default:
		return fmt.Sprintf("%v", t), nil
	}
}

func (o *Extensions) GetInt(name string) (int64, error) {
	v, err := o.Get(name)
	if err != nil {
		return 0, err
	}

	val := reflect.ValueOf(v)
	if val.CanInt() {
		return val.Int(), nil
	}

	return 0, fmt.Errorf("%s is not an integer: %v (%T)", name, v, v)
}

func (o *Extensions) Set(name string, value any) error {
	if o.IExtensionsValue == nil {
		return fmt.Errorf("%w: %s", ErrExtensionNotFound, name)
	}

	extType := reflect.TypeOf(o.IExtensionsValue)
	extVal := reflect.ValueOf(o.IExtensionsValue)
	if extType.Kind() == reflect.Pointer {
		extType = extType.Elem()
		extVal = extVal.Elem()
	}

	var fieldName, fieldJSONTag, fieldCBORTag string
	for i := 0; i < extVal.NumField(); i++ {
		typeField := extType.Field(i)
		valField := extVal.Field(i)
		fieldName = typeField.Name

		tag, ok := typeField.Tag.Lookup("json")
		if ok {
			fieldJSONTag = strings.Split(tag, ",")[0]
		}

		tag, ok = typeField.Tag.Lookup("cbor")
		if ok {
			fieldCBORTag = strings.Split(tag, ",")[0]
		}

		if fieldName == name || fieldJSONTag == name || fieldCBORTag == name {
			newVal := reflect.ValueOf(value)
			if newVal.CanConvert(valField.Type()) {
				valField.Set(newVal.Convert(valField.Type()))
				return nil
			}

			return fmt.Errorf(
				"cannot set field %q (of type %s) to %v (%T)",
				name, typeField.Type.Name(),
				value, value,
			)
		}
	}

	return fmt.Errorf("%w: %s", ErrExtensionNotFound, name)
}
