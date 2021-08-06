// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
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
