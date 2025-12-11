// Copyright 2024-2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package encoding

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func SerializeStructToJSON(source any) ([]byte, error) {
	rawMap := newStructFieldsJSON()

	structType := reflect.TypeOf(source)
	structVal := reflect.ValueOf(source)

	if err := doSerializeStructToJSON(rawMap, structType, structVal); err != nil {
		return nil, err
	}

	return rawMap.ToJSON()
}

func doSerializeStructToJSON(
	rawMap *structFieldsJSON,
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

		_, ok := typeField.Tag.Lookup("field-cache")
		if ok {
			if err := addCachedFieldsToMapJSON(valField, rawMap); err != nil {
				return err
			}
		}

		tag, ok := typeField.Tag.Lookup("json")
		if !ok {
			continue
		}

		parts := strings.Split(tag, ",")
		key := parts[0]

		if key == "-" {
			continue // field is not marshaled
		}

		isOmitEmpty := false
		if len(parts) > 1 {
			for _, option := range parts[1:] {
				if option == omitempty {
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

		data, err := json.Marshal(valField.Interface())
		if err != nil {
			return fmt.Errorf("error marshaling field %q: %w",
				typeField.Name,
				err,
			)
		}

		if err := rawMap.Add(key, json.RawMessage(data)); err != nil {
			return err
		}
	}

	for _, emb := range embeds {
		if err := doSerializeStructToJSON(rawMap, emb.Type, emb.Value); err != nil {
			return err
		}
	}

	return nil
}

func PopulateStructFromJSON(data []byte, dest any) error {
	rawMap := newStructFieldsJSON()

	if err := rawMap.FromJSON(data); err != nil {
		return err
	}

	structType := reflect.TypeOf(dest)
	structVal := reflect.ValueOf(dest)

	return doPopulateStructFromJSON(rawMap, structType, structVal)
}

func doPopulateStructFromJSON(
	rawMap *structFieldsJSON,
	structType reflect.Type,
	structVal reflect.Value,
) error {
	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
		structVal = structVal.Elem()
	}

	var embeds []embedded
	var fieldCache reflect.Value

	for i := 0; i < structVal.NumField(); i++ {
		typeField := structType.Field(i)
		valField := structVal.Field(i)

		if collectEmbedded(&typeField, valField, &embeds) {
			continue
		}

		_, ok := typeField.Tag.Lookup("field-cache")
		if ok {
			fieldCache = valField
			continue
		}

		tag, ok := typeField.Tag.Lookup("json")
		if !ok {
			continue
		}

		parts := strings.Split(tag, ",")
		key := parts[0]

		if key == "-" {
			continue // field is not marshaled
		}

		isOmitEmpty := false
		if len(parts) > 1 {
			for _, option := range parts[1:] {
				if option == omitempty {
					isOmitEmpty = true
					break
				}
			}
		}

		rawVal, ok := rawMap.Get(key)
		if !ok {
			if isOmitEmpty {
				continue
			}

			return fmt.Errorf("missing mandatory field %q (%q)",
				typeField.Name, key)
		}

		fieldPtr := valField.Addr().Interface()
		if err := json.Unmarshal(rawVal, fieldPtr); err != nil {
			return fmt.Errorf("error unmarshalling field %q: %w",
				typeField.Name,
				err,
			)
		}

		rawMap.Delete(key)
	}

	for _, emb := range embeds {
		if err := doPopulateStructFromJSON(rawMap, emb.Type, emb.Value); err != nil {
			return err
		}
	}

	// Any remaining contents of rawMap will be added to the field cache,
	// if current struct has one.
	return updateFieldCacheJSON(fieldCache, rawMap)
}

// structFieldsJSON is a specialized implementation of "OrderedMap", where the
// order of the keys is kept track of, and used when serializing the map to
// JSON. While JSON maps do not mandate any particular ordering, and so this
// isn't strictly necessary, it is useful to have a _stable_ serialization
// order for map keys to be compatible with regular Go struct serialization
// behavior. This is also useful for tests/examples that compare encoded
// []byte's.
type structFieldsJSON struct {
	Fields map[string]json.RawMessage
	Keys   []string
}

func newStructFieldsJSON() *structFieldsJSON {
	return &structFieldsJSON{
		Fields: make(map[string]json.RawMessage),
	}
}

func (o structFieldsJSON) Has(key string) bool {
	_, ok := o.Fields[key]
	return ok
}

func (o *structFieldsJSON) Add(key string, val json.RawMessage) error {
	if o.Has(key) {
		return fmt.Errorf("duplicate JSON key: %q", key)
	}

	o.Fields[key] = val
	o.Keys = append(o.Keys, key)

	return nil
}

func (o *structFieldsJSON) Get(key string) (json.RawMessage, bool) {
	val, ok := o.Fields[key]
	return val, ok
}

func (o *structFieldsJSON) Delete(key string) {
	delete(o.Fields, key)

	for i, existing := range o.Keys {
		if existing == key {
			o.Keys = append(o.Keys[:i], o.Keys[i+1:]...)
		}
	}
}

func (o *structFieldsJSON) ToJSON() ([]byte, error) {
	var out bytes.Buffer

	out.Write([]byte("{"))

	first := true
	for _, key := range o.Keys {
		if first {
			first = false
		} else {
			out.Write([]byte(","))
		}
		marshaledKey, err := json.Marshal(key)
		if err != nil {
			return nil, fmt.Errorf("problem marshaling key %s: %w", key, err)
		}
		out.Write(marshaledKey)
		out.Write([]byte(":"))
		out.Write(o.Fields[key])
	}

	out.Write([]byte("}"))

	return out.Bytes(), nil
}

func (o *structFieldsJSON) FromJSON(data []byte) error {
	if err := json.Unmarshal(data, &o.Fields); err != nil {
		return err
	}

	return o.unmarshalKeys(data)
}

func (o *structFieldsJSON) unmarshalKeys(data []byte) error {

	decoder := json.NewDecoder(bytes.NewReader(data))

	token, err := decoder.Token()
	if err != nil {
		return err
	}

	if token != json.Delim('{') {
		return errors.New("expected start of object")
	}

	var keys []string

	for {
		token, err = decoder.Token()
		if err != nil {
			return err
		}

		if token == json.Delim('}') {
			break
		}

		key, ok := token.(string)
		if !ok {
			return fmt.Errorf("expected string, found %T", token)
		}

		keys = append(keys, key)

		if err := skipValue(decoder); err != nil {
			return err
		}
	}

	o.Keys = keys

	return nil
}

var errEndOfStream = errors.New("invalid end of array or object")

func skipValue(decoder *json.Decoder) error {

	token, err := decoder.Token()
	if err != nil {
		return err
	}
	switch token {
	case json.Delim('['), json.Delim('{'):
		for {
			if err := skipValue(decoder); err != nil {
				if err == errEndOfStream {
					break
				}
				return err
			}
		}
	case json.Delim(']'), json.Delim('}'):
		return errEndOfStream
	}
	return nil
}

// TypeAndValue stores a JSON object with two attributes: a string "type"
// and a generic "value" (string) defined by type.  This type is used in
// a few places to implement the choice types that CBOR handles using tags.
type TypeAndValue struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

func (o *TypeAndValue) UnmarshalJSON(data []byte) error {
	var temp struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.Type == "" {
		return errors.New("type not set")
	}

	if len(temp.Value) == 0 {
		return fmt.Errorf("no value provided for %s", temp.Type)
	}

	o.Type = temp.Type
	o.Value = temp.Value

	return nil
}

func addCachedFieldsToMapJSON(cacheField reflect.Value, rawMap *structFieldsJSON) error {
	if !isMapStringAny(cacheField) {
		return errors.New("field-cache does not appear to be a map[string]any")
	}

	if !cacheField.IsValid() || cacheField.IsNil() {
		// field cache was never set, so nothing to do
		return nil
	}

	for _, key := range cacheField.MapKeys() {
		keyText := key.String()

		data, err := json.Marshal(cacheField.MapIndex(key).Interface())
		if err != nil {
			return fmt.Errorf(
				"error marshaling field-cache entry %q: %w",
				keyText,
				err,
			)
		}

		if err := rawMap.Add(keyText, json.RawMessage(data)); err != nil {
			return fmt.Errorf(
				"could not add field-cache entry %q to serialization map: %w",
				keyText,
				err,
			)
		}
	}

	return nil
}

func updateFieldCacheJSON(cacheField reflect.Value, rawMap *structFieldsJSON) error {
	if !cacheField.IsValid() {
		// current struct does not have a field-cache field
		return nil
	}

	if !isMapStringAny(cacheField) {
		return errors.New("field-cache does not appear to be a map[string]any")
	}

	if cacheField.IsNil() {
		cacheField.Set(reflect.MakeMap(cacheField.Type()))
	}

	for key, rawVal := range rawMap.Fields {
		var val any
		if err := json.Unmarshal(rawVal, &val); err != nil {
			return fmt.Errorf("could not unmarshal key %q: %w", key, err)
		}

		keyVal := reflect.ValueOf(key)
		valVal := reflect.ValueOf(val)
		cacheField.SetMapIndex(keyVal, valVal)
	}

	return nil
}
