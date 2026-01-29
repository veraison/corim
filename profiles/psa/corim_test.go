// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/corim"
	"github.com/veraison/eat"
)

// getTestcasePath returns the absolute path to a testcase file
func getTestcasePath(filename string) string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get current file path")
	}
	testcasesDir := filepath.Join(filepath.Dir(thisFile), "testcases")
	return filepath.Join(testcasesDir, filename)
}

// loadTestcase loads a CBOR testcase file
func loadTestcase(t *testing.T, filename string) []byte {
	// Validate filename is relative (no absolute paths or parent directory traversal)
	if filepath.IsAbs(filename) || strings.Contains(filename, "..") {
		t.Fatalf("invalid testcase filename: %s", filename)
	}
	path := getTestcasePath(filename)
	data, err := os.ReadFile(path) // #nosec G304 - path is validated above
	require.NoError(t, err, "failed to load testcase %s", filename)
	return data
}

// getComidFromCorim extracts and decodes the first CoMID from a CoRIM
// using the PSA profile extensions
func getComidFromCorim(t *testing.T, corimData []byte) *comid.Comid {
	// First, parse the CoRIM to extract the tags
	var c corim.UnsignedCorim
	err := c.FromCBOR(corimData)
	require.NoError(t, err, "failed to parse CoRIM")
	require.Greater(t, len(c.Tags), 0, "CoRIM must have at least one tag")

	// Get a CoMID with PSA profile extensions registered
	profileID, err := eat.NewProfile(ProfileURI)
	require.NoError(t, err)

	manifest, found := corim.GetProfileManifest(profileID)
	require.True(t, found, "PSA profile should be registered")

	// Create a Comid with extensions
	m := manifest.GetComid()

	// Decode the first tag (which should be a CoMID, tag 506)
	require.Equal(t, uint64(506), c.Tags[0].Number, "first tag should be a CoMID (506)")
	err = m.FromCBOR(c.Tags[0].Content)
	require.NoError(t, err, "failed to decode CoMID from tag")

	return m
}

// TestCoRIMIntegration_ValidPSA tests that a valid PSA CoRIM:
// - Parses successfully
// - The PSA profile is recognized and loaded
// - CoMID validation passes with PSA constraints
func TestCoRIMIntegration_ValidPSA(t *testing.T) {
	data := loadTestcase(t, "psa-valid.cbor")

	// Extract and decode CoMID with PSA extensions
	m := getComidFromCorim(t, data)

	// Validate the CoMID (should pass all PSA constraints)
	err := m.Valid()
	assert.NoError(t, err, "valid PSA CoMID should pass validation")
}

// TestCoRIMIntegration_InvalidImplementationID tests that a CoRIM with invalid
// Implementation ID (wrong length) is rejected by PSA profile validation.
// The CoRIM structure itself is valid, but fails PSA profile constraints.
func TestCoRIMIntegration_InvalidImplementationID(t *testing.T) {
	data := loadTestcase(t, "psa-invalid-impl-id.cbor")

	// Extract and decode CoMID with PSA extensions
	m := getComidFromCorim(t, data)

	// Validation should FAIL due to invalid Implementation ID
	err := m.Valid()
	assert.ErrorContains(t, err, "implementation-id must be exactly 32 bytes")
}

// TestCoRIMIntegration_InvalidAttestVerifKey tests that a CoRIM with invalid
// AttestVerifKeys (multiple keys instead of one) is rejected by PSA profile validation.
// The CoRIM structure itself is valid, but fails PSA profile constraints.
func TestCoRIMIntegration_InvalidAttestVerifKey(t *testing.T) {
	data := loadTestcase(t, "psa-invalid-attest-key.cbor")

	// Extract and decode CoMID with PSA extensions
	m := getComidFromCorim(t, data)

	// Validation should FAIL due to multiple attestation keys
	err := m.Valid()
	assert.ErrorContains(t, err, "verification-keys must contain exactly one entry")
}

// TestCoRIMIntegration_NoProfile tests that a CoRIM without a profile:
// - Parses successfully
// - Validation passes without PSA extensions (base validation only)
// This serves as a control case - the same CoMID that fails with PSA profile
// should pass without it.
func TestCoRIMIntegration_NoProfile(t *testing.T) {
	data := loadTestcase(t, "no-profile.cbor")

	// Parse the CoRIM
	var c corim.UnsignedCorim
	err := c.FromCBOR(data)
	require.NoError(t, err, "profile-less CoRIM should parse without error")

	// Verify no profile is set
	assert.Nil(t, c.Profile, "profile should not be present")

	// Decode CoMID WITHOUT PSA extensions (plain CoMID)
	require.Greater(t, len(c.Tags), 0, "CoRIM must have at least one tag")
	m := comid.NewComid()
	err = m.FromCBOR(c.Tags[0].Content)
	require.NoError(t, err, "failed to decode CoMID")

	// Validation should pass (no PSA profile constraints applied)
	err = m.Valid()
	assert.NoError(t, err, "profile-less CoMID should pass base validation")
}
