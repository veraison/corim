// Copyright 2023-2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package encoding

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"

	cbor "github.com/fxamacker/cbor/v2"
)

func SerializeStructToCBOR(em cbor.EncMode, source any) ([]byte, error) {
	rawMap := newStructFieldsCBOR()

	structType := reflect.TypeOf(source)
	structVal := reflect.ValueOf(source)

	if err := doSerializeStructToCBOR(em, rawMap, structType, structVal); err != nil {
		return nil, err
	}

	return rawMap.ToCBOR(em)
}

func doSerializeStructToCBOR(
	em cbor.EncMode,
	rawMap *structFieldsCBOR,
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
			if err := addCachedFieldsToMapCBOR(em, valField, rawMap); err != nil {
				return err
			}
		}

		tag, ok := typeField.Tag.Lookup("cbor")
		if !ok {
			continue
		}

		parts := strings.Split(tag, ",")
		keyString := parts[0]

		if keyString == "-" {
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

		keyInt, err := strconv.Atoi(keyString)
		if err != nil {
			return fmt.Errorf("non-integer cbor key: %s", keyString)
		}

		data, err := em.Marshal(valField.Interface())
		if err != nil {
			return fmt.Errorf("error marshaling field %q: %w",
				typeField.Name,
				err,
			)
		}

		if err := rawMap.Add(keyInt, cbor.RawMessage(data)); err != nil {
			return err
		}
	}

	for _, emb := range embeds {
		if err := doSerializeStructToCBOR(em, rawMap, emb.Type, emb.Value); err != nil {
			return err
		}
	}

	return nil
}

func PopulateStructFromCBOR(dm cbor.DecMode, data []byte, dest any) error {
	rawMap := newStructFieldsCBOR()

	if err := rawMap.FromCBOR(dm, data); err != nil {
		return err
	}

	structType := reflect.TypeOf(dest)
	structVal := reflect.ValueOf(dest)

	return doPopulateStructFromCBOR(dm, rawMap, structType, structVal)
}

func doPopulateStructFromCBOR(
	dm cbor.DecMode,
	rawMap *structFieldsCBOR,
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

		tag, ok := typeField.Tag.Lookup("cbor")
		if !ok {
			continue
		}

		parts := strings.Split(tag, ",")
		keyString := parts[0]

		if keyString == "-" { // field is not marshaled
			continue
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

		keyInt, err := strconv.Atoi(keyString)
		if err != nil {
			return fmt.Errorf("non-integer cbor key %s", keyString)
		}

		rawVal, ok := rawMap.Get(keyInt)
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

		rawMap.Delete(keyInt)
	}

	for _, emb := range embeds {
		if err := doPopulateStructFromCBOR(dm, rawMap, emb.Type, emb.Value); err != nil {
			return err
		}
	}

	// Any remaining contents of rawMap will be added to the field cache,
	// if current struct has one.
	return updateFieldCacheCBOR(dm, fieldCache, rawMap)
}

// structFieldsCBOR is a specialized implementation of "OrderedMap", where the
// order of the keys is kept track of, and used when serializing the map to
// CBOR. While CBOR maps do not mandate any particular ordering, and so this
// isn't strictly necessary, it is useful to have a _stable_ serialization
// order for map keys to be compatible with regular Go struct serialization
// behavior. This is also useful for tests/examples that compare encoded
// []byte's.
type structFieldsCBOR struct {
	Fields map[int]cbor.RawMessage
	Keys   []int
}

func newStructFieldsCBOR() *structFieldsCBOR {
	return &structFieldsCBOR{
		Fields: make(map[int]cbor.RawMessage),
	}
}

func (o structFieldsCBOR) Has(key int) bool {
	_, ok := o.Fields[key]
	return ok
}

func (o *structFieldsCBOR) Add(key int, val cbor.RawMessage) error {
	if o.Has(key) {
		return fmt.Errorf("duplicate cbor key: %d", key)
	}

	o.Fields[key] = val
	o.Keys = append(o.Keys, key)

	return nil
}

func (o *structFieldsCBOR) Get(key int) (cbor.RawMessage, bool) {
	val, ok := o.Fields[key]
	return val, ok
}

func (o *structFieldsCBOR) Delete(key int) {
	delete(o.Fields, key)

	for i, existing := range o.Keys {
		if existing == key {
			o.Keys = append(o.Keys[:i], o.Keys[i+1:]...)
		}
	}
}

func (o *structFieldsCBOR) ToCBOR(em cbor.EncMode) ([]byte, error) {
	var out []byte

	header := byte(0xa0) // 0b101_00000 -- Major Type 5 ==  map
	mapLen := len(o.Keys)
	if mapLen == 0 {
		return []byte{header}, nil
	} else if mapLen < 24 {
		header |= byte(mapLen)
		out = append(out, header)
	} else if mapLen <= math.MaxUint8 {
		header |= byte(24)
		out = append(out, header, uint8(mapLen))
	} else if mapLen <= math.MaxUint16 {
		header |= byte(25)
		out = append(out, header)
		out = binary.BigEndian.AppendUint16(out, uint16(mapLen))
	} else if mapLen <= math.MaxUint32 {
		header |= byte(26)
		out = append(out, header)
		out = binary.BigEndian.AppendUint32(out, uint32(mapLen))
	} else {
		return nil, errors.New("mapLen cannot exceed math.MaxUint32")
	}

	lexSort(em, o.Keys)
	for _, key := range o.Keys {
		marshalledKey, err := em.Marshal(key)
		if err != nil {
			return nil, fmt.Errorf("problem marshaling key %d: %w", key, err)
		}

		out = append(out, marshalledKey...)
		out = append(out, o.Fields[key]...)
	}

	return out, nil
}

