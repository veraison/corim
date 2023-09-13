// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoles_ToJSON_ok(t *testing.T) {
	tvs := []struct {
		roles    *Roles
		expected string
	}{
		{
			roles:    NewRoles().Add(RoleManifestCreator),
			expected: `[ "manifestCreator" ]`,
		},
	}

	for _, tv := range tvs {
		actual, err := tv.roles.ToJSON()

		assert.NoError(t, err)
		assert.JSONEq(t, tv.expected, string(actual))
	}
}

func TestRoles_ToJSON_fail(t *testing.T) {
	tvs := []struct {
		roles       *Roles
		expectedErr string
	}{
		{
			roles:       NewRoles(),
			expectedErr: "validation failed: empty roles",
		},
	}

	for _, tv := range tvs {
		_, err := tv.roles.ToJSON()

		assert.EqualError(t, err, tv.expectedErr)
	}
}

func TestRoles_FromJSON_fail(t *testing.T) {
	tvs := []struct {
		json        string
		expectedErr string
	}{
		{
			json:        `[]`,
			expectedErr: "validation failed: empty roles",
		},
		{
			json:        `["blabla"]`,
			expectedErr: `decoding failed: unknown role "blabla"`,
		},
		{
			json:        `[ "manifestCreator", "xyz" ]`,
			expectedErr: `decoding failed: unknown role "xyz"`,
		},
		{
			json:        `"tagCreator"`,
			expectedErr: "decoding failed: json: cannot unmarshal string into Go value of type []string",
		},
	}

	for _, tv := range tvs {
		var actual Roles

		err := actual.FromJSON([]byte(tv.json))

		assert.EqualError(t, err, tv.expectedErr)
	}
}

func Test_Role_String(t *testing.T) {
	assert.Equal(t, "manifestCreator", RoleManifestCreator.String())
	assert.Equal(t, "Role(9999)", Role(9999).String())
}

func Test_RegisterRole(t *testing.T) {
	err := RegisterRole(1, "owner")
	assert.EqualError(t, err, "role with value 1 already exists")

	err = RegisterRole(3, "manifestCreator")
	assert.EqualError(t, err, `role with name "manifestCreator" already exists`)

	err = RegisterRole(3, "owner")
	assert.NoError(t, err)

	roles := NewRoles().Add(Role(3))

	out, err := roles.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `["owner"]`, string(out))
}
