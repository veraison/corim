// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	_ "embed"
	"encoding/json"
	"testing"

	cborfx "github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/corim"
)

//go:embed testcases/example-spdm-toc.json
var testSpdmTocJSON []byte

//go:embed testcases/example-spdm-toc.cbor
var testSpdmTocCBOR []byte

func TestTaggedSpdmToc_FromJSON_OK(t *testing.T) {
	var toc TaggedSpdmToc
	require.NoError(t, toc.FromJSON(testSpdmTocJSON))
	assert.Len(t, toc.TaggedEvidence, 1)
	assert.Len(t, toc.TaggedEvidence[0].EvTriples.EvidenceTriples.Values, 11)
}

func TestTaggedSpdmToc_ToCBOR_OK(t *testing.T) {
	var toc TaggedSpdmToc
	require.NoError(t, toc.FromJSON(testSpdmTocJSON))

	got, err := toc.ToCBOR()
	require.NoError(t, err)
	assert.Equal(t, testSpdmTocCBOR, got)
}

func TestTaggedSpdmToc_FromCBOR_OK(t *testing.T) {
	var toc TaggedSpdmToc
	require.NoError(t, toc.FromCBOR(testSpdmTocCBOR))
	assert.Len(t, toc.TaggedEvidence, 1)
	assert.Len(t, toc.TaggedEvidence[0].EvTriples.EvidenceTriples.Values, 11)
}

func TestTaggedSpdmToc_ToJSON_OK(t *testing.T) {
	var toc TaggedSpdmToc
	require.NoError(t, toc.FromCBOR(testSpdmTocCBOR))

	j, err := toc.ToJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, j)
}

func TestTaggedSpdmToc_RoundTrip_JSON_CBOR(t *testing.T) {
	var toc TaggedSpdmToc
	require.NoError(t, toc.FromJSON(testSpdmTocJSON))

	cbor, err := toc.ToCBOR()
	require.NoError(t, err)

	var toc2 TaggedSpdmToc
	require.NoError(t, toc2.FromCBOR(cbor))

	j, err := toc2.ToJSON()
	require.NoError(t, err)

	var toc3 TaggedSpdmToc
	require.NoError(t, toc3.FromJSON(j))
	assert.Equal(t, len(toc.TaggedEvidence), len(toc3.TaggedEvidence))
}

func TestTaggedSpdmToc_FromCBOR_NOK_missing_tag(t *testing.T) {
	var toc TaggedSpdmToc
	err := toc.FromCBOR([]byte{0x00, 0x01, 0x02})
	assert.EqualError(t, err, "missing spdm-toc tag 570")
}

func TestTaggedSpdmToc_FromCBOR_NOK_bad_inner_cbor(t *testing.T) {
	var toc TaggedSpdmToc
	// Tag 570 followed by invalid CBOR
	data := append(append([]byte{}, SpdmTocTag...), 0xff, 0xff)
	err := toc.FromCBOR(data)
	require.Error(t, err)
}

func TestTaggedSpdmToc_FromCBOR_NOK_missing_inner_tag571(t *testing.T) {
	// Build a tagged-spdm-toc-map containing one item that lacks Tag 571.
	// FromCBOR delegates to TaggedConciseEvidence.FromCBOR which checks the tag.
	type spdmTocMap struct {
		Items []cborfx.RawMessage `cbor:"0,keyasint"`
	}
	inner, err := cborfx.Marshal(spdmTocMap{Items: []cborfx.RawMessage{{0xa0}}})
	require.NoError(t, err)

	data := append(append([]byte{}, SpdmTocTag...), inner...)
	var toc TaggedSpdmToc
	err = toc.FromCBOR(data)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tagged-concise-evidence 0")
}

