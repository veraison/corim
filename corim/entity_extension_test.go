// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"testing"

	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntities_Extensions_Valid_ok(t *testing.T) {
	ext := &EntityExtension{"my profile", 0x56}

	e := NewEntity().
		SetEntityName("ACME Ltd.").
		SetRegID("http://acme.example").
		SetRoles(RoleManifestCreator).
		SetExtension(ext)
	require.NotNil(t, e)

	es := NewEntities().AddEntity(*e)
	require.NotNil(t, es)

	err := es.Valid()
	assert.Nil(t, err)
}

func TestEntities_Extensions_ToCBOR_ok(t *testing.T) {
	ext := &EntityExtension{"my profile", 0x56}

	e := NewEntity().
		SetEntityName("ACME Ltd.").
		SetRegID("http://acme.example").
		SetRoles(RoleManifestCreator).
		SetExtension(ext)
	require.NotNil(t, e)
	data, err := e.ToCBOR()

	require.Nil(t, err)
	fmt.Printf("to CBOR Entity = %x", data)

}

func TestEntities_Extensions_FromCBOR_ok(t *testing.T) {
	data := []byte{0xa4, 0x00, 0x69, 0x41, 0x43, 0x4d, 0x45, 0x20, 0x4c, 0x74, 0x64, 0x2e, 0x01, 0xd8, 0x20, 0x73, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x61, 0x63, 0x6d, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x02, 0x81, 0x01, 0x03, 0xa2, 0x00, 0x6a, 0x6d, 0x79, 0x20, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x01, 0x18, 0x56}
	//var ext1 EntityExtension
	e := &Entity{Extension: &EntityExtension{}}

	err := e.FromCBOR(data)
	require.Nil(t, err)
	ext, err := e.GetExtension()
	require.Nil(t, err)
	require.NotNil(t, ext)
}

func TestEntities_NoExtensions_ToCBOR_ok(t *testing.T) {
	e := NewEntity().
		SetEntityName("ACME Ltd.").
		SetRegID("http://acme.example").
		SetRoles(RoleManifestCreator)
	require.NotNil(t, e)
	data, err := e.ToCBOR()

	require.Nil(t, err)
	fmt.Printf("to CBOR Entity = %x", data)

}

func TestEntities_NoExtensions_FromCBOR_ok(t *testing.T) {
	data := []byte{0xa4, 0x00, 0x69, 0x41, 0x43, 0x4d, 0x45, 0x20, 0x4c, 0x74, 0x64, 0x2e, 0x01, 0xd8, 0x20, 0x73, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x61, 0x63, 0x6d, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x02, 0x81, 0x01, 0x03, 0x0F6}
	e := &Entity{}
	err := e.FromCBOR(data)
	require.Nil(t, err)
	require.Nil(t, err)

}
