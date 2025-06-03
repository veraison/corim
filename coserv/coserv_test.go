// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
)

func TestCoserv_ToCBOR_rv_class_simple(t *testing.T) {
	class := comid.NewClassBytes([]byte{0x00, 0x11, 0x22, 0x33}).
		SetVendor("Example Vendor").
		SetModel("Example Model")
	require.NotNil(t, class)

	envSelector := NewEnvironmentSelector().
		AddClass(*class)
	require.NotNil(t, envSelector)

	query, err := NewQuery(ArtifactTypeReferenceValues, *envSelector)
	require.NoError(t, err)

	// overwrite the default query timestamp
	query.SetTimestamp(testTimestamp)

	tv, err := NewCoserv(
		`tag:example.com,2025:cc-platform#1.0.0`,
		*query,
	)
	require.NoError(t, err)

	actual, err := tv.ToCBOR()
	assert.NoError(t, err)

	// {0: "tag:example.com,2025:cc-platform#1.0.0", 1: {0: 2, 1: {0: [{0: 560(h'00112233'), 1: "Example Vendor", 2: "Example Model"}]}, 2: 0("2025-12-13T18:30:02Z")}}
	expected := comid.MustHexDecode(t, "a20078267461673a6578616d706c652e636f6d2c323032353a63632d706c6174666f726d23312e302e3001a3000201a10081a300d902304400112233016e4578616d706c652056656e646f72026d4578616d706c65204d6f64656c02c074323032352d31322d31335431383a33303a30325a")
	assert.Equal(t, expected, actual)
}

func TestCoserv_ToCBOR_exampleClassSelector(t *testing.T) {
	query, err := NewQuery(ArtifactTypeReferenceValues, *exampleClassSelector(t))
	require.NoError(t, err)

	// overwrite the default query timestamp
	query.SetTimestamp(testTimestamp)

	tv, err := NewCoserv(
		`tag:example.com,2025:cc-platform#1.0.0`,
		*query,
	)
	require.NoError(t, err)

	actual, err := tv.ToCBOR()
	assert.NoError(t, err)

	fmt.Printf("%x\n", actual)

	// {0: "tag:example.com,2025:cc-platform#1.0.0", 1: {0: 2, 1: {0: [{0: 560(h'8999786556'), 1: "Example Vendor", 2: "Example Model"}, {0: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}]}, 2: 0("2025-12-13T18:30:02Z")}}
	expected := comid.MustHexDecode(t, "a20078267461673a6578616d706c652e636f6d2c323032353a63632d706c6174666f726d23312e302e3001a3000201a10082a300d90230458999786556016e4578616d706c652056656e646f72026d4578616d706c65204d6f64656ca100d8255031fb5abf023e4992aa4e95f9c1503bfa02c074323032352d31322d31335431383a33303a30325a")
	assert.Equal(t, expected, actual)
}

func TestCoserv_ToCBOR_exampleInstanceSelector(t *testing.T) {
	query, err := NewQuery(ArtifactTypeReferenceValues, *exampleInstanceSelector(t))
	require.NoError(t, err)

	// overwrite the default query timestamp
	query.SetTimestamp(testTimestamp)

	tv, err := NewCoserv(
		`tag:example.com,2025:cc-platform#1.0.0`,
		*query,
	)
	require.NoError(t, err)

	actual, err := tv.ToCBOR()
	assert.NoError(t, err)

	// {0: "tag:example.com,2025:cc-platform#1.0.0", 1: {0: 2, 1: {1: [550(h'02DEADBEEFDEAD'), 560(h'8999786556')]}, 2: 0("2025-12-13T18:30:02Z")}}
	expected := comid.MustHexDecode(t, "a20078267461673a6578616d706c652e636f6d2c323032353a63632d706c6174666f726d23312e302e3001a3000201a10182d902264702deadbeefdeadd9023045899978655602c074323032352d31322d31335431383a33303a30325a")
	assert.Equal(t, expected, actual)
}

