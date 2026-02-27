// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cca

import (
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/corim"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/corim/profiles/psa"
	"github.com/veraison/eat"
)

const PlatformProfileURI = "tag:arm.com,2025:cca_platform#1.0.0"

// CCA Platform measurement key constants
const (
	CCASoftwareComponentMkey = "cca.software-component"
	CCAPlatformConfigMkey    = "cca.platform-config"
)

func init() {
	profileID, err := eat.NewProfile(PlatformProfileURI)
	if err != nil {
		panic(err)
	}

	extMap := extensions.NewMap().
		Add(comid.ExtTriples, &PlatformTriplesExtensions{})

	if err := corim.RegisterProfile(profileID, extMap); err != nil {
		panic(err)
	}
}

// PlatformTriplesExtensions provides CCA Platform-specific validation for Triples
type PlatformTriplesExtensions struct{}

// ValidTriples implements ITriplesConstrainer to enforce CCA Platform-specific constraints
// on the Triples structure. This is called automatically during Triples.Valid().
func (o PlatformTriplesExtensions) ValidTriples(triples *comid.Triples) error {
	// Validate Reference Values (Section 3.1.3)
	if triples.ReferenceValues != nil {
		for i, refVal := range triples.ReferenceValues.Values {
			if err := validateCCAPlatformReferenceValue(&refVal); err != nil {
				return fmt.Errorf("platform reference value at index %d: %w", i, err)
			}
		}
	}

	// Validate Attestation Verification Keys (Section 3.1.4)
	if triples.AttestVerifKeys != nil {
		for i, avk := range *triples.AttestVerifKeys {
			if err := validateCCAPlatformAttestVerifKey(&avk); err != nil {
				return fmt.Errorf("platform attestation verification key at index %d: %w", i, err)
			}
		}
	}

	return nil
}

// validateCCAPlatformReferenceValue validates a Reference Value of CCA Platform Endorsements.
func validateCCAPlatformReferenceValue(refVal *comid.ValueTriple) error {
	// Validate Implementation ID in Environment (Section 3.1.2)
	if err := validateCCAPlatformImplementationID(&refVal.Environment); err != nil {
		return fmt.Errorf("environment: %w", err)
	}

	// Track what we find
	var hasSoftwareComponent bool
	platformConfigCount := 0

	for j := range refVal.Measurements.Values {
		measurement := &refVal.Measurements.Values[j]

		// Validate mkey is set
		if measurement.Key == nil || !measurement.Key.IsSet() {
			return fmt.Errorf("measurement at index %d: mkey is mandatory but not set", j)
		}

		// mkey must be string type
		if measurement.Key.Type() != comid.StringType {
			return fmt.Errorf("measurement at index %d: mkey must be of type 'string', got '%s'",
				j, measurement.Key.Type())
		}

		mkeyVal := measurement.Key.Value.String()

		switch mkeyVal {
		case CCASoftwareComponentMkey:
			// Validate software component
			if err := validateCCASoftwareComponent(measurement); err != nil {
				return fmt.Errorf("measurement at index %d: %w", j, err)
			}
			hasSoftwareComponent = true

		case CCAPlatformConfigMkey:
			// Validate platform configuration
			if err := validateCCAPlatformConfig(measurement); err != nil {
				return fmt.Errorf("measurement at index %d: %w", j, err)
			}
			platformConfigCount++
			if platformConfigCount > 1 {
				return fmt.Errorf("only one platform-config measurement allowed per triple, found %d",
					platformConfigCount)
			}

		default:
			return fmt.Errorf("measurement at index %d: invalid mkey %q, expected %q or %q",
				j, mkeyVal, CCASoftwareComponentMkey, CCAPlatformConfigMkey)
		}
	}

	// At least one software component is required
	if !hasSoftwareComponent {
		return fmt.Errorf("at least one software component measurement is required")
	}

	// platform-config is optional, but only one entry is allowed.
	if platformConfigCount > 1 {
		return fmt.Errorf("only one platform-config measurement allowed per triple, found %d",
			platformConfigCount)
	}

	return nil
}

// validateCCASoftwareComponent validates a CCA Platform software component measurement
func validateCCASoftwareComponent(measurement *comid.Measurement) error {
	// Validate digests (key 2) - mandatory
	if err := validateCCADigests(measurement.Val.Digests); err != nil {
		return fmt.Errorf("digests: %w", err)
	}

	// Validate cryptokeys (key 13) - mandatory, exactly one signer-id
	if err := validateCCASignerID(measurement.Val.CryptoKeys); err != nil {
		return fmt.Errorf("cryptokeys: %w", err)
	}

	// version (key 0) is optional - if present, version-scheme MUST NOT be present
	if measurement.Val.Ver != nil {
		if measurement.Val.Ver.Scheme.String() != "" {
			return fmt.Errorf("version-scheme field MUST NOT be present in cca.software-component")
		}
	}
	return nil
}

