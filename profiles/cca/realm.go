// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cca

import (
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/corim"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
)

const RealmProfileURI = "tag:arm.com,2025:cca_realm#1.0.0"

// CCA Realm measurement key constants
const (
	CCARealmInitialMeasurementMkey   = "cca.rim"
	CCARealmPersonalizationMkey      = "cca.rpv"
	CCARealmExtendedMeasurement0Mkey = "cca.rem0"
	CCARealmExtendedMeasurement1Mkey = "cca.rem1"
	CCARealmExtendedMeasurement2Mkey = "cca.rem2"
	CCARealmExtendedMeasurement3Mkey = "cca.rem3"
)

func init() {
	profileID, err := eat.NewProfile(RealmProfileURI)
	if err != nil {
		panic(err)
	}

	extMap := extensions.NewMap().
		Add(comid.ExtTriples, &RealmTriplesExtensions{})

	if err := corim.RegisterProfile(profileID, extMap); err != nil {
		panic(err)
	}
}

// RealmTriplesExtensions provides CCA Realm-specific validation for Triples
type RealmTriplesExtensions struct{}

// ValidTriples implements ITriplesConstrainer to enforce CCA Realm-specific constraints
// on the Triples structure. This is called automatically during Triples.Valid().
func (o RealmTriplesExtensions) ValidTriples(triples *comid.Triples) error {
	if triples.ReferenceValues != nil {
		for i, refVal := range triples.ReferenceValues.Values {
			if err := validateCCARealmReferenceValue(&refVal); err != nil {
				return fmt.Errorf("realm reference value at index %d: %w", i, err)
			}
		}
	}

	return nil
}

// validateCCARealmReferenceValue validates a Reference Value of CCA Realm Endorsements.
// Reference Values comprise:
// 1. Realm Initial Measurements (RIM) - MANDATORY
// 2. Realm Extended Measurements (REMs) - OPTIONAL
// 3. Realm Personalization Value (RPV) - OPTIONAL
func validateCCARealmReferenceValue(refVal *comid.ValueTriple) error {
	// Validate RIM as class identifier in Environment (Section 3.2.1)
	if err := validateCCARealmRIM(&refVal.Environment); err != nil {
		return err
	}

	// Track what we find
	var hasRIM bool
	remPresent := make(map[string]bool)

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
		case CCARealmInitialMeasurementMkey:
			// Validate RIM (mandatory digest)
			if err := validateCCARealmMeasurement(measurement); err != nil {
				return fmt.Errorf("RIM measurement at index %d: %w", j, err)
			}
			hasRIM = true

		case CCARealmExtendedMeasurement0Mkey, CCARealmExtendedMeasurement1Mkey,
			CCARealmExtendedMeasurement2Mkey, CCARealmExtendedMeasurement3Mkey:
			// Validate REM (optional digest)
			if err := validateCCARealmMeasurement(measurement); err != nil {
				return fmt.Errorf("REM measurement at index %d: %w", j, err)
			}
			remPresent[mkeyVal] = true

		case CCARealmPersonalizationMkey:
			// Validate RPV (optional raw-value using tagged bytes)
			if err := validateCCARealmPersonalizationValue(measurement); err != nil {
				return fmt.Errorf("RPV measurement at index %d: %w", j, err)
			}

		default:
			return fmt.Errorf("measurement at index %d: invalid mkey %q, expected 'cca.rim', 'cca.rem0'-'cca.rem3', or 'cca.rpv'",
				j, mkeyVal)
		}
	}

	// RIM is mandatory
	if !hasRIM {
		return fmt.Errorf("RIM (cca.rim) measurement is mandatory but not found")
	}

	return nil
}

// validateCCARealmRIM validates that the RIM (Realm Initial Measurement) is present
// and correctly encoded as a class identifier using tagged-bytes in the environment.
// The RIM is the unique identifier for a Realm Target Environment.
func validateCCARealmRIM(env *comid.Environment) error {
	// RIM is in Environment.Class.ClassID as tagged-bytes
	if env.Class == nil {
		return fmt.Errorf("environment.class is required for CCA Realm profile")
	}

	if env.Class.ClassID == nil {
		return fmt.Errorf("environment.class.id (RIM) is required for CCA Realm profile")
	}

	classID := env.Class.ClassID

	// Must be of type "bytes" (tagged-bytes, CBOR tag 560)
	if classID.Type() != "bytes" {
		return fmt.Errorf("RIM must be of type 'bytes', got '%s'", classID.Type())
	}

	// Get the RIM bytes for validation
	rimBytes := classID.Bytes()

	// RIM should be a valid hash digest (32, 48, or 64 bytes for SHA-256, SHA-384, SHA-512)
	if err := ValidateHashDigestSize(rimBytes); err != nil {
		return fmt.Errorf("RIM must be 32, 48, or 64 bytes")
	}

	return nil
}

// validateCCARealmMeasurement validates a RIM or REM (Realm measurement) which must have
// a digest (key 2) encoded as text algorithm and hash value.
func validateCCARealmMeasurement(measurement *comid.Measurement) error {
	// RIM and REM must have digests (key 2) - mandatory
	if measurement.Val.Digests == nil {
		return fmt.Errorf("digests field is mandatory but not set")
	}

	digests := measurement.Val.Digests

	// Must have exactly one digest (one entry per measurement)
	if len(*digests) == 0 {
		return fmt.Errorf("digests must contain at least one entry")
	}

	if len(*digests) > 1 {
		return fmt.Errorf("digests must contain exactly one entry, got %d", len(*digests))
	}

	// Validate the single digest entry
	digest := (*digests)[0]

	// Hash value must be 32, 48, or 64 bytes (SHA-256, SHA-384, SHA-512)
	if err := ValidateHashDigestSize(digest.HashValue); err != nil {
		return fmt.Errorf("hash value must be 32, 48, or 64 bytes")
	}

	return nil
}

// validateCCARealmPersonalizationValue validates the RPV (Realm Personalization Value)
// which is encoded as a raw-value using the tagged-bytes variant.
func validateCCARealmPersonalizationValue(measurement *comid.Measurement) error {
	// RPV must have a raw-value (key 4)
	if measurement.Val.RawValue == nil {
		return fmt.Errorf("raw-value is mandatory for cca.rpv")
	}

	// Validate we can extract bytes from raw-value
	_, err := measurement.Val.RawValue.GetBytes()
	if err != nil {
		return fmt.Errorf("unable to extract bytes from raw-value: %w", err)
	}

	return nil
}
