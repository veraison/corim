// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cca

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/corim"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/corim/profiles/common/arm"
)

const PlatformProfileURI = "tag:arm.com,2025:cca_platform#1.0.0"

// CCA Platform measurement key constants
const (
	CCASoftwareComponentMkey          = "cca.software-component"
	CCAPlatformConfigMkey             = "cca.platform-config"
	CCAPlatformManufacturingConfigKey = "cca.platform-manufacturing-config"
	CCAMaxTBBRoTPKKeys                = 48 // In CCA Max key supported is 48
)

func init() {
	profileID := corim.MustNewURIProfile(PlatformProfileURI)

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
		return arm.ValidateKeyTriples(triples.AttestVerifKeys)
	}

	return nil
}

// validateCCAPlatformReferenceValue validates a Reference Value of CCA Platform Endorsements.
func validateCCAPlatformReferenceValue(refVal *comid.ValueTriple) error {
	// Validate Implementation ID in Environment (Section 3.1.2)
	if err := arm.ValidateImplementationID(&refVal.Environment, ""); err != nil {
		return fmt.Errorf("environment: %w", err)
	}

	// Track what we find
	var hasSoftwareComponent bool
	platformConfigCount := 0
	platformManufacturingConfigCount := 0
	rotpkRegexp := regexp.MustCompile(
		`^cca\.rotpk\.[CD]M\.[0-5]$`,
	)

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
			if err := arm.ValidatePlatformSoftwareComponent(measurement); err != nil {
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

		case CCAPlatformManufacturingConfigKey:
			// Validate platform manufacturing configuration
			if err := validateCCAPlatformManufacturingConfig(measurement); err != nil {
				return fmt.Errorf("measurement at index %d: %w", j, err)
			}
			platformManufacturingConfigCount++
			if platformManufacturingConfigCount > 1 {
				return fmt.Errorf("only one platform-manufacturing-config measurement allowed per triple, found %d",
					platformManufacturingConfigCount)
			}

		default:

			// Check for TBB ROTPK in the MKey
			if rotpkRegexp.MatchString(mkeyVal) {
				if err := validateCCATBBRoTPK(measurement.Val.CryptoKeys); err != nil {
					return fmt.Errorf("measurement at index %d: %w", j, err)
				} else {
					continue
				}
			}

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

// validateCCAPlatformConfig validates a CCA Platform configuration measurement
func validateCCAPlatformConfig(measurement *comid.Measurement) error {
	// raw-value (key 4), with typechoice masked-raw-value is mandatory
	if err := validateCCARawValue(measurement.Val.RawValue); err != nil {
		return fmt.Errorf("cca.platform-config raw value error: %w", err)
	}

	return nil
}

// validateCCAPlatformManufacturingConfig validates a CCA Platform manufacturing configuration measurement
func validateCCAPlatformManufacturingConfig(measurement *comid.Measurement) error {
	// raw-value (key 4), with typechoice masked-raw-value is mandatory
	if err := validateCCARawValue(measurement.Val.RawValue); err != nil {
		return fmt.Errorf("cca.platform-manufacturing-config raw value error: %w", err)
	}
	return nil
}

func validateCCATBBRoTPK(keys *comid.CryptoKeys) error {

	if keys == nil {
		return fmt.Errorf("missing crypto keys for TBB RoTPK reference value")
	}

	if len(*keys) == 0 {
		return errors.New("no keys present")
	}

	// In CCA Implementation, for a particular type of array, the maximum number
	//  of array entries SHALL be 8,
	// while the maximum entries in a single array SHALL be 6.
	// So we can have at most 48 Key Entries
	if len(*keys) > CCAMaxTBBRoTPKKeys {
		return fmt.Errorf("invalid number of keys in measurement maps: %d", len(*keys))
	}

	for i, key := range *keys {
		if err := key.Valid(); err != nil {
			return fmt.Errorf("invalid key at index %d: %w", i, err)
		}

		if key.Type() != comid.BytesType {
			return fmt.Errorf("must be of type 'bytes' %d:, %s", i, key.Type())
		}

		b := key.Value.Bytes()

		// Key Digest is always CCA Hash value that must be 32, 48, or 64 bytes (SHA-256, SHA-384, SHA-512)
		if err := arm.ValidateHashDigestSize(b); err != nil {
			return fmt.Errorf("invalid key size at index %d: %w", i, err)
		}

	}
	return nil
}

// validateCCARawValue is a common function applicable to validation of CCA Platform Configuration raw-value
// as well as CCA Platform Manufacturing Configuration raw-value
func validateCCARawValue(rv *comid.RawValue) error {
	if rv == nil {
		return fmt.Errorf("raw-value is mandatory")
	}

	if rv.Type() != comid.MaskedType {
		return fmt.Errorf("raw-value must be of type masked")
	}

	if len(rv.Bytes()) == 0 {
		return fmt.Errorf("raw-value cannot be empty")
	}

	if len(rv.Mask()) == 0 {
		return fmt.Errorf("raw-value mask cannot be empty")
	}
	if len(rv.Bytes()) != len(rv.Mask()) {
		return fmt.Errorf("length of raw-value bytes: %d, not equal to length of raw-value mask: %d", len(rv.Bytes()), len(rv.Mask()))
	}

	return nil
}
