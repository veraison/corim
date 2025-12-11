// Copyright 2023-2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package extensions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/encoding"
)

type Entity struct {
	EntityName string  `cbor:"0,keyasint" json:"entity-name"`
	Roles      []int64 `cbor:"1,keyasint,omitempty" json:"roles,omitempty"`

	Extensions
}

type TestExtensions struct {
	Address     string            `cbor:"-1,keyasint,omitempty" json:"address,omitempty"`
	Size        int               `cbor:"-2,keyasint,omitempty" json:"size,omitempty"`
	YearsOnAir  float32           `cbor:"-3,keyasint,omitempty" json:"years-on-air,omitempty"`
	StillAiring bool              `cbor:"-4,keyasint,omitempty" json:"still-airing,omitempty"`
	Ages        []int             `cbor:"-5,keyasint,omitempty" json:"ages,omitempty"`
	Jobs        map[string]string `cbor:"-6,keyasint,omitempty" json:"jobs,omitempty"`
}

func TestExtensions_Register(t *testing.T) {
	exts := Extensions{}
	assert.False(t, exts.HaveExtensions())

	exts.Register(&TestExtensions{})
	assert.True(t, exts.HaveExtensions())

	badRegister := func() {
		exts.Register(TestExtensions{})
	}

	assert.Panics(t, badRegister)
}

func TestExtensions_GetSet(t *testing.T) {
	extsVal := TestExtensions{
		Address:     "742 Evergreen Terrace",
		Size:        6,
		YearsOnAir:  33.8,
		StillAiring: true,
		Ages:        []int{2, 7, 8, 10, 37, 38},
		Jobs: map[string]string{
			"Homer": "safety inspector",
			"Marge": "housewife",
			"Bart":  "elementary school student",
			"Lisa":  "elementary school student",
		},
	}
	exts := Extensions{IMapValue: &extsVal}

	v, err := exts.GetInt("size")
	assert.NoError(t, err)
	assert.Equal(t, 6, v)

	assert.Equal(t, 6, exts.MustGetInt("size"))
	assert.Equal(t, int64(6), exts.MustGetInt64("size"))
	assert.Equal(t, int32(6), exts.MustGetInt32("size"))
	assert.Equal(t, int16(6), exts.MustGetInt16("size"))
	assert.Equal(t, int8(6), exts.MustGetInt8("size"))

	assert.Equal(t, uint(6), exts.MustGetUint("size"))
	assert.Equal(t, uint64(6), exts.MustGetUint64("size"))
	assert.Equal(t, uint32(6), exts.MustGetUint32("size"))
	assert.Equal(t, uint16(6), exts.MustGetUint16("size"))
	assert.Equal(t, uint8(6), exts.MustGetUint8("size"))

	assert.InEpsilon(t, float32(33.8), exts.MustGetFloat32("years-on-air"), 0.000001)
	assert.InEpsilon(t, float64(33.8), exts.MustGetFloat64("-3"), 0.000001)

	assert.Equal(t, true, exts.MustGetBool("StillAiring"))

	_, err = exts.GetSlice("ages")
	assert.EqualError(t, err,
		`unable to cast []int{2, 7, 8, 10, 37, 38} of type []int to []interface{}`)
	assert.Nil(t, exts.MustGetSlice("ages"))

	assert.EqualValues(t, []int{2, 7, 8, 10, 37, 38}, exts.MustGetIntSlice("ages"))
	assert.EqualValues(t, []string{"2", "7", "8", "10", "37", "38"},
		exts.MustGetStringSlice("ages"))

	assert.EqualValues(t, map[string]string{
		"Homer": "safety inspector",
		"Marge": "housewife",
		"Bart":  "elementary school student",
		"Lisa":  "elementary school student",
	}, exts.MustGetStringMapString("jobs"))

	_, err = exts.GetStringMap("jobs")
	assert.EqualError(t, err,
		`unable to cast map[string]string{"Bart":"elementary school student", "Homer":"safety inspector", "Lisa":"elementary school student", "Marge":"housewife"} of type map[string]string to map[string]interface{}`)
	m := exts.MustGetStringMap("jobs")
	assert.Equal(t, map[string]any{}, m)

	s, err := exts.GetString("address")
	assert.NoError(t, err)
	assert.Equal(t, "742 Evergreen Terrace", s)

	_, err = exts.GetInt("address")
	assert.EqualError(t, err, `unable to cast "742 Evergreen Terrace" of type string to int`)

	_, err = exts.GetInt("foo")
	assert.EqualError(t, err, "extension not found: foo")

	err = exts.Set("-1", "123 Fake Street")
	assert.NoError(t, err)

	s, err = exts.GetString("address")
	assert.NoError(t, err)
	assert.Equal(t, "123 Fake Street", s)

	err = exts.Set("Size", "foo")
	assert.EqualError(t, err, `cannot set field "Size" (of type int) to foo (string)`)

	assert.Equal(t, "", exts.MustGetString("does-not-exist"))
	assert.Equal(t, 0, exts.MustGetInt("does-not-exist"))
}

func Test_Extensions_New(t *testing.T) {
	exts := Extensions{}

	assert.Nil(t, exts.New())

	exts.Register(&TestExtensions{})

	newValOne, ok := exts.New().(*TestExtensions)
	require.True(t, ok)

	newValTwo, ok := exts.New().(*TestExtensions)
	require.True(t, ok)

	newValOne.Address = "123 Fake Street"
	assert.Equal(t, "", newValTwo.Address)
}