// validateCCADigests validates the digests field (key 2)
func validateCCADigests(digests *comid.Digests) error {
	// Digests field is mandatory
	if digests == nil {
		return fmt.Errorf("digests field is mandatory but not set")
	}

	// Must have at least one digest
	if len(*digests) == 0 {
		return fmt.Errorf("digests must contain at least one entry")
	}

	// Validate each digest entry
	for i, digest := range *digests {
		// For CCA, we accept both integer and text algorithm IDs
		// The hash value size is what matters most for validation

		// Hash value must be 32, 48, or 64 bytes (SHA-256, SHA-384, SHA-512)
		if err := ValidateHashDigestSize(digest.HashValue); err != nil {
			return fmt.Errorf("digest at index %d: %w", i, err)
		}
	}

	return nil
}

// validateCCASignerID validates the cryptokeys field (signer-id) - key 13
//
//	a) CryptoKeys field is set and contains exactly one entry
//	b) The type of CryptoKey is TaggedBytes (Type() returns "bytes")
//	c) The length of the value is 32, 48, or 64 bytes
func validateCCASignerID(keys *comid.CryptoKeys) error {
	return psa.ValidateSignerID(keys, 0, 0)
}

// validateCCAPlatformConfig validates a CCA Platform configuration measurement
func validateCCAPlatformConfig(measurement *comid.Measurement) error {
	// raw-value (key 4) and raw-value-mask (key 5) are mandatory for platform config
	if measurement.Val.RawValue == nil {
		return fmt.Errorf("raw-value is mandatory for cca.platform-config")
	}

	if measurement.Val.RawValueMask == nil {
		return fmt.Errorf("raw-value-mask is mandatory for cca.platform-config")
	}

	// Validate we can extract bytes from raw-value
	_, err := measurement.Val.RawValue.GetBytes()
	if err != nil {
		return fmt.Errorf("unable to extract bytes from raw-value: %w", err)
	}

	return nil
}

// validateCCAPlatformAttestVerifKey validates CCA Platform Attestation Verification Keys
// The CPAK public key must use pkix-base64-key-type and there must be exactly one key.
func validateCCAPlatformAttestVerifKey(avk *comid.KeyTriple) error {
	// Validate Implementation ID in Environment (Section 3.1.2)
	if err := validateCCAPlatformImplementationID(&avk.Environment); err != nil {
		return fmt.Errorf("environment: %w", err)
	}

	// Validate Instance ID in Environment (Section 3.1.2)
	// Instance ID is REQUIRED for Attestation Verification Keys
	if err := validateCCAPlatformInstanceID(&avk.Environment); err != nil {
		return fmt.Errorf("environment: %w", err)
	}

	// Must have exactly one key
	if len(avk.VerifKeys) == 0 {
		return fmt.Errorf("verification-keys must contain exactly one entry, got 0")
	}

	if len(avk.VerifKeys) != 1 {
		return fmt.Errorf("verification-keys must contain exactly one entry, got %d", len(avk.VerifKeys))
	}

	key := avk.VerifKeys[0]
	if key == nil {
		return fmt.Errorf("verification-key entry is nil")
	}

	// "The CPAK public key uses the tagged-pkix-base64-key-type variant" (CBOR tag 554)
	if key.Type() != comid.PKIXBase64KeyType {
		return fmt.Errorf("verification-key must be of type '%s', got '%s'",
			comid.PKIXBase64KeyType, key.Type())
	}

	return nil
}

// validateCCAPlatformImplementationID validates the CCA Platform Implementation ID
// The Implementation ID must be:
//   - A tagged-bytes type (CBOR tag 560)
//   - Exactly 32 bytes in length
func validateCCAPlatformImplementationID(env *comid.Environment) error {
	return psa.ValidateImplementationID(env, "")
}

// validateCCAPlatformInstanceID validates the CCA Platform Instance ID per Section 3.1.2
// The Instance ID must be:
//   - A tagged-ueid-type (CBOR tag 550)
//   - The first byte MUST be 0x01 (RAND type)
//   - Followed by exactly 32 bytes (total 33 bytes)
func validateCCAPlatformInstanceID(env *comid.Environment) error {
	return psa.ValidateInstanceID(env, "")
}
