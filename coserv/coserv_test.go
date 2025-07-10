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
		AddClass(StatefulClass{Class: class})
	require.NotNil(t, envSelector)

	query, err := NewQuery(ArtifactTypeReferenceValues, *envSelector, ResultTypeSourceArtifacts)
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

	expected := readTestVectorSlice(t, "rv-class-simple.cbor")
	assert.Equal(t, expected, actual)
}

func TestCoserv_ToCBOR_exampleClassSelector(t *testing.T) {
	query, err := NewQuery(ArtifactTypeReferenceValues, *exampleClassSelector(t), ResultTypeCollectedArtifacts)
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

	expected := readTestVectorSlice(t, "example-class-selector.cbor")
	assert.Equal(t, expected, actual)

	fmt.Printf("%x\n", expected)
}

func TestCoserv_ToCBOR_exampleInstanceSelector(t *testing.T) {
	query, err := NewQuery(ArtifactTypeReferenceValues, *exampleInstanceSelector(t), ResultTypeBoth)
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

	expected := readTestVectorSlice(t, "example-instance-selector.cbor")
	assert.Equal(t, expected, actual)
}

func TestCoserv_ToCBOR_exampleGroupSelector(t *testing.T) {
	query, err := NewQuery(ArtifactTypeReferenceValues, *exampleGroupSelector(t), ResultTypeSourceArtifacts)
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

	expected := readTestVectorSlice(t, "example-group-selector.cbor")
	assert.Equal(t, expected, actual)
}

func TestCoserv_FromCBOR_fail(t *testing.T) {
	tv := comid.MustHexDecode(t, "ff")

	var actual Coserv
	err := actual.FromCBOR(tv)
	assert.EqualError(t, err, `decoding CoSERV from CBOR: cbor: unexpected "break" code`)
}

func TestCoserv_FromBase64Url_ok_class(t *testing.T) {
	tv := readTestVectorString(t, "example-class-selector.b64u")

	var actual Coserv

	err := actual.FromBase64Url(tv)
	require.NoError(t, err)

	actualProfile, err := actual.Profile.Get()
	require.NoError(t, err)
	assert.Equal(t, `tag:example.com,2025:cc-platform#1.0.0`, actualProfile)
	assert.Equal(t, "reference-values", actual.Query.ArtifactType.String())
	assert.Equal(t, testTimestamp, actual.Query.Timestamp)
	assert.Equal(t, "collected-artifacts", actual.Query.ResultType.String())
	assert.Equal(t, *exampleClassSelector(t), actual.Query.EnvironmentSelector)
}

func TestCoserv_FromBase64Url_ok_instance(t *testing.T) {
	tv := readTestVectorString(t, "example-instance-selector.b64u")

	var actual Coserv

	err := actual.FromBase64Url(tv)
	require.NoError(t, err)

	actualProfile, err := actual.Profile.Get()
	require.NoError(t, err)
	assert.Equal(t, `tag:example.com,2025:cc-platform#1.0.0`, actualProfile)
	assert.Equal(t, "reference-values", actual.Query.ArtifactType.String())
	assert.Equal(t, testTimestamp, actual.Query.Timestamp)
	assert.Equal(t, "both", actual.Query.ResultType.String())
	assert.Equal(t, *exampleInstanceSelector(t), actual.Query.EnvironmentSelector)
}

func TestCoserv_FromBase64Url_ok_group(t *testing.T) {
	tv := readTestVectorString(t, "example-group-selector.b64u")

	var actual Coserv

	err := actual.FromBase64Url(tv)
	require.NoError(t, err)

	actualProfile, err := actual.Profile.Get()
	require.NoError(t, err)
	assert.Equal(t, `tag:example.com,2025:cc-platform#1.0.0`, actualProfile)
	assert.Equal(t, "reference-values", actual.Query.ArtifactType.String())
	assert.Equal(t, testTimestamp, actual.Query.Timestamp)
	assert.Equal(t, "source-artifacts", actual.Query.ResultType.String())
	assert.Equal(t, *exampleGroupSelector(t), actual.Query.EnvironmentSelector)
}

