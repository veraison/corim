// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"fmt"
	"regexp"

	"github.com/veraison/corim/comid"
)

// PSASwComponentVersion represents the version information for a PSA software component
type PSASwComponentVersion struct {
	Version string `cbor:"0,keyasint" json:"version"`
}

// PSADigest represents a PSA digest with algorithm and value
type PSADigest struct {
	Algorithm string `cbor:"0,keyasint" json:"alg"`
	Value     []byte `cbor:"1,keyasint" json:"val"`
}

// PSASwComponentMeasurementValues represents the measurement values for a PSA software component
// as defined in draft-fdb-rats-psa-endorsements-08 Section 3.3
type PSASwComponentMeasurementValues struct {
	Version    *PSASwComponentVersion `cbor:"0,keyasint,omitempty" json:"version,omitempty"`
	Digests    []PSADigest            `cbor:"2,keyasint" json:"digests"`
	Name       *string                `cbor:"11,keyasint,omitempty" json:"name,omitempty"`
	CryptoKeys [][]byte               `cbor:"13,keyasint" json:"cryptokeys"`
}

// Valid validates the PSA software component measurement values according to the specification
func (o PSASwComponentMeasurementValues) Valid() error {
	// Digests field is mandatory and must contain at least one entry
	if len(o.Digests) == 0 {
		return fmt.Errorf("digests field is mandatory and must contain at least one entry")
	}

	// Check that all digest algorithms are unique
	algs := make(map[string]bool)
	for _, digest := range o.Digests {
		if algs[digest.Algorithm] {
			return fmt.Errorf("duplicate digest algorithm: %s", digest.Algorithm)
		}
		algs[digest.Algorithm] = true

		// Validate hash length based on algorithm
		if err := validateHashLength(digest.Algorithm, digest.Value); err != nil {
			return err
		}
	}

	// CryptoKeys field is mandatory and must contain exactly one entry
	if len(o.CryptoKeys) != 1 {
		return fmt.Errorf("cryptokeys field is mandatory and must contain exactly one entry")
	}

	// Validate signer ID length (32, 48, or 64 bytes)
	signerID := o.CryptoKeys[0]
	switch len(signerID) {
	case 32, 48, 64:
		// Valid lengths
	default:
		return fmt.Errorf("signer-id must be 32, 48, or 64 bytes, got %d", len(signerID))
	}

	return nil
}

// validateHashLength validates that the hash value length matches the expected length for the algorithm
func validateHashLength(algorithm string, value []byte) error {
	var expectedLength int
	switch algorithm {
	case "sha-256":
		expectedLength = 32
	case "sha-384":
		expectedLength = 48
	case "sha-512":
		expectedLength = 64
	default:
		// For unknown algorithms, allow any length
		return nil
	}

	if len(value) != expectedLength {
		return fmt.Errorf("invalid hash length for %s: expected %d bytes, got %d", 
			algorithm, expectedLength, len(value))
	}

	return nil
}

// PSASwComponentMeasurement represents a PSA software component measurement
// with the mkey set to "psa.software-component"
type PSASwComponentMeasurement struct {
	comid.Measurement
}

// NewPSASwComponentMeasurement creates a new PSA software component measurement
// This is a placeholder that demonstrates the concept - full integration
// with the comid.Measurement system requires additional work
func NewPSASwComponentMeasurement(values *PSASwComponentMeasurementValues) (*PSASwComponentMeasurement, error) {
	if err := values.Valid(); err != nil {
		return nil, fmt.Errorf("invalid PSA software component values: %w", err)
	}

	// For now, just return a simple measurement structure
	// Full integration would require proper CBOR/JSON marshaling
	measurement := &PSASwComponentMeasurement{}

	return measurement, nil
}

// PSACertNumType represents a PSA Certified Security Assurance Certificate number
// The format is validated according to the specification: "[0-9]{13} - [0-9]{5}"
type PSACertNumType string

// Valid validates the PSA certificate number format
func (o PSACertNumType) Valid() error {
	pattern := `^[0-9]{13} - [0-9]{5}$`
	matched, err := regexp.MatchString(pattern, string(o))
	if err != nil {
		return fmt.Errorf("failed to validate certificate number pattern: %w", err)
	}
	if !matched {
		return fmt.Errorf("invalid PSA certificate number format: must match pattern '[0-9]{13} - [0-9]{5}', got '%s'", string(o))
	}
	return nil
}
