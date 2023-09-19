// Copyright 2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package encoding

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	cbor "github.com/fxamacker/cbor/v2"
)

type embedded struct {
	Type  reflect.Type
	Value reflect.Value
}

func SerializeStructToCBOR(em cbor.EncMode, source any) ([]byte, error) {
	rawMap := make(map[int]cbor.RawMessage)

	structType := reflect.TypeOf(source)
	structVal := reflect.ValueOf(source)

	if err := doSerializeStructToCBOR(em, rawMap, structType, structVal); err != nil {
		return nil, err
	}

	return em.Marshal(rawMap)
}

func doSerializeStructToCBOR(
	em cbor.EncMode,
	rawMap map[int]cbor.RawMessage,
	structType reflect.Type,
	structVal reflect.Value,
) error {
	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
		structVal = structVal.Elem()
	}

	var embeds []embedded

	for i := 0; i < structVal.NumField(); i++ {
		typeField := structType.Field(i)
		valField := structVal.Field(i)

		if collectEmbedded(&typeField, valField, &embeds) {
			continue
		}

		tag, ok := typeField.Tag.Lookup("cbor")
		if !ok {
			continue
		}

		parts := strings.Split(tag, ",")
		keyString := parts[0]

		isOmitEmpty := false
		if len(parts) > 1 {
			for _, option := range parts[1:] {
				if option == "omitempty" {
					isOmitEmpty = true
					break
				}
			}
		}

		// do not serialize zero values if the corresponding field is
		// omitempty
		if isOmitEmpty && valField.IsZero() {
			continue
		}

		keyInt, err := strconv.Atoi(keyString)
		if err != nil {
			return fmt.Errorf("non-integer cbor key: %s", keyString)
		}

		if _, ok := rawMap[keyInt]; ok {
			return fmt.Errorf("duplicate cbor key: %d", keyInt)
		}

		data, err := em.Marshal(valField.Interface())
		if err != nil {
			return fmt.Errorf("error marshaling field %q: %w",
				typeField.Name,
				err,
			)
		}

		rawMap[keyInt] = cbor.RawMessage(data)
	}

	for _, emb := range embeds {
		if err := doSerializeStructToCBOR(em, rawMap, emb.Type, emb.Value); err != nil {
			return err
		}
	}

	return nil
}

func PopulateStructFromCBOR(dm cbor.DecMode, data []byte, dest any) error {
	var rawMap map[int]cbor.RawMessage

	if err := dm.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	structType := reflect.TypeOf(dest)
	structVal := reflect.ValueOf(dest)

	return doPopulateStructFromCBOR(dm, rawMap, structType, structVal)
}

func doPopulateStructFromCBOR(
	dm cbor.DecMode,
	rawMap map[int]cbor.RawMessage,
	structType reflect.Type,
	structVal reflect.Value,
) error {
	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
		structVal = structVal.Elem()
	}

	var embeds []embedded

	for i := 0; i < structVal.NumField(); i++ {
		typeField := structType.Field(i)
		valField := structVal.Field(i)

		if collectEmbedded(&typeField, valField, &embeds) {
			continue
		}

		tag, ok := typeField.Tag.Lookup("cbor")
		if !ok {
			continue
		}

		parts := strings.Split(tag, ",")
		keyString := parts[0]

		isOmitEmpty := false
		if len(parts) > 1 {
			for _, option := range parts[1:] {
				if option == "omitempty" {
					isOmitEmpty = true
					break
				}
			}
		}

		keyInt, err := strconv.Atoi(keyString)
		if err != nil {
			return fmt.Errorf("non-integer cbor key %s", keyString)
		}

		rawVal, ok := rawMap[keyInt]
		if !ok {
			if isOmitEmpty {
				continue
			}

			return fmt.Errorf("missing mandatory field %q (%d)",
				typeField.Name, keyInt)
		}

		fieldPtr := valField.Addr().Interface()
		if err := dm.Unmarshal(rawVal, fieldPtr); err != nil {
			return fmt.Errorf("error unmarshalling field %q: %w",
				typeField.Name,
				err,
			)
		}

		delete(rawMap, keyInt)
	}

	for _, emb := range embeds {
		if err := doPopulateStructFromCBOR(dm, rawMap, emb.Type, emb.Value); err != nil {
			return err
		}
	}

	return nil
}

// collectEmbedded returns true if the Field is embedded (regardless of
// whether or not it was collected).
func collectEmbedded(
	typeField *reflect.StructField,
	valField reflect.Value,
	embeds *[]embedded,
) bool {
	if typeField.Name == typeField.Type.Name() &&
		(typeField.Type.Kind() == reflect.Struct ||
			typeField.Type.Kind() == reflect.Interface) {

		var fieldType reflect.Type
		var fieldValue reflect.Value

		if typeField.Type.Kind() == reflect.Interface {
			fieldValue = valField.Elem()
			if fieldValue.Kind() == reflect.Invalid {
				// no value underlying the interface
				return true
			}
			// use the interface's underlying value's real type
			fieldType = valField.Elem().Type()
		} else {
			fieldType = typeField.Type
			fieldValue = valField
		}

		*embeds = append(*embeds, embedded{Type: fieldType, Value: fieldValue})
		return true
	}

	return false
}
