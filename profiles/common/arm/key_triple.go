// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package arm

import (
	"errors"
	"fmt"

	"github.com/veraison/corim/comid"
)

func ValidateKeyTriples(kts *comid.KeyTriples) error {
	if kts == nil {
		return errors.New("no key triples in ValidateKeyTriples")
	}

	for i, avk := range *kts {
		if err := ValidateArmPlatformAttestVerifKey(&avk); err != nil {
			return fmt.Errorf("platform attestation verification key at index %d: %w", i, err)
		}
	}

	return nil
}

// ValidateArmPlatformAttestVerifKey validates CCA Platform Attestation Verification Keys
// The CPAK public key must use pkix-base64-key-type and there must be exactly one key.
func ValidateArmPlatformAttestVerifKey(avk *comid.KeyTriple) error {
	// Validate Implementation ID in Environment (Section 3.1.2)
	if err := ValidateImplementationID(&avk.Environment, ""); err != nil {
		return fmt.Errorf("environment: %w", err)
	}

	// Validate Instance ID in Environment (Section 3.1.2)
	// Instance ID is REQUIRED for Attestation Verification Keys
	if err := ValidateInstanceID(&avk.Environment, ""); err != nil {
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
