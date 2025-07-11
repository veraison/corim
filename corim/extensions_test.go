// Copyright 2023-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package corim

import (
	"errors"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/extensions"
)

type TestExtensions struct {
	Address string `cbor:"-1,keyasint,omitempty" json:"address,omitempty"`
	Size    int    `cbor:"-2,keyasint,omitempty" json:"size,omitempty"`
}

func (o TestExtensions) ConstrainEntity(ent *Entity) error {
	if ent.Name.String() != "Futurama" {
		return errors.New(`EntityName must be "Futurama"`) // nolint:golint
	}

	return nil
}

func (o TestExtensions) ConstrainCorim(_ *UnsignedCorim) error {
	return errors.New("invalid")
}

func (o TestExtensions) ConstrainSigner(_ *Signer) error {
	return errors.New("invalid")
}

func TestEntityExtensions_Valid(t *testing.T) {
	ent := NewEntity()
	ent.SetName("The Simpsons")
	ent.SetRoles(RoleManifestCreator)

	err := ent.Valid()
	assert.NoError(t, err)

	extMap := extensions.NewMap().Add(ExtEntity, &TestExtensions{})
	err = ent.RegisterExtensions(extMap)
	require.NoError(t, err)

	err = ent.Valid()
	assert.EqualError(t, err, `EntityName must be "Futurama"`)

	ent.SetName("Futurama")
	err = ent.Valid()
	assert.NoError(t, err)

	assert.EqualError(t, ent.validCorim(nil), "invalid")
	assert.EqualError(t, ent.validSigner(nil), "invalid")
}

func TestEntityExtensions_CBOR(t *testing.T) {
	data := []byte{
		0xa4, // map(4)

		0x00,                   // key 0
		0x64,                   // val tstr(4)
		0x61, 0x63, 0x6d, 0x65, // "acme"

		0x02, // key 2
		0x81, // array(1)
		0x01, // 1

		0x20,             // key -1
		0x63,             // val tstr(3)
		0x66, 0x6f, 0x6f, // "foo"

		0x21, // key -2
		0x06, // val 6
	}

	ent := NewEntity()
	extMap := extensions.NewMap().Add(ExtEntity, &TestExtensions{})
	err := ent.RegisterExtensions(extMap)
	require.NoError(t, err)

	err = cbor.Unmarshal(data, &ent)
	assert.NoError(t, err)

	assert.Equal(t, ent.Name.String(), "acme")

	address, err := ent.Get("address")
	require.NoError(t, err)
	assert.Equal(t, address, "foo")
}
