// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/corim"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/corim/profiles/common/arm"
)

const ProfileURI = "tag:arm.com,2025:psa#1.0.0"

// PSASoftwareComponentMkey is the required value for mkey in PSA reference value measurements
const PSASoftwareComponentMkey = "psa.software-component"

func init() {
	profileID := corim.MustNewURIProfile(ProfileURI)

	mvalExt := &MvalExtensions{}

	extMap := extensions.NewMap().
		Add(comid.ExtTriples, &TriplesExtensions{}).
		Add(comid.ExtReferenceValue, mvalExt).
		Add(comid.ExtEndorsedValue, mvalExt)

	if err := corim.RegisterProfile(profileID, extMap); err != nil {
		panic(err)
	}
}

// TriplesExtensions provides PSA-specific validation for Triples via the ValidTriples method
type TriplesExtensions struct{}

// ValidTriples implements ITriplesConstrainer to enforce PSA-specific constraints
// on the Triples structure. This is called automatically during Triples.Valid().
func (o TriplesExtensions) ValidTriples(triples *comid.Triples) error {
	// Validate Reference Values (Section 3.3)
	if triples.ReferenceValues != nil {
		for i, refVal := range triples.ReferenceValues.Values {
			if err := validatePSAReferenceValue(&refVal, i); err != nil {
				return err
			}
		}
	}

	// Validate Attestation Verification Keys (Section 3.4)
	if triples.AttestVerifKeys != nil {
		return arm.ValidateKeyTriples(triples.AttestVerifKeys)
	}

	return nil
}

// validatePSAReferenceValue validates a Reference Value  of PSA Endorsements.
func validatePSAReferenceValue(refVal *comid.ValueTriple, tripleIndex int) error {
	prefix := fmt.Sprintf("reference value at index %d", tripleIndex)

	// Validate Implementation ID in Environment (Section 3.2)
	if err := arm.ValidateImplementationID(&refVal.Environment, prefix); err != nil {
		return err
	}

	for j := range refVal.Measurements.Values {
		measurement := &refVal.Measurements.Values[j]

		// Validate mkey: must be the string "psa.software-component" (Section 3.3)
		if err := validatePSAMkey(measurement.Key, tripleIndex, j); err != nil {
			return err
		}

		if err := arm.ValidatePlatformSoftwareComponent(measurement); err != nil {
			return err
		}
	}

	return nil
}

// validatePSAMkey validates that the mkey is the string "psa.software-component"
// as required by Section 3.3 of the PSA Endorsements spec.
func validatePSAMkey(key *comid.Mkey, tripleIndex, measurementIndex int) error {
	prefix := fmt.Sprintf("reference value at index %d, measurement at index %d", tripleIndex, measurementIndex)

	// mkey is mandatory for PSA profile
	if key == nil || !key.IsSet() {
		return fmt.Errorf("%s: mkey is mandatory but not set", prefix)
	}

	// mkey must be of type "string"
	if key.Type() != comid.StringType {
		return fmt.Errorf("%s: mkey must be of type 'string', got '%s'", prefix, key.Type())
	}

	// The value must be exactly "psa.software-component"
	if key.Value.String() != PSASoftwareComponentMkey {
		return fmt.Errorf("%s: mkey must be %q, got %q",
			prefix, PSASoftwareComponentMkey, key.Value.String())
	}

	return nil
}

// MvalExtensions carries PSA-specific fields and constraints for Measurements
type MvalExtensions struct {
	PsaCertNum *string `cbor:"100,keyasint,omitempty" json:"psa-cert-num,omitempty"`
}