func TestCoserv_FromBase64Url_fail(t *testing.T) {
	tv := "=/+"

	var actual Coserv

	err := actual.FromBase64Url(tv)
	assert.EqualError(t, err, "decoding CoSERV: illegal base64 data at input byte 0")
}

func TestCoserv_ToBase64Url_ok_instance(t *testing.T) {
	query, err := NewQuery(ArtifactTypeReferenceValues, *exampleInstanceSelector(t), ResultTypeBoth)
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

	expected := readTestVectorString(t, "example-instance-selector.b64u")

	assert.Equal(t, expected, actual)
}

func TestCoserv_ToBase64Url_ok_class(t *testing.T) {
	query, err := NewQuery(ArtifactTypeReferenceValues, *exampleClassSelector(t), ResultTypeCollectedArtifacts)
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

	expected := readTestVectorString(t, "example-class-selector.b64u")

	assert.Equal(t, expected, actual)
}

func TestCoserv_ToBase64Url_ok_group(t *testing.T) {
	query, err := NewQuery(ArtifactTypeReferenceValues, *exampleGroupSelector(t), ResultTypeSourceArtifacts)
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

	expected := readTestVectorString(t, "example-group-selector.b64u")

	assert.Equal(t, expected, actual)
}

func TestCoserv_ToEDN_ok(t *testing.T) {
	query, err := NewQuery(ArtifactTypeReferenceValues, *exampleClassSelector(t), ResultTypeCollectedArtifacts)
	require.NoError(t, err)

	// overwrite the default query timestamp
	query.SetTimestamp(testTimestamp)

	tv, err := NewCoserv(
		`tag:example.com,2025:cc-platform#1.0.0`,
		*query,
	)
	require.NoError(t, err)

	actual, err := tv.ToEDN()
	require.NoError(t, err)

	expected := readTestVectorString(t, "example-class-selector-noindent.diag")
	assert.Equal(t, expected, actual)
}

func TestCoserv_FromCBOR_Stateful(t *testing.T) {
	tv := readTestVectorSlice(t, "rv-class-stateful.cbor")

	var actual Coserv

	err := actual.FromCBOR(tv)
	require.NoError(t, err)

	// here we only care about the measurements

	assert.Len(t, *actual.Query.EnvironmentSelector.Classes, 1)
	assert.NotNil(t, (*actual.Query.EnvironmentSelector.Classes)[0].Measurements)
}

