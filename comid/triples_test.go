// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/extensions"
)

func TestTriples_extensions(t *testing.T) {
	triples := Triples{}
	triplesExt := &struct{}{}

	extMap := extensions.NewMap().
		Add(ExtTriples, triplesExt).
		Add(ExtReferenceValue, &struct{}{}).
		Add(ExtReferenceValueFlags, &struct{}{}).
		Add(ExtEndorsedValue, &struct{}{}).
		Add(ExtEndorsedValueFlags, &struct{}{})

	err := triples.RegisterExtensions(extMap)
	assert.NoError(t, err)
	assert.Equal(t, triplesExt, triples.GetExtensions())

	badMap := extensions.NewMap().Add(extensions.Point("test"), &struct{}{})
	err = triples.RegisterExtensions(badMap)
	assert.EqualError(t, err, `unexpected extension point: "test"`)
}

func TestTriples_marshaling(t *testing.T) {
	triples := Triples{}

	extMap := extensions.NewMap().
		Add(ExtReferenceValue, &struct{}{}).
		Add(ExtEndorsedValue, &struct{}{})

	require.NoError(t, triples.RegisterExtensions(extMap))

	data, err := triples.MarshalCBOR()
	assert.NoError(t, err)
	assert.Equal(t, data, []byte{0xa0})

	data, err = triples.MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, "{}", string(data))
}

func TestTriples_Valid(t *testing.T) {
	triples := Triples{}
	triples.ReferenceValues = &ValueTriples{}
	triples.EndorsedValues = &ValueTriples{}

	err := triples.Valid()
	assert.EqualError(t, err, "triples struct must not be empty")

	triples.ReferenceValues.Add(&ValueTriple{})
	err = triples.Valid()
	assert.EqualError(t, err, "reference values: error at index 0: environment validation failed: environment must not be empty")

	triples.ReferenceValues = nil
	triples.EndorsedValues.Add(&ValueTriple{})
	err = triples.Valid()
	assert.EqualError(t, err, "endorsed values: error at index 0: environment validation failed: environment must not be empty")

	triples.EndorsedValues = nil
	triples.AttestVerifKeys = &KeyTriples{{}}
	err = triples.Valid()
	assert.EqualError(t, err, "attestation verification key at index 0: environment validation failed: environment must not be empty")

	triples.AttestVerifKeys = nil
	triples.DevIdentityKeys = &KeyTriples{{}}
	err = triples.Valid()
	assert.EqualError(t, err, "device identity key at index 0: environment validation failed: environment must not be empty")
}

func TestTriples_adders(t *testing.T) {
	triples := Triples{}

	triples.AddReferenceValue(ValueTriple{}).AddEndorsedValue(ValueTriple{})
	assert.Len(t, triples.ReferenceValues.Values, 1)
	assert.Len(t, triples.EndorsedValues.Values, 1)
}
