// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/eat"
	"github.com/veraison/swid"
)

func TestMeasurement_NewUUIDMeasurement_good_uuid(t *testing.T) {
	tv := NewUUIDMeasurement(TestUUID)

	assert.NotNil(t, tv)
}

func TestMeasurement_NewUUIDMeasurement_empty_uuid(t *testing.T) {
	emptyUUID := UUID{}

	tv := NewUUIDMeasurement(emptyUUID)

	assert.Nil(t, tv)
}

func TestMeasurement_NewPSAMeasurement_empty(t *testing.T) {
	emptyPSARefValID := PSARefValID{}

	tv := NewPSAMeasurement(emptyPSARefValID)

	assert.Nil(t, tv)
}

func TestMeasurement_NewPSAMeasurement_no_values(t *testing.T) {
	psaRefValID :=
		NewPSARefValID(TestSignerID).
			SetLabel("PRoT").
			SetVersion("1.2.3")
	require.NotNil(t, psaRefValID)

	tv := NewPSAMeasurement(*psaRefValID)
	assert.NotNil(t, tv)

	err := tv.Valid()
	assert.EqualError(t, err, "no measurement value set")
}

func TestMeasurement_NewPSAMeasurement_one_value(t *testing.T) {
	psaRefValID :=
		NewPSARefValID(TestSignerID).
			SetLabel("PRoT").
			SetVersion("1.2.3")
	require.NotNil(t, psaRefValID)

	tv := NewPSAMeasurement(*psaRefValID).SetIPaddr(TestIPaddr)
	assert.NotNil(t, tv)

	err := tv.Valid()
	assert.Nil(t, err)
}

func TestMeasurement_NewUUIDMeasurement_no_values(t *testing.T) {
	tv := NewUUIDMeasurement(TestUUID)
	require.NotNil(t, tv)

	err := tv.Valid()
	assert.EqualError(t, err, "no measurement value set")
}

func TestMeasurement_NewUUIDMeasurement_one_value(t *testing.T) {
	tv := NewUUIDMeasurement(TestUUID).SetMinSVN(2)
	require.NotNil(t, tv)

	err := tv.Valid()
	assert.Nil(t, err)
}

func TestMeasurement_NewUUIDMeasurement_bad_digest(t *testing.T) {
	tv := NewUUIDMeasurement(TestUUID)
	require.NotNil(t, tv)

	assert.Nil(t, tv.AddDigest(swid.Sha256, []byte{0xff}))
}

func TestMeasurement_NewUUIDMeasurement_bad_ueid(t *testing.T) {
	tv := NewUUIDMeasurement(TestUUID)
	require.NotNil(t, tv)

	badUEID := eat.UEID{
		0xFF, // Invalid
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	}

	assert.Nil(t, tv.SetUEID(badUEID))
}

func TestMeasurement_NewUUIDMeasurement_bad_uuid(t *testing.T) {
	tv := NewUUIDMeasurement(TestUUID)
	require.NotNil(t, tv)

	nonRFC4122UUID, err := ParseUUID("f47ac10b-58cc-4372-c567-0e02b2c3d479")
	require.Nil(t, err)

	assert.Nil(t, tv.SetUUID(nonRFC4122UUID))
}
