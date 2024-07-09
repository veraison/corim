// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package corim

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/extensions"
)

type signerExtensions struct {
	Address string `cbor:"-1" json:"address"`
}

func TestSigner_RegisterExtensions(t *testing.T) {
	signer := NewSigner()
	assert.False(t, signer.Extensions.HaveExtensions())

	exts := &signerExtensions{}
	extMap := extensions.NewMap().Add(ExtSigner, exts)

	err := signer.RegisterExtensions(extMap)
	assert.NoError(t, err)
	assert.True(t, signer.Extensions.HaveExtensions())
	assert.Equal(t, exts, signer.GetExtensions())

	badMap := extensions.NewMap().Add(extensions.Point("test"), exts)
	err = signer.RegisterExtensions(badMap)
	assert.EqualError(t, err, `unexpected extension point: "test"`)
}

func TestSigner_Valid(t *testing.T) {
	var signer Signer

	assert.EqualError(t, signer.Valid(), "empty name")

	signer.Name = "test-signer"
	uri := comid.TaggedURI("@@@")
	signer.URI = &uri

	assert.EqualError(t, signer.Valid(), `invalid URI: "@@@" is not an absolute URI`)
}

func TestSigner_JSON(t *testing.T) {
	signer := NewSigner()
	signer.Name = "test-signer"
	uri := comid.TaggedURI("https://example.com")
	signer.URI = &uri

	buf, err := signer.MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `{"name": "test-signer", "uri": "https://example.com"}`, string(buf))

	other := NewSigner()
	err = other.UnmarshalJSON(buf)
	assert.NoError(t, err)
	assert.Equal(t, signer.Name, other.Name)
	assert.Equal(t, signer.URI, other.URI)
}