func TestCoserv_ToCBOR_exampleGroupSelector(t *testing.T) {
	query, err := NewQuery(ArtifactTypeReferenceValues, *exampleGroupSelector(t))
	require.NoError(t, err)

	// overwrite the default query timestamp
	query.SetTimestamp(testTimestamp)

	tv, err := NewCoserv(
		`tag:example.com,2025:cc-platform#1.0.0`,
		*query,
	)
	require.NoError(t, err)

	actual, err := tv.ToCBOR()
	assert.NoError(t, err)

	// {0: "tag:example.com,2025:cc-platform#1.0.0", 1: {0: 2, 1: {2: [560(h'8999786556'), 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')]}, 2: 0("2025-12-13T18:30:02Z")}}
	expected := comid.MustHexDecode(t, "a20078267461673a6578616d706c652e636f6d2c323032353a63632d706c6174666f726d23312e302e3001a3000201a10282d90230458999786556d8255031fb5abf023e4992aa4e95f9c1503bfa02c074323032352d31322d31335431383a33303a30325a")
	assert.Equal(t, expected, actual)
}

func TestCoserv_FromCBOR_fail(t *testing.T) {
	tv := comid.MustHexDecode(t, "ff")

	var actual Coserv
	err := actual.FromCBOR(tv)
	assert.EqualError(t, err, `decoding CoSERV from CBOR: cbor: unexpected "break" code`)
}

func TestCoserv_FromBase64Url_ok_class(t *testing.T) {
	tv := "ogB4JnRhZzpleGFtcGxlLmNvbSwyMDI1OmNjLXBsYXRmb3JtIzEuMC4wAaMAAgGhAIKjANkCMEWJmXhlVgFuRXhhbXBsZSBWZW5kb3ICbUV4YW1wbGUgTW9kZWyhANglUDH7Wr8CPkmSqk6V-cFQO_oCwHQyMDI1LTEyLTEzVDE4OjMwOjAyWg"
	var actual Coserv

	err := actual.FromBase64Url(tv)
	require.NoError(t, err)

	actualProfile, err := actual.Profile.Get()
	require.NoError(t, err)
	assert.Equal(t, `tag:example.com,2025:cc-platform#1.0.0`, actualProfile)
	assert.Equal(t, "reference-values", actual.Query.ArtifactType.String())
	assert.Equal(t, testTimestamp, actual.Query.Timestamp)
	assert.False(t, actual.Query.GetIncludeSourceMaterial())
	assert.Equal(t, *exampleClassSelector(t), actual.Query.EnvironmentSelector)
}

func TestCoserv_FromBase64Url_ok_instance(t *testing.T) {
	tv := "ogB4JnRhZzpleGFtcGxlLmNvbSwyMDI1OmNjLXBsYXRmb3JtIzEuMC4wAaMAAgGhAYLZAiZHAt6tvu_erdkCMEWJmXhlVgLAdDIwMjUtMTItMTNUMTg6MzA6MDJa"

	var actual Coserv

	err := actual.FromBase64Url(tv)
	require.NoError(t, err)

	actualProfile, err := actual.Profile.Get()
	require.NoError(t, err)
	assert.Equal(t, `tag:example.com,2025:cc-platform#1.0.0`, actualProfile)
	assert.Equal(t, "reference-values", actual.Query.ArtifactType.String())
	assert.Equal(t, testTimestamp, actual.Query.Timestamp)
	assert.False(t, actual.Query.GetIncludeSourceMaterial())
	assert.Equal(t, *exampleInstanceSelector(t), actual.Query.EnvironmentSelector)
}