func Test_Extensions_IsEmpty(t *testing.T) {
	exts := Extensions{}
	assert.True(t, exts.IsEmpty())

	exts.Register(&TestExtensions{})
	assert.True(t, exts.IsEmpty())

	exts.Register(&TestExtensions{Address: "123 Fake Street"})
	assert.False(t, exts.IsEmpty())
}

func Test_Extensions_Values(t *testing.T) {
	extStruct := struct {
		Foo string `cbor:"0,keyasint,omitempty" json:"foo,omitempty"`
		Bar int    `cbor:"1,keyasint,omitempty" json:"bar,omitempty"`
	}{"baz", 42}

	exts := Extensions{}
	exts.Register(&extStruct)

	vals := exts.Values()
	assert.Len(t, vals, 2)
	assert.Equal(t, ExtensionValue{CBORTag: "0", JSONTag: "foo", FieldName: "Foo", Value: "baz"}, vals[0])
	assert.Equal(t, ExtensionValue{CBORTag: "1", JSONTag: "bar", FieldName: "Bar", Value: 42}, vals[1])

	exts = Extensions{}
	vals = exts.Values()
	assert.Len(t, vals, 0)
}

func Test_Extensions_unknown_handling_CBOR(t *testing.T) {
	// nolint: gocritic
	data := []byte{
		0xa4, // map(4) [entity]

		0x00,             // key: 0 [entity-name]
		0x63,             // value: tstr(3)
		0x66, 0x6f, 0x6f, // "foo"

		0x1,        // key: 1 [roles]
		0x82,       // value: array(2)
		0x01, 0x02, // [1, 2]

		0x21, // key: -2 [extension(size)]
		0x07, // value: 7

		0x27, // key: -8 [extension(<unknown>)]
		0xf5, // value: true
	}

	entity := Entity{}
	err := encoding.PopulateStructFromCBOR(dm, data, &entity)
	assert.NoError(t, err)
	assert.Equal(t, "foo", entity.EntityName)
	assert.Equal(t, []int64{1, 2}, entity.Roles)
	assert.Equal(t, uint64(7), entity.Extensions.Cached["-2"]) // nolint: staticcheck

	// Check that the cached value has been populated into the
	// newly-registered struct.
	entity.Register(&TestExtensions{})
	assert.Equal(t, 7, entity.MustGetInt("size"))

	// Check that the populated value is no longer cached.
	_, ok := entity.Extensions.Cached["-2"] // nolint: staticcheck
	assert.False(t, ok)

	entity = Entity{}
	entity.Register(&TestExtensions{})
	err = encoding.PopulateStructFromCBOR(dm, data, &entity)
	assert.NoError(t, err)

	// If extensions were registered before unmarshalling, the value gets
	// populated directly into the registered struct, bypassing the cache.
	assert.Equal(t, 7, entity.MustGetInt("size"))
	_, ok = entity.Extensions.Cached["-2"] // nolint: staticcheck
	assert.False(t, ok)

	// Values for keys in in the registered struct still go into cache.
	val, ok := entity.Extensions.Cached["-8"] // nolint: staticcheck
	assert.True(t, ok)
	assert.True(t, val.(bool))

	encoded, err := encoding.SerializeStructToCBOR(em, &entity)
	assert.NoError(t, err)
	assert.Equal(t, data, encoded)
}

func Test_Extensions_unknown_handling_JSON(t *testing.T) {
	data := []byte(`{"entity-name":"foo","roles":[1,2],"size":7,"-8":true}`)

	entity := Entity{}
	err := encoding.PopulateStructFromJSON(data, &entity)
	assert.NoError(t, err)
	assert.Equal(t, "foo", entity.EntityName)
	assert.Equal(t, []int64{1, 2}, entity.Roles)
	assert.Equal(t, float64(7), entity.Extensions.Cached["size"]) // nolint: staticcheck

	// since we only have the JSON field name "size", and we don't know
	// what extension it corresponds to, CBOR encoding fails.
	_, err = encoding.SerializeStructToCBOR(em, &entity)
	assert.ErrorContains(t, err, "cached field name not an integer")

	// Check that the cached value has been populated into the
	// newly-registered struct.
	entity.Register(&TestExtensions{})
	assert.Equal(t, 7, entity.MustGetInt("size"))

	// Check that the populated value is no longer cached.
	_, ok := entity.Extensions.Cached["size"] // nolint: staticcheck
	assert.False(t, ok)

	// "size" has been recognized and removed form cache; we can now
	// serialize it to CBOR as we now know its code point. The only
	// remaining unknown extension has a name that can parse to an integer,
	// so we can use that as the code point for CBOR, and serialization
	// should succeed.
	_, err = encoding.SerializeStructToCBOR(em, &entity)
	assert.NoError(t, err)

	entity = Entity{}
	entity.Register(&TestExtensions{})
	err = encoding.PopulateStructFromJSON(data, &entity)
	assert.NoError(t, err)

	// If extensions were registered before unmarshalling, the value gets
	// populated directly into the registered struct, bypassing the cache.
	assert.Equal(t, 7, entity.MustGetInt("size"))
	_, ok = entity.Extensions.Cached["size"] // nolint: staticcheck
	assert.False(t, ok)

	// Values for keys in in the registered struct still go into cache.
	val, ok := entity.Extensions.Cached["-8"] // nolint: staticcheck
	assert.True(t, ok)
	assert.True(t, val.(bool))

	encoded, err := encoding.SerializeStructToJSON(&entity)
	assert.NoError(t, err)

	assert.JSONEq(t, string(data), string(encoded))
}
