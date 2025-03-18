// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironment_Valid_empty(t *testing.T) {
	tv := Environment{}

	err := tv.Valid()

	assert.EqualError(t, err, "environment must not be empty")
}

func TestEnvironment_Valid_empty_class(t *testing.T) {
	tv := Environment{
		Class: &Class{},
	}

	err := tv.Valid()

	assert.EqualError(t, err, "class validation failed: class must not be empty")
}

func TestEnvironment_Valid_empty_instance(t *testing.T) {
	tv := Environment{
		Instance: &Instance{},
	}

	err := tv.Valid()

	assert.EqualError(t, err, "instance validation failed: invalid instance id")
}

func TestEnvironment_Valid_empty_group(t *testing.T) {
	tv := Environment{
		Group: &Group{},
	}

	err := tv.Valid()

	assert.EqualError(t, err, "group validation failed: no value set")
}
func TestEnvironment_Valid_ok_with_class(t *testing.T) {
	tv := Environment{
		Class: NewClassUUID(TestUUID),
	}

	err := tv.Valid()

	assert.Nil(t, err)
}

func TestEnvironment_ToCBOR_empty(t *testing.T) {
	var actual Environment
	_, err := actual.ToCBOR()

	assert.EqualError(t, err, "environment must not be empty")
}