func (o *structFieldsCBOR) FromCBOR(dm cbor.DecMode, data []byte) error {
	if len(data) == 0 {
		return errors.New("empty input")
	}

	header := data[0]
	rest := data[1:]
	additionalInfo := 0x1f & header

	var err error

	majorType := (0xe0 & header) >> 5
	if majorType == 6 { // tag
		_, rest, err = processAdditionalInfo(additionalInfo, rest)
		if err != nil {
			return err
		}

		header = rest[0]
		rest = rest[1:]
		majorType = (0xe0 & header) >> 5
		additionalInfo = 0x1f & header
	}

	if majorType != 5 {
		return fmt.Errorf("expected map (CBOR Major Type 5), found Major Type %d", majorType)
	}

	var mapLen int

	mapLen, rest, err = processAdditionalInfo(additionalInfo, rest)
	if err != nil {
		return err
	}

	if mapLen != 0 {
		o.Fields = make(map[int]cbor.RawMessage, mapLen)

		for i := 0; i < mapLen; i++ {
			rest, err = o.unmarshalKeyValue(dm, rest)
			if err != nil {
				return fmt.Errorf("map item %d: %w", i, err)
			}
		}
	} else { // mapLen == 0 --> indefinite encoding
		o.Fields = make(map[int]cbor.RawMessage)

		i := 0
		done := false
		for len(rest) > 0 {
			if rest[0] == 0xFF {
				done = true
				break
			}

			rest, err = o.unmarshalKeyValue(dm, rest)
			if err != nil {
				return fmt.Errorf("map item %d: %w", i, err)
			}

			i++
		}

		if !done {
			return errors.New("unexpected EOF")
		}
	}

	return nil
}

func (o *structFieldsCBOR) unmarshalKeyValue(dm cbor.DecMode, rest []byte) ([]byte, error) {
	var key int
	var val cbor.RawMessage
	var err error

	rest, err = dm.UnmarshalFirst(rest, &key)
	if err != nil {
		return rest, fmt.Errorf("could not unmarshal key: %w", err)
	}

	rest, err = dm.UnmarshalFirst(rest, &val)
	if err != nil {
		return rest, fmt.Errorf("could not unmarshal value: %w", err)
	}

	if err := o.Add(key, val); err != nil {
		return rest, err
	}

	return rest, nil
}

func processAdditionalInfo(
	additionalInfo byte,
	data []byte,
) (mapLen int, rest []byte, err error) {
	rest = data

	if additionalInfo < 24 {
		mapLen = int(additionalInfo)
	} else if additionalInfo < 28 {
		switch additionalInfo - 23 {
		case 1:
			if len(data) < 1 {
				return 0, nil, errors.New("unexpected EOF")
			}
			mapLen = int(data[0])
			rest = data[1:]
		case 2:
			if len(data) < 2 {
				return 0, nil, errors.New("unexpected EOF")
			}
			mapLen = int(binary.BigEndian.Uint16(data[:2]))
			rest = data[2:]
		case 3:
			if len(data) < 4 {
				return 0, nil, errors.New("unexpected EOF")
			}
			mapLen = int(binary.BigEndian.Uint32(data[:4]))
			rest = data[4:]
		default:
			return 0, nil, errors.New("cbor: cannot decode length value of 8 bytes")
		}
	} else if additionalInfo == 31 {
		mapLen = 0 // indefinite encoding
	} else {
		return 0, nil, fmt.Errorf("cbor: unexpected additional information value %d", additionalInfo)
	}

	return mapLen, rest, nil
}

func addCachedFieldsToMapCBOR(em cbor.EncMode, cacheField reflect.Value, rawMap *structFieldsCBOR) error {
	if !isMapStringAny(cacheField) {
		return errors.New("field-cache does not appear to be a map[string]any")
	}

	if !cacheField.IsValid() || cacheField.IsNil() {
		// field cache was never set, so nothing to do
		return nil
	}

	for _, key := range cacheField.MapKeys() {
		keyText := key.String()
		keyInt, err := strconv.Atoi(keyText)
		if err != nil {
			return fmt.Errorf(
				"cached field name not an integer (cannot encode to CBOR): %s",
				keyText,
			)
		}

		data, err := em.Marshal(cacheField.MapIndex(key).Interface())
		if err != nil {
			return fmt.Errorf(
				"error marshaling field-cache entry %q: %w",
				keyText,
				err,
			)
		}

		if err := rawMap.Add(keyInt, cbor.RawMessage(data)); err != nil {
			return fmt.Errorf(
				"could not add field-cache entry %q to serialization map: %w",
				keyText,
				err,
			)
		}
	}

	return nil
}

func updateFieldCacheCBOR(dm cbor.DecMode, cacheField reflect.Value, rawMap *structFieldsCBOR) error {
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
		if err := dm.Unmarshal(rawVal, &val); err != nil {
			return fmt.Errorf("could not unmarshal key %d: %w", key, err)
		}

		keyText := fmt.Sprint(key)
		keyVal := reflect.ValueOf(keyText)
		valVal := reflect.ValueOf(val)
		cacheField.SetMapIndex(keyVal, valVal)
	}

	return nil
}

// Lexicographic sorting of CBOR integer keys. See:
// https://www.ietf.org/archive/id/draft-ietf-cbor-cde-13.html#name-the-lexicographic-map-sorti
func lexSort(em cbor.EncMode, v []int) {
	sort.Slice(v, func(i, j int) bool {
		a, err := em.Marshal(v[i])
		if err != nil {
			panic(err) // integer encoding cannot fail
		}

		b, err := em.Marshal(v[j])
		if err != nil {
			panic(err) // integer encoding cannot fail
		}

		for k, v := range a {
			if v < b[k] {
				return true
			} else if v > b[k] {
				return false
			}
		}

		return false
	})
}