func TestCoserv_FromBase64Url_ok_group(t *testing.T) {
	tv := "ogB4JnRhZzpleGFtcGxlLmNvbSwyMDI1OmNjLXBsYXRmb3JtIzEuMC4wAaMAAgGhAoLZAjBFiZl4ZVbYJVAx-1q_Aj5JkqpOlfnBUDv6AsB0MjAyNS0xMi0xM1QxODozMDowMlo"

	var actual Coserv

	err := actual.FromBase64Url(tv)
	require.NoError(t, err)

	actualProfile, err := actual.Profile.Get()
	require.NoError(t, err)
	assert.Equal(t, `tag:example.com,2025:cc-platform#1.0.0`, actualProfile)
	assert.Equal(t, "reference-values", actual.Query.ArtifactType.String())
	assert.Equal(t, testTimestamp, actual.Query.Timestamp)
	assert.False(t, actual.Query.GetIncludeSourceMaterial())
	assert.Equal(t, *exampleGroupSelector(t), actual.Query.EnvironmentSelector)
}

func TestCoserv_FromBase64Url_ok_include_source_material(t *testing.T) {
	tv := "ogB4JnRhZzpleGFtcGxlLmNvbSwyMDI1OmNjLXBsYXRmb3JtIzEuMC4wAaQAAgGhAIKjANkCMEWJmXhlVgFuRXhhbXBsZSBWZW5kb3ICbUV4YW1wbGUgTW9kZWyhANglUDH7Wr8CPkmSqk6V-cFQO_oCwHQyMDI1LTEyLTEzVDE4OjMwOjAyWgP1"

	var actual Coserv

	err := actual.FromBase64Url(tv)
	require.NoError(t, err)

	assert.True(t, actual.Query.GetIncludeSourceMaterial())
}

func TestCoserv_FromBase64Url_fail(t *testing.T) {
	tv := "=/+"

	var actual Coserv

	err := actual.FromBase64Url(tv)
	assert.EqualError(t, err, "decoding CoSERV: illegal base64 data at input byte 0")
}

func TestCoserv_ToBase64Url_ok_instance(t *testing.T) {
	query, err := NewQuery(ArtifactTypeReferenceValues, *exampleInstanceSelector(t))
	require.NoError(t, err)

	// overwrite the default query timestamp
	query.SetTimestamp(testTimestamp)

	tv, err := NewCoserv(
		`tag:example.com,2025:cc-platform#1.0.0`,
		*query,
	)
	require.NoError(t, err)

	actual, err := tv.ToBase64Url()
	assert.NoError(t, err)

	fmt.Printf("%s\n", actual)

	expected := "ogB4JnRhZzpleGFtcGxlLmNvbSwyMDI1OmNjLXBsYXRmb3JtIzEuMC4wAaMAAgGhAYLZAiZHAt6tvu_erdkCMEWJmXhlVgLAdDIwMjUtMTItMTNUMTg6MzA6MDJa"

	assert.Equal(t, expected, actual)
}

func TestCoserv_ToBase64Url_ok_class(t *testing.T) {
	query, err := NewQuery(ArtifactTypeReferenceValues, *exampleClassSelector(t))
	require.NoError(t, err)

	// overwrite the default query timestamp
	query.SetTimestamp(testTimestamp)

	tv, err := NewCoserv(
		`tag:example.com,2025:cc-platform#1.0.0`,
		*query,
	)
	require.NoError(t, err)

	actual, err := tv.ToBase64Url()
	assert.NoError(t, err)

	fmt.Printf("%s\n", actual)

	expected := "ogB4JnRhZzpleGFtcGxlLmNvbSwyMDI1OmNjLXBsYXRmb3JtIzEuMC4wAaMAAgGhAIKjANkCMEWJmXhlVgFuRXhhbXBsZSBWZW5kb3ICbUV4YW1wbGUgTW9kZWyhANglUDH7Wr8CPkmSqk6V-cFQO_oCwHQyMDI1LTEyLTEzVDE4OjMwOjAyWg"

	assert.Equal(t, expected, actual)
}

