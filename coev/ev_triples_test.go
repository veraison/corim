// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/extensions"
)

func TestEvTriples_NewEvTriples(t *testing.T) {
	evTriples := NewEvTriples()
	assert.NotNil(t, evTriples)
	assert.Nil(t, evTriples.EvidenceTriples)
	assert.Nil(t, evTriples.IdentityTriples)
	assert.Nil(t, evTriples.CoSWIDTriples)
	assert.Nil(t, evTriples.AttestKeysTriples)
}

func TestEvTriples_Valid_empty(t *testing.T) {
	evTriples := NewEvTriples()
	err := evTriples.Valid()
	assert.EqualError(t, err, "no Triples set inside EvTriples")
}

func TestEvTriples_Valid_with_evidence_triples(t *testing.T) {
	evTriples := NewEvTriples()

	// Add a valid evidence triple with proper measurement value
	measurement := comid.MustNewUUIDMeasurement(TestUUID)
	measurement.SetRawValueBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{0xff, 0xff, 0xff, 0xff})

	valueTriple := &comid.ValueTriple{
		Environment: comid.Environment{
			Class: comid.NewClassOID(comid.TestOID).
				SetVendor("Test Vendor").
				SetModel("Test Model"),
		},
		Measurements: *comid.NewMeasurements().Add(measurement),
	}

	evTriples.AddEvidenceTriple(valueTriple)

	err := evTriples.Valid()
	assert.NoError(t, err)
}

func TestEvTriples_Valid_with_identity_triples(t *testing.T) {
	evTriples := NewEvTriples()

	// Add a valid identity triple
	keyTriple := &comid.KeyTriple{
		Environment: comid.Environment{
			Class: comid.NewClassUUID(TestUUID).
				SetVendor("Test Vendor"),
		},
		VerifKeys: *comid.NewCryptoKeys().
			Add(comid.MustNewPKIXBase64Key(comid.TestECPubKey)),
	}

	evTriples.AddIdentityTriple(keyTriple)

	err := evTriples.Valid()
	assert.NoError(t, err)
}

func TestEvTriples_Valid_with_coswid_triples(t *testing.T) {
	evTriples := NewEvTriples()

	// Add a valid CoSWID triple with evidence and empty evidence map
	coswidEvidence := NewCoSWIDEvidence().AddCoSWIDEvidenceMap(&CoSWIDEvidenceMap{
		Evidence: TestEvidence,
	})

	coswidTriple := &CoSWIDTriple{
		Environment: comid.Environment{
			Instance: comid.MustNewUEIDInstance(comid.TestUEID),
		},
		Evidence: *coswidEvidence,
	}

	evTriples.AddCoSWIDTriple(coswidTriple)

	err := evTriples.Valid()
	assert.NoError(t, err)
}

func TestEvTriples_Valid_with_attest_key_triples(t *testing.T) {
	evTriples := NewEvTriples()

	// Add a valid attest key triple
	keyTriple := &comid.KeyTriple{
		Environment: comid.Environment{
			Class: comid.NewClassUUID(TestUUID).
				SetVendor("Test Vendor"),
		},
		VerifKeys: *comid.NewCryptoKeys().
			Add(comid.MustNewPKIXBase64Key(comid.TestECPubKey)),
	}

	evTriples.AddAttestKeyTriple(keyTriple)

	err := evTriples.Valid()
	assert.NoError(t, err)
}

func TestEvTriples_AddEvidenceTriple_nil_receiver(t *testing.T) {
	var evTriples *EvTriples
	valueTriple := &comid.ValueTriple{}

	result := evTriples.AddEvidenceTriple(valueTriple)
	assert.Nil(t, result)
}

func TestEvTriples_AddEvidenceTriple_success(t *testing.T) {
	evTriples := NewEvTriples()

	measurement := comid.MustNewUUIDMeasurement(TestUUID)
	measurement.SetRawValueBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{0xff, 0xff, 0xff, 0xff})

	valueTriple := &comid.ValueTriple{
		Environment: comid.Environment{
			Class: comid.NewClassOID(comid.TestOID).
				SetVendor("Test Vendor").
				SetModel("Test Model"),
		},
		Measurements: *comid.NewMeasurements().Add(measurement),
	}

	result := evTriples.AddEvidenceTriple(valueTriple)

	assert.Equal(t, evTriples, result)
	assert.NotNil(t, evTriples.EvidenceTriples)
	assert.Len(t, evTriples.EvidenceTriples.Values, 1)
	assert.Equal(t, *valueTriple, evTriples.EvidenceTriples.Values[0])
}

