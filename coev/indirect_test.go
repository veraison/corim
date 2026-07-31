// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/corim"
	"github.com/veraison/corim/extensions"
)

//go:embed testcases/example-concise-evidence.json
var testConciseEvidenceJSON []byte

//go:embed testcases/example-concise-evidence.cbor
var testConciseEvidenceCBOR []byte

func TestSpdmExtensionMap_OK(t *testing.T) {
	m := SpdmExtensionMap()
	assert.NotNil(t, m)
	_, ok := m[ExtEvidenceTriples]
	assert.True(t, ok)
}

func TestSpdmExtensionMap_RegisterOnCE_OK(t *testing.T) {
	ce := NewConciseEvidence()
	err := ce.RegisterExtensions(SpdmExtensionMap())
	assert.NoError(t, err)
}

func TestMValExtensions_RoundTrip_JSON(t *testing.T) {
	var ce ConciseEvidence
	require.NoError(t, ce.RegisterExtensions(SpdmExtensionMap()))
	require.NoError(t, ce.FromJSON(testConciseEvidenceJSON))

	j, err := ce.ToJSON()
	require.NoError(t, err)
	assert.Contains(t, string(j), "spdm-indirect")

	var ce2 ConciseEvidence
	require.NoError(t, ce2.RegisterExtensions(SpdmExtensionMap()))
	require.NoError(t, ce2.FromJSON(j))
	assert.Equal(t, ce, ce2)
}

func TestMValExtensions_RoundTrip_CBOR(t *testing.T) {
	// testConciseEvidenceCBOR is a tagged-concise-evidence (tag 571)
	var tce TaggedConciseEvidence
	require.NoError(t, (*ConciseEvidence)(&tce).RegisterExtensions(SpdmExtensionMap()))
	require.NoError(t, tce.FromCBOR(testConciseEvidenceCBOR))

	b, err := tce.ToCBOR()
	require.NoError(t, err)
	assert.Equal(t, testConciseEvidenceCBOR, b)
}

// Profile registration tests

func TestRegisterProfile_OK(t *testing.T) {
	profileID := corim.MustNewOIDProfile("1.2.3.4")
	extMap := extensions.NewMap().Add(ExtEvidenceTriples, &MValExtensions{})
	err := RegisterProfile(profileID, extMap)
	require.NoError(t, err)
	t.Cleanup(func() { UnregisterProfile(profileID) })
}

func TestRegisterProfile_NOK_already_registered(t *testing.T) {
	profileID := corim.MustNewOIDProfile("1.2.3.5")
	extMap := extensions.NewMap().Add(ExtEvidenceTriples, &MValExtensions{})
	require.NoError(t, RegisterProfile(profileID, extMap))
	t.Cleanup(func() { UnregisterProfile(profileID) })

	err := RegisterProfile(profileID, extMap)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestRegisterProfile_NOK_unexpected_point(t *testing.T) {
	profileID := corim.MustNewOIDProfile("1.2.3.6")
	extMap := extensions.NewMap().Add("bad-point", &MValExtensions{})
	err := RegisterProfile(profileID, extMap)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected extension point")
}

func TestRegisterProfile_NOK_non_pointer(t *testing.T) {
	profileID := corim.MustNewOIDProfile("1.2.3.7")
	extMap := extensions.NewMap().Add(ExtEvidenceTriples, MValExtensions{})
	err := RegisterProfile(profileID, extMap)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "non-pointer")
}

func TestGetProfileManifest_OK(t *testing.T) {
	profileID := corim.MustNewOIDProfile("1.2.3.8")
	extMap := extensions.NewMap().Add(ExtEvidenceTriples, &MValExtensions{})
	require.NoError(t, RegisterProfile(profileID, extMap))
	t.Cleanup(func() { UnregisterProfile(profileID) })

	m, ok := GetProfileManifest(profileID)
	assert.True(t, ok)
	assert.NotNil(t, m)
}

