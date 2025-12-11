// Copyright 2023-2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package extensions

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

var ErrExtensionNotFound = errors.New("extension not found")

type IMapValue any

type ExtensionValue struct {
	CBORTag   string
	JSONTag   string
	FieldName string
	Value     any
}

type Extensions struct {
	IMapValue `json:"extensions,omitempty"`

	Cached map[string]any `field-cache:"" cbor:"-" json:"-"`
}

func (o *Extensions) Register(exts IMapValue) {
	if reflect.TypeOf(exts).Kind() != reflect.Pointer {
		panic("attempting to register a non-pointer IMapValue")
	}

	// Ensure that the values of any existing extensions are preserved.
	// The contents of the existing IMapValue (if there is one) are added
	// to the cache, which is then applied to the new IMapValue. If the new
	// IMapValue has fields corresponding to the old extensions, they will be
	// populated into the new IMapValue; any old extensions that are not
	// recognized by the new IMapValue will be cached.
	updateMapFromInterface(&o.Cached, o.IMapValue)
	o.IMapValue = exts
	updateInterfaceFromMap(o.IMapValue, o.Cached)
}

func (o *Extensions) HaveExtensions() bool {
	return o.IMapValue != nil
}

func (o Extensions) New() IMapValue {
	return newIMapValue(o.IMapValue)
}

