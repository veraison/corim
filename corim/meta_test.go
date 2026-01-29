// Copyright 2021-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/extensions"
)

var (
	// {0: {0: "ACME Ltd.", 1: 32("https://acme.example")}, 1: {0: 1(1601424000), 1: 1(1632960000)}}
	metaFull = []byte{
		0xa2, 0x00, 0xa2, 0x00, 0x69, 0x41, 0x43, 0x4d, 0x45, 0x20, 0x4c, 0x74,
		0x64, 0x2e, 0x01, 0xd8, 0x20, 0x74, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a,
		0x2f, 0x2f, 0x61, 0x63, 0x6d, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70,
		0x6c, 0x65, 0x01, 0xa2, 0x00, 0xc1, 0x1a, 0x5f, 0x73, 0xca, 0x80, 0x01,
		0xc1, 0x1a, 0x61, 0x54, 0xfe, 0x00,
	}
	// {0: {0: "ACME Ltd."}, 1: {1: 1(1605181526)}}
	mandatoryOnly = []byte{
		0xa2, 0x00, 0xa1, 0x00, 0x69, 0x41, 0x43, 0x4d, 0x45, 0x20, 0x4c, 0x74,
		0x64, 0x2e, 0x01, 0xa1, 0x01, 0xc1, 0x1a, 0x61, 0x54, 0xfe, 0x00,
	}
)

func TestMeta_SetSigner_empty_name(t *testing.T) {
	tv := NewMeta()
	require.NotNil(t, tv)

	emptyName := ""

	assert.Nil(t, tv.SetSigner(emptyName, nil))
}

func TestMeta_SetSigner_name_only(t *testing.T) {
	actual := NewMeta()
	require.NotNil(t, actual)

	name := "ACME Ltd."

	expected := Meta{
		Signer: Signer{Name: name},
	}

	assert.NotNil(t, actual.SetSigner(name, nil))
	assert.Equal(t, expected, *actual)
}

func TestMeta_SetSigner_empty_uri(t *testing.T) {
	tv := NewMeta()
	require.NotNil(t, tv)

	emptyURI := ""

	assert.Nil(t, tv.SetSigner("ACME Ltd.", &emptyURI))
}

func TestMeta_SetSigner_bad_uri(t *testing.T) {
	tv := NewMeta()
	require.NotNil(t, tv)

	badURI := "z/a"

	assert.Nil(t, tv.SetSigner("ACME Ltd.", &badURI))
}

func TestMeta_SetSigner_full(t *testing.T) {
	actual := NewMeta()
	require.NotNil(t, actual)

	var (
		name      = "ACME Ltd."
		uri       = "https://acme.example"
		taggedURI = comid.TaggedURI(uri)
	)

	expected := Meta{
		Signer: Signer{
			Name: name,
			URI:  &taggedURI,
		},
	}

	assert.NotNil(t, actual.SetSigner(name, &uri))
	assert.Equal(t, expected, *actual)
}

func TestMeta_SetValidity_ok(t *testing.T) {
	var (
		notBefore = time.Now()
		notAfter  = time.Now().Add(time.Hour)
	)

	actual := NewMeta().
		SetValidity(notAfter, &notBefore)

	expected := Meta{
		Validity: &Validity{
			NotBefore: &notBefore,
			NotAfter:  notAfter,
		},
	}

	assert.NotNil(t, actual)
	assert.Equal(t, expected, *actual)
}

func TestMeta_Valid_ok(t *testing.T) {
	tv := NewMeta().
		SetSigner("ACME Ltd.", nil).
		SetValidity(time.Now(), nil)

	require.NotNil(t, tv)

	assert.Nil(t, tv.Valid())
}