func TestCoserv_ToBase64Url_ok_group(t *testing.T) {
	query, err := NewQuery(ArtifactTypeReferenceValues, *exampleGroupSelector(t))
	require.NoError(t, err)

	// overwrite the default query timestamp
	query.SetTimestamp(testTimestamp)

	tv, err := NewCoserv(
		`tag:example.com,2025:cc-platform#1.0.0`,
		*query,
	)
	require.NoError(t, err)

	actual, err := tv.ToBase64Url()
	assert.NoError(t, err)

	fmt.Printf("%s\n", actual)

	expected := "ogB4JnRhZzpleGFtcGxlLmNvbSwyMDI1OmNjLXBsYXRmb3JtIzEuMC4wAaMAAgGhAoLZAjBFiZl4ZVbYJVAx-1q_Aj5JkqpOlfnBUDv6AsB0MjAyNS0xMi0xM1QxODozMDowMlo"

	assert.Equal(t, expected, actual)
}

func TestCoserv_ToEDN_ok(t *testing.T) {
	query, err := NewQuery(ArtifactTypeReferenceValues, *exampleClassSelector(t))
	require.NoError(t, err)

	// overwrite the default query timestamp
	query.SetTimestamp(testTimestamp)

	// Set the include-source-material flag
	query.SetIncludeSourceMaterial()

	tv, err := NewCoserv(
		`tag:example.com,2025:cc-platform#1.0.0`,
		*query,
	)
	require.NoError(t, err)

	actual, err := tv.ToEDN()
	require.NoError(t, err)

	expected := `{0: "tag:example.com,2025:cc-platform#1.0.0", 1: {0: 2, 1: {0: [{0: 560(h'8999786556'), 1: "Example Vendor", 2: "Example Model"}, {0: 37(h'31fb5abf023e4992aa4e95f9c1503bfa')}]}, 2: 0("2025-12-13T18:30:02Z"), 3: true}}`

	assert.Equal(t, expected, actual)
}

func TestCoserv_FromCBOR_Results(t *testing.T) {
	// {0: "tag:example.com,2025:cc-platform#1.0.0", 1: {0: 2, 1: {0: [{0: 560(h'8999786556'), 1: "Example Vendor", 2: "Example Model"}, {0: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}]}, 2: 0("2025-12-13T18:30:02Z")}, 2: {10: 0("2030-12-13T18:30:02Z"), 0: [[{0: {0: 560(h'8999786556'), 1: "Example Vendor", 2: "Example Model"}}, [{1: {11: "Component A", 2: [[1, h'AA'], [2, h'BB']]}}, {1: {11: "Component B", 2: [[1, h'CC'], [2, h'DD']]}}]]]}}
	tv := comid.MustHexDecode(t, `
a30078267461673a6578616d706c652e636f6d2c323032353a63632d706c
6174666f726d23312e302e3001a3000201a10082a300d902304589997865
56016e4578616d706c652056656e646f72026d4578616d706c65204d6f64
656ca100d8255031fb5abf023e4992aa4e95f9c1503bfa02c07432303235
2d31322d31335431383a33303a30325a02a20ac074323033302d31322d31
335431383a33303a30325a008182a100a300d90230458999786556016e45
78616d706c652056656e646f72026d4578616d706c65204d6f64656c82a1
01a20b6b436f6d706f6e656e7420410282820141aa820241bba101a20b6b
436f6d706f6e656e7420420282820141cc820241dd`)

	var actual Coserv

	err := actual.FromCBOR(tv)
	require.NoError(t, err)

	actualProfile, err := actual.Profile.Get()
	require.NoError(t, err)
	assert.Equal(t, `tag:example.com,2025:cc-platform#1.0.0`, actualProfile)
	assert.Equal(t, "reference-values", actual.Query.ArtifactType.String())
	assert.Equal(t, testTimestamp, actual.Query.Timestamp)
	assert.Equal(t, *exampleClassSelector(t), actual.Query.EnvironmentSelector)

	// results-related assertions
	assert.NotNil(t, actual.Results)
	assert.NotNil(t, actual.Results.Expiry)
	assert.Equal(t, *actual.Results.Expiry, testExpiry)
	assert.NotNil(t, actual.Results.ReferenceValues)
	assert.Len(t, *actual.Results.ReferenceValues, 1)

	assert.Equal(t, comid.TestBytes, (*actual.Results.ReferenceValues)[0].Environment.Class.ClassID.Bytes())

	assert.Len(t, (*actual.Results.ReferenceValues)[0].Measurements.Values, 2)
	assert.Equal(t, "Component A", *(*actual.Results.ReferenceValues)[0].Measurements.Values[0].Val.Name)
	assert.Equal(t, "Component B", *(*actual.Results.ReferenceValues)[0].Measurements.Values[1].Val.Name)

	assert.Len(t, *(*actual.Results.ReferenceValues)[0].Measurements.Values[1].Val.Digests, 2)

	assert.Equal(t, (*(*actual.Results.ReferenceValues)[0].Measurements.Values[0].Val.Digests)[0].HashAlgID, swid.Sha256)
	assert.Equal(t, (*(*actual.Results.ReferenceValues)[0].Measurements.Values[0].Val.Digests)[0].HashValue, []byte{0xaa})
	assert.Equal(t, (*(*actual.Results.ReferenceValues)[0].Measurements.Values[0].Val.Digests)[1].HashAlgID, swid.Sha256_128)
	assert.Equal(t, (*(*actual.Results.ReferenceValues)[0].Measurements.Values[0].Val.Digests)[1].HashValue, []byte{0xbb})

	assert.Equal(t, (*(*actual.Results.ReferenceValues)[0].Measurements.Values[1].Val.Digests)[0].HashAlgID, swid.Sha256)
	assert.Equal(t, (*(*actual.Results.ReferenceValues)[0].Measurements.Values[1].Val.Digests)[0].HashValue, []byte{0xcc})
	assert.Equal(t, (*(*actual.Results.ReferenceValues)[0].Measurements.Values[1].Val.Digests)[1].HashAlgID, swid.Sha256_128)
	assert.Equal(t, (*(*actual.Results.ReferenceValues)[0].Measurements.Values[1].Val.Digests)[1].HashValue, []byte{0xdd})
}