func TestCoserv_FromCBOR_Results(t *testing.T) {
	tv := readTestVectorSlice(t, "rv-class-simple-results.cbor")

	var actual Coserv

	err := actual.FromCBOR(tv)
	require.NoError(t, err)

	actualProfile, err := actual.Profile.Get()
	require.NoError(t, err)
	assert.Equal(t, `tag:example.com,2025:cc-platform#1.0.0`, actualProfile)
	assert.Equal(t, "reference-values", actual.Query.ArtifactType.String())
	assert.Equal(t, testTimestamp, actual.Query.Timestamp)
	assert.Equal(t, *exampleClassSelector2(t), actual.Query.EnvironmentSelector)

	// results-related assertions
	assert.NotNil(t, actual.Results)
	assert.NotNil(t, actual.Results.Expiry)
	assert.Equal(t, *actual.Results.Expiry, testExpiry)
	assert.NotNil(t, actual.Results.RVQ)
	assert.Len(t, *actual.Results.RVQ, 1)

	assert.Equal(t, testBytes, (*actual.Results.RVQ)[0].RVTriple.Environment.Class.ClassID.Bytes())

	assert.Len(t, (*actual.Results.RVQ)[0].RVTriple.Measurements.Values, 2)
	assert.Equal(t, "Component A", *(*actual.Results.RVQ)[0].RVTriple.Measurements.Values[0].Val.Name)
	assert.Equal(t, "Component B", *(*actual.Results.RVQ)[0].RVTriple.Measurements.Values[1].Val.Name)

	assert.Len(t, *(*actual.Results.RVQ)[0].RVTriple.Measurements.Values[1].Val.Digests, 2)

	assert.Equal(t, (*(*actual.Results.RVQ)[0].RVTriple.Measurements.Values[0].Val.Digests)[0].HashAlgID, swid.Sha256)
	assert.Equal(t, (*(*actual.Results.RVQ)[0].RVTriple.Measurements.Values[0].Val.Digests)[0].HashValue, []byte{0xaa})
	assert.Equal(t, (*(*actual.Results.RVQ)[0].RVTriple.Measurements.Values[0].Val.Digests)[1].HashAlgID, swid.Sha256_128)
	assert.Equal(t, (*(*actual.Results.RVQ)[0].RVTriple.Measurements.Values[0].Val.Digests)[1].HashValue, []byte{0xbb})

	assert.Equal(t, (*(*actual.Results.RVQ)[0].RVTriple.Measurements.Values[1].Val.Digests)[0].HashAlgID, swid.Sha256)
	assert.Equal(t, (*(*actual.Results.RVQ)[0].RVTriple.Measurements.Values[1].Val.Digests)[0].HashValue, []byte{0xcc})
	assert.Equal(t, (*(*actual.Results.RVQ)[0].RVTriple.Measurements.Values[1].Val.Digests)[1].HashAlgID, swid.Sha256_128)
	assert.Equal(t, (*(*actual.Results.RVQ)[0].RVTriple.Measurements.Values[1].Val.Digests)[1].HashValue, []byte{0xdd})
}

func TestCoserv_FromCBOR_Results_Source_Artifacts(t *testing.T) {
	tv := readTestVectorSlice(t, "rv-class-simple-results-source-artifacts.cbor")

	var actual Coserv

	err := actual.FromCBOR(tv)
	require.NoError(t, err)

	assert.Equal(t, "source-artifacts", actual.Query.ResultType.String())

	// results-related assertions
	assert.NotNil(t, actual.Results)

	assert.NotNil(t, actual.Results.RVQ)
	assert.Len(t, *actual.Results.RVQ, 0)

	assert.NotNil(t, actual.Results.SourceArtifacts)
	assert.Len(t, *actual.Results.SourceArtifacts, 2)

	cmw0 := (*actual.Results.SourceArtifacts)[0]

	assert.Equal(t, "monad", cmw0.GetKind().String())

	t0, err := cmw0.GetMonadType()
	require.NoError(t, err)
	assert.Equal(t, "application/vnd.example.refvals", t0)

	v0, err := cmw0.GetMonadValue()
	require.NoError(t, err)
	assert.Equal(t, []byte{0xaf, 0xae, 0xad, 0xac}, v0)

	cmw1 := (*actual.Results.SourceArtifacts)[1]

	assert.Equal(t, "monad", cmw1.GetKind().String())

	t1, err := cmw1.GetMonadType()
	require.NoError(t, err)
	assert.Equal(t, "application/vnd.example.refvals", t1)

	v1, err := cmw1.GetMonadValue()
	require.NoError(t, err)
	assert.Equal(t, []byte{0xad, 0xac, 0xab, 0xaa}, v1)
}

func TestCoserv_results_ToCBOR_ok(t *testing.T) {
	class := comid.NewClassBytes(comid.TestBytes)
	require.NotNil(t, class)

	envSelector := NewEnvironmentSelector().
		AddClass(StatefulClass{Class: class})
	require.NotNil(t, envSelector)

	query, err := NewQuery(ArtifactTypeReferenceValues, *envSelector, ResultTypeCollectedArtifacts)
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

	expected := readTestVectorSlice(t, "rv-results.cbor")

	actual, err := tv.ToCBOR()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
