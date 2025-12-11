// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package encoding

import (
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
)

// The following structs emulate the embedding pattern used for extensions

type IEmbeddedValue any

type Embedded struct {
	IEmbeddedValue `json:"embedded,omitempty"`

	FieldCache map[string]any `field-cache:"" cbor:"-" json:"-"`
}

type MyStruct struct {
	Field0 string `cbor:"0,keyasint,omitempty" json:"field0,omitempty"`
	Field1 int    `cbor:"1,keyasint,omitempty" json:"field1,omitempty"`

	Embedded
}

type MyEmbed struct {
	Foo string `cbor:"-1,keyasint,omitempty" json:"foo,omitempty"`
}

type EmbeddedNoCache struct {
	IEmbeddedValue `json:"embedded,omitempty"`
}

type MyStructNoCache struct {
	Field0 string `cbor:"0,keyasint,omitempty" json:"field0,omitempty"`
	Field1 int    `cbor:"1,keyasint,omitempty" json:"field1,omitempty"`

	EmbeddedNoCache
}

func mustInitEncMode() cbor.EncMode {
	encOpt := cbor.EncOptions{
		Sort:        cbor.SortCoreDeterministic,
		IndefLength: cbor.IndefLengthForbidden,
		TimeTag:     cbor.EncTagRequired,
	}

	em, err := encOpt.EncMode()
	if err != nil {
		panic(err)
	}

	return em
}

func mustInitDecMode() cbor.DecMode {
	decOpt := cbor.DecOptions{
		IndefLength: cbor.IndefLengthAllowed,
	}

	dm, err := decOpt.DecMode()
	if err != nil {
		panic(err)
	}

	return dm
}

func Test_preserve_unknown_embeds_CBOR(t *testing.T) {
	// nolint: gocritic
	data := []byte{
		0xa3, // map(3)

		0x00,       // key: 0
		0x62,       // value: tstr(2)
		0x66, 0x31, // "f1"

		0x01, // key: 1
		0x02, // value: 2

		0x20,             // key: -1
		0x63,             // value: tstr(3)
		0x62, 0x61, 0x72, // "bar"
	}

	em := mustInitEncMode()
	dm := mustInitDecMode()

	// First, a sanity test to make sure that embedded data
	// is preserved when there is a concrete struct to
	// contain it.
	embed := MyEmbed{}
	myStruct := MyStruct{Embedded: Embedded{IEmbeddedValue: &embed, FieldCache: map[string]any{}}}

	err := PopulateStructFromCBOR(dm, data, &myStruct)
	assert.NoError(t, err)

	outData, err := SerializeStructToCBOR(em, &myStruct)
	assert.NoError(t, err)
	assert.Equal(t, data, outData)

	// Now, the same test with IEmbeddedValue not set. This simulates the
	// case where extensions are present in the data but the struct needed
	// to understand them has not been registered.
	myStruct = MyStruct{}

	err = PopulateStructFromCBOR(dm, data, &myStruct)
	assert.NoError(t, err)

	outData, err = SerializeStructToCBOR(em, &myStruct)
	assert.NoError(t, err)
	assert.Equal(t, data, outData)

	// Make sure that, without caching, unknown values are simply ignored
	// without causing errors.
	noCache := MyStructNoCache{}
	err = PopulateStructFromCBOR(dm, data, &noCache)
	assert.NoError(t, err)

	// nolint: gocritic
	expectedNoEmbed := []byte{
		0xa2, // map(2)

		0x00,       // key: 0
		0x62,       // value: tstr(2)
		0x66, 0x31, // "f1"

		0x01, // key: 1
		0x02, // value: 2
	}

	outData, err = SerializeStructToCBOR(em, &noCache)
	assert.NoError(t, err)
	assert.Equal(t, expectedNoEmbed, outData)
}

func Test_preserve_unknown_embeds_JSON(t *testing.T) {
	data := []byte(`{"field0":"f1","field1":2,"foo":"bar"}`)

	// First, a sanity test to make sure that embedded data
	// is presevered when there is a concrete struct to
	// contain it.
	embed := MyEmbed{}
	myStruct := MyStruct{Embedded: Embedded{IEmbeddedValue: &embed, FieldCache: map[string]any{}}}

	err := PopulateStructFromJSON(data, &myStruct)
	assert.NoError(t, err)

	outData, err := SerializeStructToJSON(&myStruct)
	assert.NoError(t, err)
	assert.Equal(t, data, outData)

	// Now, the same test with IEmbeddedValue not set. This simulates the
	// case where extensions are present int the data but the struct needed
	// to understand them has not been registered.
	myStruct = MyStruct{}

	err = PopulateStructFromJSON(data, &myStruct)
	assert.NoError(t, err)

	outData, err = SerializeStructToJSON(&myStruct)
	assert.NoError(t, err)
	assert.Equal(t, data, outData)

	// Make sure that, without caching, unknown values are simply ignored
	// without causing errors.
	noCache := MyStructNoCache{}
	err = PopulateStructFromJSON(data, &noCache)
	assert.NoError(t, err)

	expectedNoCache := []byte(`{"field0":"f1","field1":2}`)

	outData, err = SerializeStructToJSON(&noCache)
	assert.NoError(t, err)
	assert.Equal(t, expectedNoCache, outData)
}
