package psa

// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/corim"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
)

const ProfileURI = "tag:arm.com,2025:psa#1.0.0"

// PSARefValIDTag is the CBOR tag for PSA reference value identifiers
const PSARefValIDTag = 601

func init() {
	profileID, err := eat.NewProfile(ProfileURI)
	if err != nil {
		panic(err)
	}

	// Register PSA-specific ClassID factory (psa.impl-id)
	// CBOR tag 560 is already registered for TaggedBytes in comid/cbor.go
	// so we only register the factory here
	if err := comid.RegisterClassIDFactory(NewImplIDClassID); err != nil {
		panic(err)
	}

	// Register PSA-specific Mkey type (psa.refval-id) with CBOR tag 601
	// This registers both the CBOR tag and the factory
	if err := comid.RegisterMkeyType(PSARefValIDTag, NewMkeyPSARefvalID); err != nil {
		panic(err)
	}

	mvalExt := &MvalExtensions{}

	extMap := extensions.NewMap().
		Add(comid.ExtTriples, &TriplesExtensions{}).
		Add(comid.ExtReferenceValue, mvalExt).
		Add(comid.ExtEndorsedValue, mvalExt).
		Add(corim.ExtUnsignedCorim, &CorimExtensions{})

	if err := corim.RegisterProfile(profileID, extMap); err != nil {
		panic(err)
	}
}

// TriplesExtensions carries PSA-specific fields and constraints for Triples
type TriplesExtensions struct {
	PsaSwRelTriples *PsaSwRelTriples `cbor:"5,keyasint,omitempty" json:"psa-swrel-triples,omitempty"`
}

// ValidTriples implements ITriplesConstrainer to enforce PSA-specific constraints
// on the Triples structure. This is called automatically during Triples.Valid().
func (o TriplesExtensions) ValidTriples(triples *comid.Triples) error {
	// Validate Reference Values (Section 3.3)
	if triples.ReferenceValues != nil {
		refVals := (*extensions.Collection[comid.ValueTriple, *comid.ValueTriple])(triples.ReferenceValues)
		for i, refVal := range refVals.Values {
			if err := validatePSAReferenceValue(&refVal, i); err != nil {
				return err
			}
		}
	}

	// Validate Attestation Verification Keys (Section 3.4)
	if triples.AttestVerifKeys != nil {
		for i, avk := range *triples.AttestVerifKeys {
			if err := validatePSAAttestVerifKey(&avk, i); err != nil {
				return err
			}
		}
	}

	return nil
}

// validatePSAReferenceValue validates a Reference Value  of PSA Endorsements.
func validatePSAReferenceValue(refVal *comid.ValueTriple, tripleIndex int) error {
	prefix := fmt.Sprintf("reference value at index %d", tripleIndex)

	// Validate Implementation ID in Environment (Section 3.2)
	if err := validatePSAImplementationID(&refVal.Environment, prefix); err != nil {
		return err
	}

	// Access the Measurements collection
	measurements := (*extensions.Collection[comid.Measurement, *comid.Measurement])(&refVal.Measurements)

	for j, m := range measurements.Values {
		// Validate cryptokeys field (signer-id):
		// - cryptokeys (key 13) is MANDATORY
		// - Must have exactly one entry
		// - Entry must be tagged-bytes (type "bytes")
		// - Byte length must be 32, 48, or 64 (psa-hash-type)
		if err := validatePSASignerID(m.Val.CryptoKeys, tripleIndex, j); err != nil {
			return err
		}
	}

	return nil
}

// validatePSASignerID validates the cryptokeys field (signer-id)
//
//	a) CryptoKeys field is set and contains exactly one entry
//	b) The type of CryptoKey is TaggedBytes (Type() returns "bytes")
//	c) The length of the value is 32, 48, or 64
func validatePSASignerID(keys *comid.CryptoKeys, tripleIndex, measurementIndex int) error {
	prefix := fmt.Sprintf("reference value at index %d, measurement at index %d", tripleIndex, measurementIndex)

	// a) CryptoKeys field MUST be set (mandatory per PSA profile)
	if keys == nil {
		return fmt.Errorf("%s: cryptokeys (signer-id) is mandatory but not set", prefix)
	}

	// a) Must contain exactly one entry
	if len(*keys) == 0 {
		return fmt.Errorf("%s: cryptokeys must contain exactly one entry, got 0", prefix)
	}

	if len(*keys) != 1 {
		return fmt.Errorf("%s: cryptokeys must contain exactly one entry, got %d", prefix, len(*keys))
	}

	key := (*keys)[0]
	if key == nil {
		return fmt.Errorf("%s: cryptokeys entry is nil", prefix)
	}

	// b) The CryptoKey must be of type "bytes" (TaggedBytes)
	if key.Type() != "bytes" {
		return fmt.Errorf("%s: cryptokeys (signer-id) must be of type 'bytes', got '%s'", prefix, key.Type())
	}

	// c) Byte length must be 32, 48, or 64 (psa-hash-type: SHA-256, SHA-384, or SHA-512)
	tb, ok := key.Value.(comid.TaggedBytes)
	if !ok {
		tbPtr, ok := key.Value.(*comid.TaggedBytes)
		if !ok {
			return fmt.Errorf("%s: failed to extract TaggedBytes from cryptokeys", prefix)
		}
		tb = *tbPtr
	}

	switch len(tb) {
	case 32, 48, 64:
		return nil
	default:
		return fmt.Errorf("%s: signer-id must be 32, 48, or 64 bytes (got %d)", prefix, len(tb))
	}
}