func TestEnvironment_ToCBOR_class_only(t *testing.T) {
	tv := Environment{
		Class: NewClassUUID(TestUUID),
	}
	require.NotNil(t, tv.Class)

	// {0: {0: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}}
	expected := MustHexDecode(t, "a100a100d8255031fb5abf023e4992aa4e95f9c1503bfa")

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestEnvironment_ToCBOR_instance_only(t *testing.T) {
	tv := Environment{
		Instance: MustNewUEIDInstance(TestUEID),
	}
	require.NotNil(t, tv.Instance)

	// {1: 550(h'02DEADBEEFDEAD')}
	expected := MustHexDecode(t, "a101d902264702deadbeefdead")

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestEnvironment_ToCBOR_group_only(t *testing.T) {
	tv := Environment{
		Group: MustNewUUIDGroup(TestUUID),
	}
	require.NotNil(t, tv.Group)

	// {2: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}
	expected := MustHexDecode(t, "a102d8255031fb5abf023e4992aa4e95f9c1503bfa")

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestEnvironment_ToCBOR_class_and_instance(t *testing.T) {
	tv := Environment{
		Class:    NewClassUUID(TestUUID),
		Instance: MustNewUEIDInstance(TestUEID),
	}
	require.NotNil(t, tv.Class)
	require.NotNil(t, tv.Instance)

	// {0: {0: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}, 1: 550(h'02DEADBEEFDEAD')}
	expected := MustHexDecode(t, "a200a100d8255031fb5abf023e4992aa4e95f9c1503bfa01d902264702deadbeefdead")

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestEnvironment_ToCBOR_class_and_group(t *testing.T) {
	tv := Environment{
		Class: NewClassUUID(TestUUID),
		Group: MustNewUUIDGroup(TestUUID),
	}
	require.NotNil(t, tv.Class)
	require.NotNil(t, tv.Group)

	// {0: {0: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}, 2: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}
	expected := MustHexDecode(t, "a200a100d8255031fb5abf023e4992aa4e95f9c1503bfa02d8255031fb5abf023e4992aa4e95f9c1503bfa")

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestEnvironment_ToCBOR_instance_and_group(t *testing.T) {
	tv := Environment{
		Instance: MustNewUEIDInstance(TestUEID),
		Group:    MustNewUUIDGroup(TestUUID),
	}
	require.NotNil(t, tv.Instance)
	require.NotNil(t, tv.Group)

	// {1: 550(h'02DEADBEEFDEAD'), 2: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}
	expected := MustHexDecode(t, "a201d902264702deadbeefdead02d8255031fb5abf023e4992aa4e95f9c1503bfa")

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestEnvironment_ToCBOR_class_and_instance_and_group(t *testing.T) {
	tv := Environment{
		Class:    NewClassUUID(TestUUID),
		Instance: MustNewUEIDInstance(TestUEID),
		Group:    MustNewUUIDGroup(TestUUID),
	}
	require.NotNil(t, tv.Class)
	require.NotNil(t, tv.Instance)
	require.NotNil(t, tv.Group)

	// {0: {0: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}, 1: 550(h'02DEADBEEFDEAD'), 2: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}
	expected := MustHexDecode(t, "a300a100d8255031fb5abf023e4992aa4e95f9c1503bfa01d902264702deadbeefdead02d8255031fb5abf023e4992aa4e95f9c1503bfa")

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestEnvironment_FromCBOR_unknown_map_entry(t *testing.T) {
	// { 3: "unknown" }
	tv := MustHexDecode(t, "a10367756e6b6e6f776e")

	var actual Environment
	err := actual.FromCBOR(tv)

	// since the unknown map entry is ignored, the resulting Environment is empty
	assert.EqualError(t, err, "environment must not be empty")
}

func TestEnvironment_FromCBOR_empty(t *testing.T) {
	tv := MustHexDecode(t, "a0")

	var actual Environment
	err := actual.FromCBOR(tv)

	assert.EqualError(t, err, "environment must not be empty")
}

func TestEnvironment_FromCBOR_class_only(t *testing.T) {
	// {0: {0: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}}
	tv := MustHexDecode(t, "a100a100d8255031fb5abf023e4992aa4e95f9c1503bfa")

	var actual Environment
	err := actual.FromCBOR(tv)

	assert.Nil(t, err)
	assert.NotNil(t, actual.Class)
	assert.Equal(t, TestUUIDString, actual.Class.ClassID.String())
	assert.Nil(t, actual.Instance)
	assert.Nil(t, actual.Group)
}

func TestEnvironment_FromCBOR_instance_only(t *testing.T) {
	// {1: 550(h'02DEADBEEFDEAD')}
	tv := MustHexDecode(t, "a101d902264702deadbeefdead")

	var actual Environment
	err := actual.FromCBOR(tv)

	assert.Nil(t, err)
	assert.Nil(t, actual.Class)
	assert.NotNil(t, actual.Instance)
	assert.Equal(t, []byte(TestUEID), actual.Instance.Bytes())
	assert.Nil(t, actual.Group)
}

func TestEnvironment_FromCBOR_group_only(t *testing.T) {
	// {2: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}
	tv := MustHexDecode(t, "a102d8255031fb5abf023e4992aa4e95f9c1503bfa")

	var actual Environment
	err := actual.FromCBOR(tv)

	assert.Nil(t, err)
	assert.Nil(t, actual.Class)
	assert.Nil(t, actual.Instance)
	assert.NotNil(t, actual.Group)
	assert.Equal(t, TestUUIDString, actual.Group.String())
}

func TestEnvironment_FromCBOR_class_and_instance(t *testing.T) {
	// {0: {0: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}, 1: 550(h'02DEADBEEFDEAD')}
	tv := MustHexDecode(t, "a200a100d8255031fb5abf023e4992aa4e95f9c1503bfa01d902264702deadbeefdead")

	var actual Environment
	err := actual.FromCBOR(tv)

	assert.Nil(t, err)
	assert.NotNil(t, actual.Class)
	assert.Equal(t, TestUUIDString, actual.Class.ClassID.String())
	assert.NotNil(t, actual.Instance)
	assert.Equal(t, []byte(TestUEID), actual.Instance.Bytes())
	assert.Nil(t, actual.Group)
}

func TestEnvironment_FromCBOR_class_and_group(t *testing.T) {
	// {0: {0: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}, 2: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}
	tv := MustHexDecode(t, "a200a100d8255031fb5abf023e4992aa4e95f9c1503bfa02d8255031fb5abf023e4992aa4e95f9c1503bfa")

	var actual Environment
	err := actual.FromCBOR(tv)

	assert.Nil(t, err)
	assert.NotNil(t, actual.Class)
	assert.Equal(t, TestUUIDString, actual.Class.ClassID.String())
	assert.Nil(t, actual.Instance)
	assert.NotNil(t, actual.Group)
	assert.Equal(t, TestUUIDString, actual.Group.String())
}

func TestEnvironment_FromCBOR_instance_and_group(t *testing.T) {
	// {1: 550(h'02DEADBEEFDEAD'), 2: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}
	tv := MustHexDecode(t, "a201d902264702deadbeefdead02d8255031fb5abf023e4992aa4e95f9c1503bfa")

	var actual Environment
	err := actual.FromCBOR(tv)

	assert.Nil(t, err)
	assert.Nil(t, actual.Class)
	assert.NotNil(t, actual.Instance)
	assert.Equal(t, []byte(TestUEID), actual.Instance.Bytes())
	assert.NotNil(t, actual.Group)
	assert.Equal(t, TestUUIDString, actual.Group.String())
}

func TestEnvironment_FromCBOR_class_and_instance_and_group(t *testing.T) {
	// {0: {0: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}, 1: 550(h'02DEADBEEFDEAD'), 2: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}
	tv := MustHexDecode(t, "a300a100d8255031fb5abf023e4992aa4e95f9c1503bfa01d902264702deadbeefdead02d8255031fb5abf023e4992aa4e95f9c1503bfa")

	var actual Environment
	err := actual.FromCBOR(tv)

	assert.Nil(t, err)
	assert.NotNil(t, actual.Class)
	assert.Equal(t, TestUUIDString, actual.Class.ClassID.String())
	assert.NotNil(t, actual.Instance)
	assert.Equal(t, []byte(TestUEID), actual.Instance.Bytes())
	assert.NotNil(t, actual.Group)
	assert.Equal(t, TestUUIDString, actual.Group.String())
}

func TestEnviroment_JSON(t *testing.T) {
	testEnv := Environment{
		Class: NewClassUUID(TestUUID),
	}

	out, err := testEnv.ToJSON()
	require.NoError(t, err)
	assert.Equal(t, `{"class":{"id":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}}}`, string(out))

	var outEnv Environment

	err = outEnv.FromJSON(out)
	require.NoError(t, err)
	assert.Equal(t, testEnv, outEnv)

	_, err = Environment{}.ToJSON()
	assert.EqualError(t, err, "environment must not be empty")

	err = outEnv.FromJSON([]byte(`{"class": 7}`))
	assert.EqualError(t, err, "json: cannot unmarshal number into Go struct field Environment.class of type comid.Class")
}
