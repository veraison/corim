// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"testing"
)

func TestEntity_Valid_uninitialized(t *testing.T) {
	tv := Entity{}

	err := tv.Valid()

	assert.EqualError(t, err, "invalid entity: empty entity-name")
}

func TestEntity_Valid_empty_name(t *testing.T) {
	tv := Entity{
		EntityName: "",
	}

	err := tv.Valid()

	assert.EqualError(t, err, "invalid entity: empty entity-name")
}

func TestEntity_Valid_non_nil_empty_URI(t *testing.T) {
	emptyRegID := comid.TaggedURI("")

	tv := Entity{
		EntityName: "ACME Ltd.",
		RegID:      &emptyRegID,
	}

	err := tv.Valid()

	assert.EqualError(t, err, "invalid entity: empty reg-id")
}

func TestEntity_Valid_missing_roles(t *testing.T) {
	regID := comid.TaggedURI("http://acme.example")

	tv := Entity{
		EntityName: "ACME Ltd.",
		RegID:      &regID,
	}

	err := tv.Valid()

	assert.EqualError(t, err, "invalid entity: empty roles")
}

func TestEntity_Valid_unknown_role(t *testing.T) {
	regID := comid.TaggedURI("http://acme.example")

	tv := Entity{
		EntityName: "ACME Ltd.",
		RegID:      &regID,
		Roles:      Roles{Role(666)},
	}

	err := tv.Valid()

	assert.EqualError(t, err, "invalid entity: unknown role 666 at index 0")
}

func TestEntities_Valid_ok(t *testing.T) {
	e := NewEntity().
		SetEntityName("ACME Ltd.").
		SetRegID("http://acme.example").
		SetRoles(RoleManifestCreator)
	require.NotNil(t, e)

	es := NewEntities().AddEntity(e)
	require.NotNil(t, es)

	err := es.Valid()
	assert.Nil(t, err)
}

func TestEntities_Valid_empty(t *testing.T) {
	e := Entity{}

	es := NewEntities().AddEntity(&e)
	require.NotNil(t, es)

	err := es.Valid()
	assert.EqualError(t, err, "entity at index 0: invalid entity: empty entity-name")
}

func TestEntities_Valid1_ok(t *testing.T) {
	e := NewEntity().
		SetEntityName("ACME Ltd.").
		SetRegID("http://acme.example").
		SetRoles(RoleManifestCreator)
	require.NotNil(t, e)

	err := e.Valid()
	assert.Nil(t, err)
	data, err := e.ToCBOR()
	assert.Nil(t, err)
	fmt.Printf("Encoded CBOR Payload= %x", data)
	assert.NotNil(t, data)
}

func TestEntities_Valid2_ok(t *testing.T) {
	//ext := &EntityExtension{"myname", 0x20}
	e := NewEntity().
		SetEntityName("ACME Ltd.").
		SetRegID("http://acme.example").
		SetRoles(RoleManifestCreator)
	require.NotNil(t, e)
	//e.SetEntityExtension(ext)

	err := e.Valid()
	assert.Nil(t, err)
	data, err := e.ToCBOR()
	assert.Nil(t, err)
	fmt.Printf("Encoded CBOR Payload= %x", data)
	assert.NotNil(t, data)
}