// validatePSAAttestVerifKey validates Attestation Verification Keys
// The IAK public key must use pkix-base64-key-type and there must be exactly one key.
func validatePSAAttestVerifKey(avk *comid.KeyTriple, index int) error {
	prefix := fmt.Sprintf("attester verification key at index %d", index)

	// Validate Implementation ID in Environment (Section 3.2)
	if err := validatePSAImplementationID(&avk.Environment, prefix); err != nil {
		return err
	}

	// Validate Instance ID in Environment (Section 3.2)
	// Instance ID is required for Attestation Verification Keys
	if err := validatePSAInstanceID(&avk.Environment, prefix); err != nil {
		return err
	}

	// Must have exactly one key
	if len(avk.VerifKeys) == 0 {
		return fmt.Errorf("%s: verification-keys must contain exactly one entry, got 0", prefix)
	}

	if len(avk.VerifKeys) != 1 {
		return fmt.Errorf("%s: verification-keys must contain exactly one entry, got %d", prefix, len(avk.VerifKeys))
	}

	key := avk.VerifKeys[0]
	if key == nil {
		return fmt.Errorf("%s: verification-key entry is nil", prefix)
	}

	// "The IAK public key uses the tagged-pkix-base64-key-type variant"
	if key.Type() != comid.PKIXBase64KeyType {
		return fmt.Errorf("%s: verification-key must be of type '%s', got '%s'",
			prefix, comid.PKIXBase64KeyType, key.Type())
	}

	return nil
}

// validatePSAImplementationID validates the Implementation ID
// The Implementation ID must be:
//   - A tagged-bytes type (CBOR tag 560)
//   - Exactly 32 bytes in length
func validatePSAImplementationID(env *comid.Environment, prefix string) error {
	// Implementation ID is in Environment.Class.ClassID
	if env.Class == nil {
		return fmt.Errorf("%s: environment.class is required for PSA profile", prefix)
	}

	if env.Class.ClassID == nil {
		return fmt.Errorf("%s: environment.class.id (implementation-id) is required for PSA profile", prefix)
	}

	classID := env.Class.ClassID

	// Must be of type "bytes" (tagged-bytes, CBOR tag 560)
	if classID.Type() != "bytes" {
		return fmt.Errorf("%s: implementation-id must be of type 'bytes', got '%s'", prefix, classID.Type())
	}

	// Must be exactly 32 bytes
	idBytes := classID.Bytes()
	if len(idBytes) != 32 {
		return fmt.Errorf("%s: implementation-id must be exactly 32 bytes (got %d)", prefix, len(idBytes))
	}

	return nil
}

// validatePSAInstanceID validates the Instance ID (UEID) per PSA Endorsements Section 3.2
// The Instance ID must be:
//   - A tagged-ueid-type (CBOR tag 550)
//   - The first byte MUST be 0x01 (RAND type)
//   - Followed by exactly 32 bytes (total 33 bytes)
func validatePSAInstanceID(env *comid.Environment, prefix string) error {
	// Instance ID is in Environment.Instance
	if env.Instance == nil {
		return fmt.Errorf("%s: environment.instance (instance-id) is required for PSA profile", prefix)
	}

	instance := env.Instance

	// Must be of type "ueid" (tagged-ueid-type, CBOR tag 550)
	if instance.Type() != comid.UEIDType {
		return fmt.Errorf("%s: instance-id must be of type '%s', got '%s'",
			prefix, comid.UEIDType, instance.Type())
	}

	// Get the UEID bytes
	ueidBytes := instance.Bytes()

	// PSA Instance ID must be exactly 33 bytes (1 byte RAND type + 32 bytes identifier)
	if len(ueidBytes) != 33 {
		return fmt.Errorf("%s: instance-id must be exactly 33 bytes (got %d)", prefix, len(ueidBytes))
	}

	// The first byte MUST be 0x01 (RAND type per EAT UEID specification)
	if ueidBytes[0] != 0x01 {
		return fmt.Errorf("%s: instance-id must have RAND type (0x01), got 0x%02x", prefix, ueidBytes[0])
	}

	return nil
}

// MvalExtensions carries PSA-specific fields and constraints for Measurements
type MvalExtensions struct {
	PsaCertNum *string `cbor:"100,keyasint,omitempty" json:"psa-cert-num,omitempty"`
}

// CorimExtensions carries PSA-specific constraints for the top-level CoRIM
type CorimExtensions struct{}

// PsaSwRelTriples represents the psa-swrel-triples field (key 5)
type PsaSwRelTriples struct {
	// To be Implemented
}