func TestMeta_ToCBOR_mandatory_only(t *testing.T) {
	var (
		notAfter = time.Date(2021, time.October, 0, 0, 0, 0, 0, time.UTC)
	)

	tv := NewMeta().
		SetSigner("ACME Ltd.", nil).
		SetValidity(notAfter, nil)
	require.NotNil(t, tv)

	expected := mandatoryOnly

	actual, err := tv.ToCBOR()

	fmt.Printf("%x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestMeta_ToCBOR_full(t *testing.T) {
	var (
		notAfter  = time.Date(2021, time.October, 0, 0, 0, 0, 0, time.UTC)
		notBefore = time.Date(2020, time.October, 0, 0, 0, 0, 0, time.UTC)
		name      = "ACME Ltd."
		uri       = "https://acme.example"
	)

	tv := NewMeta().
		SetSigner(name, &uri).
		SetValidity(notAfter, &notBefore)
	require.NotNil(t, tv)

	expected := metaFull

	actual, err := tv.ToCBOR()

	fmt.Printf("%x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestMeta_FromCBOR_full(t *testing.T) {
	tv := metaFull

	var (
		notAfter  = time.Date(2021, time.October, 0, 0, 0, 0, 0, time.UTC)
		notBefore = time.Date(2020, time.October, 0, 0, 0, 0, 0, time.UTC)
		name      = "ACME Ltd."
		taggedURI = comid.TaggedURI("https://acme.example")
	)

	var actual Meta
	err := actual.FromCBOR(tv)

	assert.Nil(t, err)
	assert.Equal(t, name, actual.Signer.Name)
	assert.Equal(t, taggedURI, *actual.Signer.URI)
	assert.Equal(t, notBefore.Unix(), actual.Validity.NotBefore.Unix())
	assert.Equal(t, notAfter.Unix(), actual.Validity.NotAfter.Unix())
}

func Test_Meata_extensions(t *testing.T) {
	meta := NewMeta()

	exts := &signerExtensions{}
	extMap := extensions.NewMap().Add(ExtSigner, exts)

	err := meta.RegisterExtensions(extMap)
	assert.NoError(t, err)
	assert.True(t, meta.Signer.HaveExtensions())
	assert.Equal(t, exts, meta.Signer.GetExtensions())

	badMap := extensions.NewMap().Add(extensions.Point("test"), exts)
	err = meta.RegisterExtensions(badMap)
	assert.EqualError(t, err, `unexpected extension point: "test"`)
}

func Test_meta_Valid(t *testing.T) {
	meta := NewMeta()

	assert.EqualError(t, meta.Valid(), "invalid signer: empty name")

	uri := "http://example.com"
	meta.SetSigner("test", &uri)
	assert.NoError(t, meta.Valid())

	ts := time.Now()
	meta.Validity = &Validity{NotBefore: &ts}
	err := meta.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid validity: invalid not-before / not-after")
}

func TestMeta_ToJSON(t *testing.T) {
	var (
		notAfter  = time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC)
		notBefore = time.Date(2020, time.October, 1, 0, 0, 0, 0, time.UTC)
		name      = "ACME Ltd."
		uri       = "https://acme.example"
	)

	tv := NewMeta().
		SetSigner(name, &uri).
		SetValidity(notAfter, &notBefore)
	require.NotNil(t, tv)

	jsonData, err := tv.ToJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)
	assert.Contains(t, string(jsonData), "ACME Ltd.")
	assert.Contains(t, string(jsonData), "acme.example")
}

func TestMeta_FromJSON(t *testing.T) {
	jsonData := []byte(`{"signer":{"name":"ACME Ltd.","uri":"https://acme.example"},"validity":{"not-after":"2021-10-01T00:00:00Z","not-before":"2020-10-01T00:00:00Z"}}`)

	var meta Meta
	err := meta.FromJSON(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, "ACME Ltd.", meta.Signer.Name)
}

func TestMeta_FromJSON_Invalid(t *testing.T) {
	jsonData := []byte(`{invalid}`)

	var meta Meta
	err := meta.FromJSON(jsonData)
	assert.Error(t, err)
}
