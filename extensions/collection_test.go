// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package extensions

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/encoding"
)

type testExtensible struct {
	FieldOne string `cbor:"-1,keyasint,omitempty"  json:"field-one,omitempty"`
	Error    error  `cbor:"-"  json:"-"`

	Extensions
}

const testPoint Point = "test"

func (o *testExtensible) RegisterExtensions(exts Map) error {
	for p, v := range exts {
		switch p {
		case testPoint:
			o.Extensions.Register(v)
		default:
			return fmt.Errorf("%w: %q", ErrUnexpectedPoint, p)
		}
	}

	return nil
}

func (o testExtensible) GetExtensions() IMapValue {
	return o.Extensions.IMapValue
}

func (o testExtensible) Valid() error {
	return o.Error
}

func (o *testExtensible) UnmarshalCBOR(data []byte) error {
	return encoding.PopulateStructFromCBOR(dm, data, o)
}

func (o testExtensible) MarshalCBOR() ([]byte, error) {
	return encoding.SerializeStructToCBOR(em, o)
}

func (o *testExtensible) UnmarshalJSON(data []byte) error {
	return encoding.PopulateStructFromJSON(data, o)
}

func (o testExtensible) MarshalJSON() ([]byte, error) {
	return encoding.SerializeStructToJSON(o)
}

type testExtensions struct {
	FieldTwo string `cbor:"-2,keyasint,omitempty"  json:"field-two,omitempty"`
}

func Test_Collection_JSON(t *testing.T) {
	testData := []byte(`
	[
		{
			"field-one": "foo",
			"field-two": "bar"
		},
		{
			"field-one": "buzz"
		}
	]
	`)

	c := NewCollection[testExtensible]()

	err := json.Unmarshal(testData, c)
	require.NoError(t, err)

	// as we've not registred the extensions, the unknown field will be ignored.
	decoded, err := c.Values[0].Extensions.GetString("field-two")
	assert.Equal(t, "", decoded)
	assert.EqualError(t, err, "extension not found: field-two")

	err = c.RegisterExtensions(Map{testPoint: &testExtensions{}})
	require.NoError(t, err)
	c.Clear() // clear only clears the values array, it does not "unregister" extensions.

	err = json.Unmarshal(testData, c)
	require.NoError(t, err)

	decoded, err = c.Values[0].Extensions.GetString("field-two")
	require.NoError(t, err)
	assert.Equal(t, "bar", decoded)

	data, err := json.Marshal(c)
	require.NoError(t, err)
	assert.JSONEq(t, string(testData), string(data))
}

func Test_Collection_CBOR(t *testing.T) {
	testData := []byte{
		0x82, 0xa2, 0x20, 0x63, 0x66, 0x6f, 0x6f, 0x21, 0x63, 0x62,
		0x61, 0x72, 0xa1, 0x20, 0x64, 0x62, 0x75, 0x7a, 0x7a,
	}

	c := NewCollection[testExtensible]()

	err := dm.Unmarshal(testData, c)
	require.NoError(t, err)

	// as we've not registred the extensions, the unknown field will be ignored.
	decoded, err := c.Values[0].Extensions.GetString("field-two")
	assert.Equal(t, "", decoded)
	assert.EqualError(t, err, "extension not found: field-two")

	err = c.RegisterExtensions(Map{testPoint: &testExtensions{}})
	require.NoError(t, err)
	c.Clear() // clear only clears the values array, it does not "unregister" extensions.

	err = dm.Unmarshal(testData, c)
	require.NoError(t, err)

	decoded, err = c.Values[0].Extensions.GetString("field-two")
	require.NoError(t, err)
	assert.Equal(t, "bar", decoded)

	data, err := c.MarshalCBOR()
	require.NoError(t, err)
	assert.Equal(t, testData, data)
}

func Test_Collection_Add(t *testing.T) {
	c := NewCollection[testExtensible]()

	assert.True(t, c.IsEmpty())
	assert.Len(t, c.Values, 0)

	e := testExtensible{}
	c.Add(&e)

	assert.False(t, c.IsEmpty())
	assert.Len(t, c.Values, 1)
	assert.Equal(t, e, c.Values[0])
}

func Test_Collection_Valid(t *testing.T) {
	c := NewCollection[testExtensible]()
	assert.NoError(t, c.Valid())

	c.Add(&testExtensible{})
	assert.NoError(t, c.Valid())

	c.Add(&testExtensible{Error: errors.New("test error")})
	assert.EqualError(t, c.Valid(), "error at index 1: test error")
}

func Test_Collection_GetExtensions(t *testing.T) {
	c := NewCollection[testExtensible]()
	assert.Nil(t, c.GetExtensions())

	err := c.RegisterExtensions(Map{testPoint: &testExtensions{}})
	require.NoError(t, err)

	assert.NotNil(t, c.GetExtensions())
}

func Test_Map(t *testing.T) {
	m := NewMap()
	assert.Equal(t, 0, len(m))

	other := m.Add(testPoint, &testExtensible{})
	assert.Equal(t, 1, len(m))
	assert.Equal(t, m, other)
}