func TestEvTriples_AddEvidenceTriple_multiple(t *testing.T) {
	evTriples := NewEvTriples()

	measurement1 := comid.MustNewUUIDMeasurement(TestUUID)
	measurement1.SetRawValueBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{0xff, 0xff, 0xff, 0xff})

	valueTriple1 := &comid.ValueTriple{
		Environment: comid.Environment{
			Class: comid.NewClassOID(comid.TestOID).
				SetVendor("Test Vendor 1").
				SetModel("Test Model 1"),
		},
		Measurements: *comid.NewMeasurements().Add(measurement1),
	}

	measurement2 := comid.MustNewUUIDMeasurement(TestUUID)
	measurement2.SetRawValueBytes([]byte{0x05, 0x06, 0x07, 0x08}, []byte{0xff, 0xff, 0xff, 0xff})

	valueTriple2 := &comid.ValueTriple{
		Environment: comid.Environment{
			Class: comid.NewClassOID(comid.TestOID).
				SetVendor("Test Vendor 2").
				SetModel("Test Model 2"),
		},
		Measurements: *comid.NewMeasurements().Add(measurement2),
	}

	evTriples.AddEvidenceTriple(valueTriple1)
	evTriples.AddEvidenceTriple(valueTriple2)

	assert.NotNil(t, evTriples.EvidenceTriples)
	assert.Len(t, evTriples.EvidenceTriples.Values, 2)
	assert.Equal(t, *valueTriple1, evTriples.EvidenceTriples.Values[0])
	assert.Equal(t, *valueTriple2, evTriples.EvidenceTriples.Values[1])
}

func TestEvTriples_AddCoSWIDTriple_nil_receiver(t *testing.T) {
	var evTriples *EvTriples
	coswidTriple := &CoSWIDTriple{}

	result := evTriples.AddCoSWIDTriple(coswidTriple)
	assert.Nil(t, result)
}

func TestEvTriples_AddCoSWIDTriple_success(t *testing.T) {
	evTriples := NewEvTriples()
	coswidTriple := &CoSWIDTriple{
		Environment: comid.Environment{
			Instance: comid.MustNewUEIDInstance(comid.TestUEID),
		},
		Evidence: *NewCoSWIDEvidence(),
	}

	result := evTriples.AddCoSWIDTriple(coswidTriple)

	assert.Equal(t, evTriples, result)
	assert.NotNil(t, evTriples.CoSWIDTriples)
	assert.Len(t, *evTriples.CoSWIDTriples, 1)
	assert.Equal(t, *coswidTriple, (*evTriples.CoSWIDTriples)[0])
}

func TestEvTriples_AddIdentityTriple_nil_receiver(t *testing.T) {
	var evTriples *EvTriples
	keyTriple := &comid.KeyTriple{}

	result := evTriples.AddIdentityTriple(keyTriple)
	assert.Nil(t, result)
}

func TestEvTriples_AddIdentityTriple_success(t *testing.T) {
	evTriples := NewEvTriples()
	keyTriple := &comid.KeyTriple{
		Environment: comid.Environment{
			Class: comid.NewClassUUID(TestUUID).
				SetVendor("Test Vendor"),
		},
		VerifKeys: *comid.NewCryptoKeys().
			Add(comid.MustNewPKIXBase64Key(comid.TestECPubKey)),
	}

	result := evTriples.AddIdentityTriple(keyTriple)

	assert.Equal(t, evTriples, result)
	assert.NotNil(t, evTriples.IdentityTriples)
	assert.Len(t, *evTriples.IdentityTriples, 1)
	assert.Equal(t, *keyTriple, (*evTriples.IdentityTriples)[0])
}

func TestEvTriples_AddAttestKeyTriple_nil_receiver(t *testing.T) {
	var evTriples *EvTriples
	keyTriple := &comid.KeyTriple{}

	result := evTriples.AddAttestKeyTriple(keyTriple)
	assert.Nil(t, result)
}

func TestEvTriples_AddAttestKeyTriple_success(t *testing.T) {
	evTriples := NewEvTriples()
	keyTriple := &comid.KeyTriple{
		Environment: comid.Environment{
			Class: comid.NewClassUUID(TestUUID).
				SetVendor("Test Vendor"),
		},
		VerifKeys: *comid.NewCryptoKeys().
			Add(comid.MustNewPKIXBase64Key(comid.TestECPubKey)),
	}

	result := evTriples.AddAttestKeyTriple(keyTriple)

	assert.Equal(t, evTriples, result)
	assert.NotNil(t, evTriples.AttestKeysTriples)
	assert.Len(t, *evTriples.AttestKeysTriples, 1)
	assert.Equal(t, *keyTriple, (*evTriples.AttestKeysTriples)[0])
}