func TestCoserv_results_ToCBOR_ok(t *testing.T) {
	class := comid.NewClassBytes(comid.TestBytes)
	require.NotNil(t, class)

	envSelector := NewEnvironmentSelector().
		AddClass(*class)
	require.NotNil(t, envSelector)

	query, err := NewQuery(ArtifactTypeReferenceValues, *envSelector)
	require.NoError(t, err)

	// overwrite the default query timestamp
	query.SetTimestamp(testTimestamp)

	tv, err := NewCoserv(
		`tag:example.com,2025:cc-platform#1.0.0`,
		*query,
	)
	require.NoError(t, err)

	err = tv.AddResults(*exampleReferenceValuesResultSet(t))
	require.NoError(t, err)

	// {
	//   0:"tag:example.com,2025:cc-platform#1.0.0",
	//   1:{0: 2, 1: {0: [{0: 560(h'8999786556')}]}, 2: 0("2025-12-13T18:30:02Z")},
	//   2:{
	//     10:0("2030-12-13T18:30:02Z"),
	//     0:[
	//       [
	//         {0: {0: 560(h'8999786556')}},
	//         [
	//           {
	//             0:37(h'31FB5ABF023E4992AA4E95F9C1503BFA'),
	//             1:{0: {0: "1.2.3", 1: 16384}, 1: 553(2)}
	//           }
	//         ]
	//       ]
	//     ]
	//   }
	// }
	expected := comid.MustHexDecode(t, `a30078267461673a6578616d706c652e636f6d2c323032353a63632d706c6174666f726d23312e302e3001a3000201a10081a100d9023045899978655602c074323032352d31322d31335431383a33303a30325a02a20ac074323033302d31322d31335431383a33303a30325a008182a100a100d9023045899978655681a200d8255031fb5abf023e4992aa4e95f9c1503bfa01a200a20065312e322e330119400001d9022902`)

	actual, err := tv.ToCBOR()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	fmt.Printf("%x\n", actual)
}