func (o *Extensions) Get(name string) (any, error) {
	if o.IMapValue == nil {
		return nil, fmt.Errorf("%w: %s", ErrExtensionNotFound, name)
	}

	extType := reflect.TypeOf(o.IMapValue)
	extVal := reflect.ValueOf(o.IMapValue)
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

func (o *Extensions) IsEmpty() bool {
	if o.IMapValue == nil {
		return true
	}

	extVal := reflect.ValueOf(o.IMapValue)
	if reflect.TypeOf(o.IMapValue).Kind() == reflect.Pointer {
		extVal = extVal.Elem()
	}

	for i := 0; i < extVal.NumField(); i++ {
		if !extVal.Field(i).IsZero() {
			return false
		}
	}

	return true
}

func (o *Extensions) Values() []ExtensionValue {
	var ret []ExtensionValue
	if o.IMapValue == nil {
		return ret
	}

	extType := reflect.TypeOf(o.IMapValue)
	extVal := reflect.ValueOf(o.IMapValue)
	if extType.Kind() == reflect.Pointer {
		extType = extType.Elem()
		extVal = extVal.Elem()
	}

	var fieldJSONTag, fieldCBORTag string
	for i := 0; i < extVal.NumField(); i++ {
		typeField := extType.Field(i)

		tag, ok := typeField.Tag.Lookup("json")
		if ok {
			fieldJSONTag = strings.Split(tag, ",")[0]
		}

		tag, ok = typeField.Tag.Lookup("cbor")
		if ok {
			fieldCBORTag = strings.Split(tag, ",")[0]
		}

		ret = append(ret, ExtensionValue{
			CBORTag:   fieldCBORTag,
			JSONTag:   fieldJSONTag,
			FieldName: typeField.Name,
			Value:     extVal.Field(i).Interface(),
		})
	}

	return ret
}

func (o *Extensions) MustGetString(name string) string {
	v, _ := o.GetString(name)
	return v
}

func (o *Extensions) GetString(name string) (string, error) {
	v, err := o.Get(name)
	if err != nil {
		return "", err
	}

	return cast.ToStringE(v)
}

func (o *Extensions) MustGetInt(name string) int {
	v, _ := o.GetInt(name)
	return v
}

func (o *Extensions) GetInt(name string) (int, error) {
	v, err := o.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToIntE(v)
}

func (o *Extensions) MustGetInt64(name string) int64 {
	v, _ := o.GetInt64(name)
	return v
}

func (o *Extensions) GetInt64(name string) (int64, error) {
	v, err := o.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToInt64E(v)
}

func (o *Extensions) MustGetInt32(name string) int32 {
	v, _ := o.GetInt32(name)
	return v
}

func (o *Extensions) GetInt32(name string) (int32, error) {
	v, err := o.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToInt32E(v)
}

func (o *Extensions) MustGetInt16(name string) int16 {
	v, _ := o.GetInt16(name)
	return v
}

func (o *Extensions) GetInt16(name string) (int16, error) {
	v, err := o.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToInt16E(v)
}

func (o *Extensions) MustGetInt8(name string) int8 {
	v, _ := o.GetInt8(name)
	return v
}

func (o *Extensions) GetInt8(name string) (int8, error) {
	v, err := o.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToInt8E(v)
}

func (o *Extensions) MustGetUint(name string) uint {
	v, _ := o.GetUint(name)
	return v
}

func (o *Extensions) GetUint(name string) (uint, error) {
	v, err := o.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToUintE(v)
}

func (o *Extensions) MustGetUint64(name string) uint64 {
	v, _ := o.GetUint64(name)
	return v
}

func (o *Extensions) GetUint64(name string) (uint64, error) {
	v, err := o.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToUint64E(v)
}

func (o *Extensions) MustGetUint32(name string) uint32 {
	v, _ := o.GetUint32(name)
	return v
}

func (o *Extensions) GetUint32(name string) (uint32, error) {
	v, err := o.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToUint32E(v)
}

func (o *Extensions) MustGetUint16(name string) uint16 {
	v, _ := o.GetUint16(name)
	return v
}

func (o *Extensions) GetUint16(name string) (uint16, error) {
	v, err := o.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToUint16E(v)
}

func (o *Extensions) MustGetUint8(name string) uint8 {
	v, _ := o.GetUint8(name)
	return v
}

func (o *Extensions) GetUint8(name string) (uint8, error) {
	v, err := o.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToUint8E(v)
}

func (o *Extensions) MustGetFloat32(name string) float32 {
	v, _ := o.GetFloat32(name)
	return v
}

func (o *Extensions) GetFloat32(name string) (float32, error) {
	v, err := o.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToFloat32E(v)
}

func (o *Extensions) MustGetFloat64(name string) float64 {
	v, _ := o.GetFloat64(name)
	return v
}

func (o *Extensions) GetFloat64(name string) (float64, error) {
	v, err := o.Get(name)
	if err != nil {
		return 0, err
	}

	return cast.ToFloat64E(v)
}

func (o *Extensions) MustGetBool(name string) bool {
	v, _ := o.GetBool(name)
	return v
}

func (o *Extensions) GetBool(name string) (bool, error) {
	v, err := o.Get(name)
	if err != nil {
		return false, err
	}

	return cast.ToBoolE(v)
}

func (o *Extensions) MustGetSlice(name string) []any {
	v, _ := o.GetSlice(name)
	return v
}

func (o *Extensions) GetSlice(name string) ([]any, error) {
	v, err := o.Get(name)
	if err != nil {
		return []any{}, err
	}

	return cast.ToSliceE(v)
}

func (o *Extensions) MustGetIntSlice(name string) []int {
	v, _ := o.GetIntSlice(name)
	return v
}

func (o *Extensions) GetIntSlice(name string) ([]int, error) {
	v, err := o.Get(name)
	if err != nil {
		return []int{}, err
	}

	return cast.ToIntSliceE(v)
}

func (o *Extensions) MustGetStringSlice(name string) []string {
	v, _ := o.GetStringSlice(name)
	return v
}

func (o *Extensions) GetStringSlice(name string) ([]string, error) {
	v, err := o.Get(name)
	if err != nil {
		return []string{}, err
	}

	return cast.ToStringSliceE(v)
}

func (o *Extensions) MustGetStringMap(name string) map[string]any {
	v, _ := o.GetStringMap(name)
	return v
}

func (o *Extensions) GetStringMap(name string) (map[string]any, error) {
	v, err := o.Get(name)
	if err != nil {
		return map[string]any{}, err
	}

	return cast.ToStringMapE(v)
}

func (o *Extensions) MustGetStringMapString(name string) map[string]string {
	v, _ := o.GetStringMapString(name)
	return v
}

func (o *Extensions) GetStringMapString(name string) (map[string]string, error) {
	v, err := o.Get(name)
	if err != nil {
		return map[string]string{}, err
	}

	return cast.ToStringMapStringE(v)
}

func (o *Extensions) Set(name string, value any) error {
	if o.IMapValue == nil {
		return fmt.Errorf("%w: %s", ErrExtensionNotFound, name)
	}

	extType := reflect.TypeOf(o.IMapValue)
	extVal := reflect.ValueOf(o.IMapValue)
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

func newIMapValue(v IMapValue) IMapValue {
	if v == nil {
		return nil
	}

	valType := reflect.Indirect(reflect.ValueOf(v)).Type()

	return reflect.New(valType).Interface()
}

func updateMapFromInterface(mp *map[string]any, iface any) { // nolint: gocritic
	if iface == nil {
		return
	}

	if *mp == nil {
		*mp = make(map[string]any)
	}

	ifType := reflect.TypeOf(iface)
	ifVal := reflect.ValueOf(iface)
	if ifType.Kind() == reflect.Pointer {
		ifType = ifType.Elem()
		ifVal = ifVal.Elem()
	}

	for i := 0; i < ifVal.NumField(); i++ {
		typeField := ifType.Field(i)
		tag, ok := typeField.Tag.Lookup("cbor")
		if !ok {
			continue
		}

		codePointText := strings.Split(tag, ",")[0]
		valField := ifVal.Field(i)
		if !valField.IsZero() {
			(*mp)[codePointText] = valField.Interface()
		}
	}
}

func updateInterfaceFromMap(iface any, m map[string]any) {
	if iface == nil {
		panic("nil interface")
	}

	ifType := reflect.TypeOf(iface)
	if ifType.Kind() != reflect.Pointer {
		panic("interface must be a pointer")
	}

	ifType = ifType.Elem()
	ifVal := reflect.ValueOf(iface).Elem()

	for i := 0; i < ifVal.NumField(); i++ {
		var fieldJSONTag, fieldCBORTag string
		typeField := ifType.Field(i)
		valField := ifVal.Field(i)

		tag, ok := typeField.Tag.Lookup("json")
		if ok {
			fieldJSONTag = strings.Split(tag, ",")[0]
		}

		tag, ok = typeField.Tag.Lookup("cbor")
		if ok {
			fieldCBORTag = strings.Split(tag, ",")[0]
		}

		mapKey := fieldJSONTag
		rawMapVal, ok := m[mapKey]
		if !ok {
			mapKey = fieldCBORTag
			rawMapVal, ok = m[mapKey]
			if !ok {
				continue
			}
		}

		mapVal := reflect.ValueOf(rawMapVal)
		if !mapVal.Type().AssignableTo(typeField.Type) {
			if mapVal.Type().ConvertibleTo(typeField.Type) {
				mapVal = mapVal.Convert(typeField.Type)
			} else {
				// We cannot return an error here, and we don't
				// want to panic, so we're just going to keep the
				// entry in the cache.
				continue
			}
		}

		valField.Set(mapVal)
		delete(m, mapKey)
	}
}
