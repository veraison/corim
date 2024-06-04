// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntity_Valid_empty(t *testing.T) {
	tv := Entity{}

	err := tv.Valid()
	assert.EqualError(t, err, "invalid entity: empty entity-name")
}

func TestEntity_Valid_name_but_no_roles(t *testing.T) {
	tv := Entity{}

	require.NotNil(t, tv.SetEntityName("ACME Ltd."))

	err := tv.Valid()
	assert.EqualError(t, err, "invalid entity: empty roles")
}

func TestEntity_Valid_name_regid_but_no_roles(t *testing.T) {
	tv := Entity{}

	require.NotNil(t, tv.SetEntityName("ACME Ltd."))
	require.NotNil(t, tv.SetRegID("https://acme.example"))

	err := tv.Valid()
	assert.EqualError(t, err, "invalid entity: empty roles")
}

func TestEntity_Valid_name_regid_and_roles(t *testing.T) {
	tv := Entity{}

	require.NotNil(t, tv.SetEntityName("ACME Ltd."))
	require.NotNil(t, tv.SetRegID("https://acme.example"))
	require.NotNil(t, tv.SetRoles(RoleTagCreator))

	err := tv.Valid()
	assert.Nil(t, err)
}

func TestEntities_Valid_empty(t *testing.T) {
	e := Entity{}
	tv := NewEntities().AddEntity(e)
	require.NotNil(t, tv)

	err := tv.Valid()
	assert.EqualError(t, err, "entity at index 0: invalid entity: empty entity-name")
}

func TestEntities_Valid_ok(t *testing.T) {
	e := Entity{}

	require.NotNil(t,
		e.SetEntityName("ACME Ltd.").
			SetRegID("https://acme.example").
			SetRoles(RoleTagCreator, RoleCreator),
	)

	tv := NewEntities().AddEntity(e)
	require.NotNil(t, tv)

	err := tv.Valid()
	assert.Nil(t, err)
}

func TestEntity_SetEntityName_empty(t *testing.T) {
	e := Entity{}

	assert.Nil(t, e.SetEntityName(""))
}

func TestEntity_SetRegID_empty(t *testing.T) {
	e := Entity{}

	assert.Nil(t, e.SetRegID(""))
}

type testEntityName uint64

func newTestEntityName(val any) (*EntityName, error) {
	if val == nil {
		v := testEntityName(0)
		return &EntityName{&v}, nil
	}

	u, ok := val.(uint64)
	if !ok {
		return nil, errors.New("must be uint64")
	}

	v := testEntityName(u)
	return &EntityName{&v}, nil
}

func (o testEntityName) Type() string {
	return "test"
}

func (o testEntityName) String() string {
	return fmt.Sprint(uint64(o))
}

func (o testEntityName) Valid() error {
	return nil
}

type testEntityNameBadType struct {
	testEntityName
}

func newTestEntityNameBadType(_ any) (*EntityName, error) {
	v := testEntityNameBadType{testEntityName(7)}
	return &EntityName{&v}, nil
}

func (o testEntityNameBadType) Type() string {
	return "string"
}

func Test_RegisterEntityNameType(t *testing.T) {
	err := RegisterEntityNameType(32, newTestEntityName)
	assert.EqualError(t, err, "tag 32 is already registered")

	err = RegisterEntityNameType(99994, newTestEntityNameBadType)
	assert.EqualError(t, err, `entity name type with name "string" already exists`)

	registerTestEntityNameType(t)
}

// Since there only one, untagged, entity name type in the core spec, we use
// the test type define above in order to test the marshaling code works
// properly. Since global environment is not reset when running multiple tests,
// we cannot simply call RegisterEntityNameType() inside each test that relies
// on the test type, as that will cause the "tag already registered" error. On
// the other hand, we do not want to create inter-test dependencies by relying
// that the test registering the type is run before the others that rely on it.
// To get around this, use this global flag to only register the test type if a
// previous test hasn't already done so.
var testEntityNameTypeRegistered = false

func registerTestEntityNameType(t *testing.T) {
	if !testEntityNameTypeRegistered {
		err := RegisterEntityNameType(99994, newTestEntityName)
		require.NoError(t, err)

		testEntityNameTypeRegistered = true
	}
}

func TestEntityName_CBOR(t *testing.T) {
	registerTestEntityNameType(t)

	for _, tv := range []struct {
		Value          any
		Type           string
		ExpectedBytes  []byte
		ExpectedString string
	}{
		{
			Value: "test",
			Type:  "string",
			ExpectedBytes: []byte{
				0x64,                   // tstr(4)
				0x74, 0x65, 0x73, 0x74, // "test"
			},
			ExpectedString: "test",
		},
		{
			Value: uint64(7),
			Type:  "test",
			ExpectedBytes: []byte{
				0xda, 0x0, 0x1, 0x86, 0x9a, // tag 99994
				0x07, // unsigned int(7)
			},
			ExpectedString: "7",
		},
	} {
		t.Run(tv.Type, func(t *testing.T) {
			en, err := NewEntityName(tv.Value, tv.Type)
			require.NoError(t, err)

			data, err := en.MarshalCBOR()
			require.NoError(t, err)

			assert.Equal(t, tv.ExpectedBytes, data)

			var out EntityName

			err = out.UnmarshalCBOR(data)
			require.NoError(t, err)

			assert.Equal(t, tv.ExpectedString, out.String())
		})
	}
}

func TestEntityName_JSON(t *testing.T) {
	registerTestEntityNameType(t)

	for _, tv := range []struct {
		Value          any
		Type           string
		ExpectedBytes  []byte
		ExpectedString string
	}{
		{
			Value:          "test",
			Type:           "string",
			ExpectedBytes:  []byte(`"test"`),
			ExpectedString: "test",
		},
		{
			Value:          uint64(7),
			Type:           "test",
			ExpectedBytes:  []byte(`{"type":"test","value":7}`),
			ExpectedString: "7",
		},
	} {
		t.Run(tv.Type, func(t *testing.T) {
			en, err := NewEntityName(tv.Value, tv.Type)
			require.NoError(t, err)

			data, err := en.MarshalJSON()
			require.NoError(t, err)

			assert.Equal(t, tv.ExpectedBytes, data)

			var out EntityName

			err = out.UnmarshalJSON(data)
			require.NoError(t, err)

			assert.Equal(t, tv.ExpectedString, out.String())
		})
	}
}

func Test_NewStringEntityName(t *testing.T) {
	out, err := NewStringEntityName(nil)
	require.NoError(t, err)
	assert.EqualError(t, out.Valid(), "empty entity-name")

	out, err = NewStringEntityName([]byte("test"))
	require.NoError(t, err)
	assert.Equal(t, "test", out.String())

	_, err = NewStringEntityName(7)
	assert.EqualError(t, err, "unexpected type for string entity name: int")
}

func Test_MustNewEntityName(t *testing.T) {
	out := MustNewEntityName("test", "string")
	assert.Equal(t, "test", out.String())

	assert.Panics(t, func() {
		MustNewEntityName(7, "int")
	})
}
