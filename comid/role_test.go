// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoles_NewRoles_empty(t *testing.T) {
	tv := NewRoles()
	require.NotNil(t, tv)

	expected := "empty roles"

	assert.EqualError(t, tv.Valid(), expected)
}

func TestRoles_NewRoles_set_multi(t *testing.T) {
	actual := NewRoles().
		Add(RoleCreator).
		Add(RoleMaintainer, RoleTagCreator)
	require.NotNil(t, actual)

	expected := &Roles{
		RoleCreator,
		RoleMaintainer,
		RoleTagCreator,
	}

	assert.Nil(t, actual.Valid())
	assert.Equal(t, expected, actual)
}

func TestRoles_FromCBOR_ok(t *testing.T) {
	tvs := []struct {
		cbor     []byte
		expected Roles
	}{
		{
			// 83020100 = [2, 1, 0]
			cbor:     MustHexDecode(t, "83020100"),
			expected: Roles{RoleMaintainer, RoleCreator, RoleTagCreator},
		},
		{
			// use non-registered value -123
			// 8202387a = [2, -123]
			cbor:     MustHexDecode(t, "8202387a"),
			expected: Roles{RoleMaintainer, Role(-123)},
		},
	}

	for _, tv := range tvs {
		var actual Roles

		err := actual.FromCBOR(tv.cbor)

		assert.Nil(t, err)
		assert.Equal(t, tv.expected, actual)
	}
}

func TestRoles_FromCBOR_fail(t *testing.T) {
	tvs := []struct {
		cbor        []byte
		expectedErr string
	}{
		{
			// 80 = []
			cbor:        MustHexDecode(t, "80"),
			expectedErr: "empty roles",
		},
		{
			// 8163626c61 = ["bla"]
			cbor:        MustHexDecode(t, "8163626c61"),
			expectedErr: "cbor: cannot unmarshal UTF-8 text string into Go value of type comid.Role",
		},
		/*
			{
				// XXX(tho) - the cbor library treats null as 0
				// 81f6 = [null]
				cbor:        mustHexDecode(t, "81f6"),
				expectedErr: "..."
			},
		*/
	}

	for _, tv := range tvs {
		var actual Roles

		err := actual.FromCBOR(tv.cbor)

		assert.EqualError(t, err, tv.expectedErr)
	}
}

func TestRoles_ToCBOR_ok(t *testing.T) {
	tvs := []struct {
		roles    *Roles
		expected []byte
	}{
		{
			// [0, 1, 2]
			roles:    NewRoles().Add(RoleTagCreator, RoleCreator, RoleMaintainer),
			expected: MustHexDecode(t, "83000102"),
		},
		{
			// [0]
			roles:    NewRoles().Add(RoleTagCreator),
			expected: MustHexDecode(t, "8100"),
		},
		{
			// [1]
			roles:    NewRoles().Add(RoleCreator),
			expected: MustHexDecode(t, "8101"),
		},
		{
			// [2]
			roles:    NewRoles().Add(RoleMaintainer),
			expected: MustHexDecode(t, "8102"),
		},
		{
			// it is possible to force non-registered values
			// [1, 123]
			roles:    NewRoles().Add(RoleCreator, Role(123)),
			expected: MustHexDecode(t, "8201187b"),
		},
	}

	for _, tv := range tvs {
		actual, err := tv.roles.ToCBOR()

		fmt.Printf("CBOR: %x\n", actual)

		assert.Nil(t, err)
		assert.Equal(t, tv.expected, actual)
	}
}

func TestRoles_ToCBOR_fail_unset(t *testing.T) {
	r := NewRoles()
	require.NotNil(t, r)

	_, err := r.ToCBOR()

	assert.EqualError(t, err, "empty roles")
}

func TestRoles_UnmarshalJSON_ok(t *testing.T) {
	tvs := []struct {
		json     string
		expected Roles
	}{
		{
			json:     `[ "tagCreator", "creator", "maintainer" ]`,
			expected: Roles{RoleTagCreator, RoleCreator, RoleMaintainer},
		},
		{
			json:     `[ "creator" ]`,
			expected: Roles{RoleCreator},
		},
		{
			json:     `[ "maintainer" ]`,
			expected: Roles{RoleMaintainer},
		},
		{
			json:     `[ "tagCreator" ]`,
			expected: Roles{RoleTagCreator},
		},
	}

	for _, tv := range tvs {
		var actual Roles

		err := actual.UnmarshalJSON([]byte(tv.json))

		assert.Nil(t, err)
		assert.Equal(t, tv.expected, actual)
	}
}

func TestRoles_UnmarshalJSON_fail(t *testing.T) {
	tvs := []struct {
		json        string
		expectedErr string
	}{
		{
			json:        `[]`,
			expectedErr: "no roles found",
		},
		{
			json:        `["blabla"]`,
			expectedErr: `unknown role "blabla"`,
		},
		{
			json:        `"tagCreator"`,
			expectedErr: "json: cannot unmarshal string into Go value of type []string",
		},
	}

	for _, tv := range tvs {
		var actual Roles

		err := actual.UnmarshalJSON([]byte(tv.json))

		assert.EqualError(t, err, tv.expectedErr)
	}
}