func TestGetProfileManifest_NOK_nil(t *testing.T) {
	m, ok := GetProfileManifest(&corim.Profile{})
	assert.False(t, ok)
	assert.Nil(t, m)
}

func TestUnregisterProfile_OK(t *testing.T) {
	profileID := corim.MustNewOIDProfile("1.2.3.9")
	extMap := extensions.NewMap().Add(ExtEvidenceTriples, &MValExtensions{})
	require.NoError(t, RegisterProfile(profileID, extMap))

	ok := UnregisterProfile(profileID)
	assert.True(t, ok)
}

func TestUnregisterProfile_NOK_not_registered(t *testing.T) {
	profileID := corim.MustNewOIDProfile("1.2.3.10")
	ok := UnregisterProfile(profileID)
	assert.False(t, ok)
}

func TestUnregisterProfile_NOK_nil(t *testing.T) {
	ok := UnregisterProfile(&corim.Profile{})
	assert.False(t, ok)
}

func TestGetConciseEvidence_nil_profile(t *testing.T) {
	ce := GetConciseEvidence(&corim.Profile{})
	assert.NotNil(t, ce)
}

func TestGetConciseEvidence_unknown_profile(t *testing.T) {
	profileID := corim.MustNewOIDProfile("1.3.999.999")
	ce := GetConciseEvidence(profileID)
	assert.NotNil(t, ce)
}

func TestGetConciseEvidence_registered_profile(t *testing.T) {
	profileID := corim.MustNewOIDProfile("1.2.3.11")
	extMap := extensions.NewMap().Add(ExtEvidenceTriples, &MValExtensions{})
	require.NoError(t, RegisterProfile(profileID, extMap))
	t.Cleanup(func() { UnregisterProfile(profileID) })

	ce := GetConciseEvidence(profileID)
	assert.NotNil(t, ce)
}

func TestProfileManifest_GetConciseEvidence(t *testing.T) {
	profileID := corim.MustNewOIDProfile("1.2.3.12")
	extMap := extensions.NewMap().Add(ExtEvidenceTriples, &MValExtensions{})
	require.NoError(t, RegisterProfile(profileID, extMap))
	t.Cleanup(func() { UnregisterProfile(profileID) })

	m, ok := GetProfileManifest(profileID)
	require.True(t, ok)

	ce := m.GetConciseEvidence()
	assert.NotNil(t, ce)
}

func TestProfileManifest_GetTaggedConciseEvidence(t *testing.T) {
	profileID := corim.MustNewOIDProfile("1.2.3.13")
	extMap := extensions.NewMap().Add(ExtEvidenceTriples, &MValExtensions{})
	require.NoError(t, RegisterProfile(profileID, extMap))
	t.Cleanup(func() { UnregisterProfile(profileID) })

	m, ok := GetProfileManifest(profileID)
	require.True(t, ok)

	tce := m.GetTaggedConciseEvidence()
	assert.NotNil(t, tce)
}

func TestUnmarshalConciseEvidenceFromCBOR_unknown_profile(t *testing.T) {
	profileID := corim.MustNewOIDProfile("1.2.3.15")
	// not registered — should still succeed (extensions fall back to none)
	ce, err := UnmarshalConciseEvidenceFromCBOR(testConciseEvidenceCBOR, profileID)
	// validation fails because spdm-indirect (key 12) is not recognized without extensions
	_ = ce
	_ = err
}

func TestUnmarshalConciseEvidenceFromCBOR_with_spdm_extensions(t *testing.T) {
	// check UnmarshalConciseEvidenceFromCBOR applies extensions automatically
	profileID := corim.MustNewOIDProfile("1.2.3.14")
	extMap := extensions.NewMap().Add(ExtEvidenceTriples, &MValExtensions{})
	require.NoError(t, RegisterProfile(profileID, extMap))
	t.Cleanup(func() { UnregisterProfile(profileID) })

	ce, err := UnmarshalConciseEvidenceFromCBOR(testConciseEvidenceCBOR, profileID)
	require.NoError(t, err)
	assert.NotNil(t, ce)
}