func TestTaggedSpdmToc_FromCBOR_NOK_invalid_ce_cbor(t *testing.T) {
	// Build a tagged-spdm-toc-map where the inner item an empty tagged-concise-evidence.
	// This is valid CBOR but fails ConciseEvidence.Valid()
	type spdmTocMap struct {
		Items []cborfx.RawMessage `cbor:"0,keyasint"`
	}
	badCE := append(append([]byte{}, ConciseEvidenceTag...), 0xa0) // Tag 571 + empty map
	inner, err := cborfx.Marshal(spdmTocMap{Items: []cborfx.RawMessage{badCE}})
	require.NoError(t, err)

	data := append(append([]byte{}, SpdmTocTag...), inner...)
	var toc TaggedSpdmToc
	err = toc.FromCBOR(data)
	require.Error(t, err)
}

func TestTaggedSpdmToc_FromJSON_NOK_bad_json(t *testing.T) {
	var toc TaggedSpdmToc
	err := toc.FromJSON([]byte("{"))
	assert.Error(t, err)
}

func TestTaggedSpdmToc_FromJSON_NOK_invalid_ce(t *testing.T) {
	// tagged-evidence contains an empty object — CE Valid() will fail
	var toc TaggedSpdmToc
	err := toc.FromJSON([]byte(`{"tagged-evidence":[{}]}`))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tagged-concise-evidence 0")
}

func TestTaggedSpdmToc_RimLocators_RoundTrip(t *testing.T) {
	href := "https://example.com/rim.cbor"
	locators := []corim.Locator{
		{Href: corim.OneOrMore[comid.TaggedURI]{comid.TaggedURI(href)}},
	}

	var toc TaggedSpdmToc
	require.NoError(t, toc.FromJSON(testSpdmTocJSON))
	toc.RimLocators = &locators

	// JSON round-trip
	j, err := toc.ToJSON()
	require.NoError(t, err)
	assert.Contains(t, string(j), "rim-locators")
	assert.Contains(t, string(j), href)

	var toc2 TaggedSpdmToc
	require.NoError(t, toc2.FromJSON(j))
	require.NotNil(t, toc2.RimLocators)
	require.Len(t, *toc2.RimLocators, 1)

	// CBOR round-trip
	b, err := toc.ToCBOR()
	require.NoError(t, err)

	var toc3 TaggedSpdmToc
	require.NoError(t, toc3.FromCBOR(b))
	require.NotNil(t, toc3.RimLocators)
	assert.Len(t, *toc3.RimLocators, 1)
}

func TestTaggedSpdmToc_Profile_RoundTrip(t *testing.T) {
	profile := corim.MustNewOIDProfile("1.3.6.1.4.1.5703.1300.1.1")

	var toc TaggedSpdmToc
	require.NoError(t, toc.FromJSON(testSpdmTocJSON))
	toc.Profile = profile

	// JSON round-trip
	j, err := toc.ToJSON()
	require.NoError(t, err)
	var probed struct {
		Profile json.RawMessage `json:"profile"`
	}
	require.NoError(t, json.Unmarshal(j, &probed))
	assert.NotEmpty(t, probed.Profile)

	var toc2 TaggedSpdmToc
	require.NoError(t, toc2.FromJSON(j))
	require.NotNil(t, toc2.Profile)
	assert.Equal(t, profile.String(), toc2.Profile.String())

	// CBOR round-trip
	b, err := toc.ToCBOR()
	require.NoError(t, err)

	var toc3 TaggedSpdmToc
	require.NoError(t, toc3.FromCBOR(b))
	require.NotNil(t, toc3.Profile)
	assert.Equal(t, profile.String(), toc3.Profile.String())
}

func TestTaggedSpdmToc_OptionalFields_Absent(t *testing.T) {
	// Baseline: neither rim-locators nor profile set
	var toc TaggedSpdmToc
	require.NoError(t, toc.FromJSON(testSpdmTocJSON))
	assert.Nil(t, toc.RimLocators)
	assert.Nil(t, toc.Profile)

	b, err := toc.ToCBOR()
	require.NoError(t, err)
	assert.Equal(t, testSpdmTocCBOR, b)
}