func TestEvTriples_RegisterExtensions(t *testing.T) {
	evTriples := NewEvTriples()

	// Create test extensions
	testExts := extensions.NewMap().Add(ExtEvTriples, &struct{}{})

	err := evTriples.RegisterExtensions(testExts)
	assert.NoError(t, err)
}

func TestEvTriples_CBOR_roundtrip(t *testing.T) {
	original := NewEvTriples()

	// Add some test data with proper measurement value
	measurement := comid.MustNewUUIDMeasurement(TestUUID)
	measurement.SetRawValueBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{0xff, 0xff, 0xff, 0xff})

	valueTriple := &comid.ValueTriple{
		Environment: comid.Environment{
			Class: comid.NewClassOID(comid.TestOID).
				SetVendor("Test Vendor").
				SetModel("Test Model"),
		},
		Measurements: *comid.NewMeasurements().Add(measurement),
	}
	original.AddEvidenceTriple(valueTriple)

	// Marshal to CBOR
	cbor, err := original.MarshalCBOR()
	require.NoError(t, err)
	assert.NotEmpty(t, cbor)

	// Unmarshal from CBOR
	decoded := &EvTriples{}
	err = decoded.UnmarshalCBOR(cbor)
	require.NoError(t, err)

	// Verify the roundtrip
	assert.NotNil(t, decoded.EvidenceTriples)
	assert.Len(t, decoded.EvidenceTriples.Values, 1)

	err = decoded.Valid()
	assert.NoError(t, err)
}

func TestEvTriples_JSON_roundtrip(t *testing.T) {
	original := NewEvTriples()

	// Add some test data
	keyTriple := &comid.KeyTriple{
		Environment: comid.Environment{
			Class: comid.NewClassUUID(TestUUID).
				SetVendor("Test Vendor"),
		},
		VerifKeys: *comid.NewCryptoKeys().
			Add(comid.MustNewPKIXBase64Key(comid.TestECPubKey)),
	}
	original.AddIdentityTriple(keyTriple)

	// Marshal to JSON
	jsonData, err := original.MarshalJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Unmarshal from JSON
	decoded := &EvTriples{}
	err = decoded.UnmarshalJSON(jsonData)
	require.NoError(t, err)

	// Verify the roundtrip
	assert.NotNil(t, decoded.IdentityTriples)
	assert.Len(t, *decoded.IdentityTriples, 1)

	err = decoded.Valid()
	assert.NoError(t, err)
}

func TestEvTriples_Valid_invalid_evidence_triple(t *testing.T) {
	evTriples := NewEvTriples()

	// Add an invalid evidence triple (empty environment and measurements)
	valueTriple := &comid.ValueTriple{}
	evTriples.AddEvidenceTriple(valueTriple)

	err := evTriples.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid EvidenceTriples")
}

func TestEvTriples_Valid_invalid_identity_triple(t *testing.T) {
	evTriples := NewEvTriples()

	// Add an invalid identity triple (empty environment and keys)
	keyTriple := &comid.KeyTriple{}
	evTriples.AddIdentityTriple(keyTriple)

	err := evTriples.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid IdentityTriple at index: 0")
}

func TestEvTriples_Valid_invalid_coswid_triple(t *testing.T) {
	evTriples := NewEvTriples()

	// Add an invalid CoSWID triple (empty environment)
	coswidTriple := &CoSWIDTriple{}
	evTriples.AddCoSWIDTriple(coswidTriple)

	err := evTriples.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid CoSWIDTriple at index: 0")
}

func TestEvTriples_Valid_invalid_attest_key_triple(t *testing.T) {
	evTriples := NewEvTriples()

	// Add an invalid attest key triple (empty environment and keys)
	keyTriple := &comid.KeyTriple{}
	evTriples.AddAttestKeyTriple(keyTriple)

	err := evTriples.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid AttestKeysTriple at index: 0")
}

func TestEvTriples_GetExtensions(t *testing.T) {
	evTriples := NewEvTriples()

	// Initially should return nil
	assert.Nil(t, evTriples.GetExtensions())

	// Register extensions
	testExts := extensions.NewMap().Add(ExtEvTriples, &struct{}{})
	err := evTriples.RegisterExtensions(testExts)
	require.NoError(t, err)

	// Should now return the extensions
	assert.NotNil(t, evTriples.GetExtensions())
}
