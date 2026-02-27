// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cca

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
// using the CCA Platform profile extensions
func getComidFromCorim(t *testing.T, corimData []byte) *comid.Comid {
	// First, parse the CoRIM to extract the tags
	var c corim.UnsignedCorim
	err := c.FromCBOR(corimData)
	require.NoError(t, err, "failed to parse CoRIM")
	require.Greater(t, len(c.Tags), 0, "CoRIM must have at least one tag")

	// Get a CoMID with CCA Platform profile extensions registered
	profileID, err := eat.NewProfile(PlatformProfileURI)
	require.NoError(t, err)

	manifest, found := corim.GetProfileManifest(profileID)
	require.True(t, found, "CCA Platform profile should be registered")

	// Create a Comid with extensions
	m := manifest.GetComid()

	// Decode the first tag (which should be a CoMID, tag 506)
	require.Equal(t, uint64(506), c.Tags[0].Number, "first tag should be a CoMID (506)")
	err = m.FromCBOR(c.Tags[0].Content)
	require.NoError(t, err, "failed to decode CoMID from tag")

	return m
}

// getComidFromCorimWithProfile extracts and decodes the first CoMID from a CoRIM
// using the specified profile extensions
func getComidFromCorimWithProfile(t *testing.T, corimData []byte, profileURI string) *comid.Comid {
	// First, parse the CoRIM to extract the tags
	var c corim.UnsignedCorim
	err := c.FromCBOR(corimData)
	require.NoError(t, err, "failed to parse CoRIM")
	require.Greater(t, len(c.Tags), 0, "CoRIM must have at least one tag")

	// Get a CoMID with the specified profile extensions registered
	profileID, err := eat.NewProfile(profileURI)
	require.NoError(t, err)

	manifest, found := corim.GetProfileManifest(profileID)
	require.True(t, found, "profile %s should be registered", profileURI)

	// Create a Comid with extensions
	m := manifest.GetComid()

	// Decode the first tag (which should be a CoMID, tag 506)
	require.Equal(t, uint64(506), c.Tags[0].Number, "first tag should be a CoMID (506)")
	err = m.FromCBOR(c.Tags[0].Content)
	require.NoError(t, err, "failed to decode CoMID from tag")

	return m
}

// TestCoRIMIntegration_ValidCCAPlatform tests that a valid CCA Platform CoRIM:
// - Parses successfully
// - The CCA Platform profile is recognized and loaded
// - CoMID validation passes with CCA Platform constraints
func TestCoRIMIntegration_ValidCCAPlatform(t *testing.T) {
	data := loadTestcase(t, "cca-platform-valid.cbor")

	// Extract and decode CoMID with CCA Platform extensions
	m := getComidFromCorim(t, data)

	// Validate the CoMID (should pass all CCA Platform constraints)
	err := m.Valid()
	assert.NoError(t, err, "valid CCA Platform CoMID should pass validation")
}

// TestCoRIMIntegration_InvalidImplementationID tests that a CoRIM with invalid
// Implementation ID (wrong length) is rejected by CCA Platform profile validation.
// The CoRIM structure itself is valid, but fails CCA Platform profile constraints.
func TestCoRIMIntegration_InvalidImplementationID(t *testing.T) {
	data := loadTestcase(t, "cca-platform-invalid-impl-id.cbor")

	// Extract and decode CoMID with CCA Platform extensions
	m := getComidFromCorim(t, data)

	// Validation should FAIL due to invalid Implementation ID
	err := m.Valid()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "implementation-id must be exactly 32 bytes")
}

// TestCoRIMIntegration_InvalidDigests tests that a CoRIM with invalid
// digests (wrong size) is rejected by CCA Platform profile validation.
func TestCoRIMIntegration_InvalidDigests(t *testing.T) {
	data := loadTestcase(t, "cca-platform-invalid-digests.cbor")

	// Extract and decode CoMID with CCA Platform extensions
	m := getComidFromCorim(t, data)

	// Validation should FAIL due to invalid digest size
	err := m.Valid()
	assert.Error(t, err)
	// The error could come from either base validation or CCA profile validation
	assert.True(t,
		strings.Contains(err.Error(), "length mismatch") ||
			strings.Contains(err.Error(), "must be 32, 48, or 64 bytes"),
		"error should mention digest size validation failure")
}

// TestCoRIMIntegration_NoProfile tests that a CoRIM without a profile:
// - Parses successfully
// - Validation passes without CCA Platform extensions (base validation only)
// This serves as a control case - the same CoMID that fails with CCA Platform profile
// should pass without it.
func TestCoRIMIntegration_NoProfile(t *testing.T) {
	data := loadTestcase(t, "no-profile.cbor")

	// Parse the CoRIM
	var c corim.UnsignedCorim
	err := c.FromCBOR(data)
	require.NoError(t, err, "profile-less CoRIM should parse without error")

	// Verify no profile is set
	assert.Nil(t, c.Profile, "profile should not be present")

	// Decode CoMID WITHOUT CCA Platform extensions (plain CoMID)
	require.Greater(t, len(c.Tags), 0, "CoRIM must have at least one tag")
	m := comid.NewComid()
	err = m.FromCBOR(c.Tags[0].Content)
	require.NoError(t, err, "failed to decode CoMID")

	// Validation should pass (no CCA Platform profile constraints applied)
	err = m.Valid()
	assert.NoError(t, err, "profile-less CoMID should pass base validation")
}

// TestCoRIMIntegration_ValidCCARealm tests that a valid CCA Realm CoRIM:
// - Parses successfully
// - The CCA Realm profile is recognized and loaded
// - CoMID validation passes with CCA Realm constraints
// - Includes RIM (mandatory), REMs, and RPV (optional)
func TestCoRIMIntegration_ValidCCARealm(t *testing.T) {
	data := loadTestcase(t, "cca-realm-valid.cbor")

	// Extract and decode CoMID with CCA Realm extensions
	m := getComidFromCorimWithProfile(t, data, RealmProfileURI)

	// Validate the CoMID (should pass all CCA Realm constraints)
	err := m.Valid()
	assert.NoError(t, err, "valid CCA Realm CoMID should pass validation")
}

// TestCoRIMIntegration_InvalidCCARealm_InvalidRIMSize tests that a CoRIM with invalid
// RIM size (wrong length) is rejected by CCA Realm profile validation.
// The CoRIM structure itself is valid, but fails CCA Realm profile constraints.
func TestCoRIMIntegration_InvalidCCARealm_InvalidRIMSize(t *testing.T) {
	data := loadTestcase(t, "cca-realm-invalid-rim-size.cbor")

	// Extract and decode CoMID with CCA Realm extensions
	m := getComidFromCorimWithProfile(t, data, RealmProfileURI)

	// Validation should FAIL due to invalid RIM size
	err := m.Valid()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "RIM must be 32, 48, or 64 bytes")
}

// TestCoRIMIntegration_InvalidCCARealm_MissingRIM tests that a CoRIM without
// the mandatory RIM measurement is rejected by CCA Realm profile validation.
// RIM is mandatory, so CoRIM must contain it.
func TestCoRIMIntegration_InvalidCCARealm_MissingRIM(t *testing.T) {
	data := loadTestcase(t, "cca-realm-invalid-missing-rim.cbor")

	// Extract and decode CoMID with CCA Realm extensions
	m := getComidFromCorimWithProfile(t, data, RealmProfileURI)

	// Validation should FAIL due to missing RIM measurement
	err := m.Valid()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "RIM (cca.rim) measurement is mandatory")
}
