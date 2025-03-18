// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package corim

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
)

func TestProfileManifest_registration(t *testing.T) {
	exts := extensions.NewMap()

	err := RegisterProfile(&eat.Profile{}, exts)
	assert.EqualError(t, err, "no valid EAT profile")

	p1, err := eat.NewProfile("1.2.3")
	require.NoError(t, err)

	err = RegisterProfile(p1, exts)
	assert.NoError(t, err)

	p2, err := eat.NewProfile("1.2.3")
	require.NoError(t, err)

	err = RegisterProfile(p2, exts)
	assert.EqualError(t, err, `profile with id "1.2.3" already registered`)

	ret := UnregisterProfile(p2)
	assert.True(t, ret)
	ret = UnregisterProfile(p2)
	assert.False(t, ret)
	ret = UnregisterProfile(nil)
	assert.False(t, ret)

	err = RegisterProfile(p2, exts)
	assert.NoError(t, err)

	prof, ok := GetProfileManifest(p1)
	assert.True(t, ok)
	assert.Equal(t, exts, prof.MapExtensions)

	_, ok = GetProfileManifest(&eat.Profile{})
	assert.False(t, ok)

	p3, err := eat.NewProfile("2.3.4")
	require.NoError(t, err)

	exts2 := extensions.NewMap().Add(extensions.Point("test"), &struct{}{})
	err = RegisterProfile(p3, exts2)
	assert.EqualError(t, err, `unexpected extension point: "test"`)

	exts3 := extensions.NewMap().Add(ExtEntity, struct{}{})
	err = RegisterProfile(p3, exts3)
	assert.EqualError(t, err, `attempting to register a non-pointer IMapValue for "CorimEntity"`)

	UnregisterProfile(p1)
}

func TestProfileManifest_getters(t *testing.T) {
	id, err := eat.NewProfile("1.2.3")
	require.NoError(t, err)

	profileManifest := ProfileManifest{
		ID: id,
		MapExtensions: extensions.NewMap().
			Add(comid.ExtComid, &struct{}{}).
			Add(ExtUnsignedCorim, &struct{}{}).
			Add(ExtSigner, &struct{}{}),
	}

	c := profileManifest.GetComid()
	assert.NotNil(t, c.Extensions.IMapValue)

	u := profileManifest.GetUnsignedCorim()
	assert.NotNil(t, u.Extensions.IMapValue)

	s := profileManifest.GetSignedCorim()
	assert.NotNil(t, s.UnsignedCorim.Extensions.IMapValue)
	assert.NotNil(t, s.Meta.Signer.Extensions.IMapValue)
}

func TestProfileManifest_marshaling(t *testing.T) {
	type corimExtensions struct {
		Extension1 *string `cbor:"-1,keyasint,omitempty" json:"ext1,omitempty"`
	}

	type entityExtensions struct {
		Address *string `cbor:"-1,keyasint,omitempty" json:"address,omitempty"`
	}

	type refValExtensions struct {
		Timestamp *int `cbor:"-1,keyasint,omitempty" json:"timestamp,omitempty"`
	}

	profID, err := eat.NewProfile("http://example.com/test-profile")
	require.NoError(t, err)

	extMap := extensions.NewMap().
		Add(ExtUnsignedCorim, &corimExtensions{}).
		Add(comid.ExtEntity, &entityExtensions{}).
		Add(comid.ExtReferenceValue, &refValExtensions{})
	err = RegisterProfile(profID, extMap)
	require.NoError(t, err)

	c, err := UnmarshalUnsignedCorimFromCBOR(testGoodUnsignedCorimCBOR)
	assert.NoError(t, err)
	assert.Nil(t, c.Profile)

	c, err = UnmarshalUnsignedCorimFromCBOR(testUnsignedCorimWithExtensionsCBOR)
	assert.NoError(t, err)

	assert.Equal(t, profID, c.Profile)
	assert.Equal(t, "foo", c.Extensions.MustGetString("Extension1"))

	profileManifest, ok := GetProfileManifest(c.Profile)
	assert.True(t, ok)

	cmd, err := UnmarshalComidFromCBOR(c.Tags[0], c.Profile)
	assert.NoError(t, err)

	address := cmd.Entities.Values[0].Extensions.MustGetString("Address")
	assert.Equal(t, "123 Fake Street", address)

	ts := cmd.Triples.ReferenceValues.Values[0].Measurements.Values[0].
		Val.Extensions.MustGetInt("timestamp")
	assert.Equal(t, 1720782190, ts)

	unregProfID, err := eat.NewProfile("http://example.com")
	require.NoError(t, err)

	cmdNoExt, err := UnmarshalComidFromCBOR(c.Tags[0], unregProfID)
	assert.NoError(t, err)

	address = cmdNoExt.Entities.Values[0].Extensions.MustGetString("Address")
	assert.Equal(t, "", address)

	out, err := c.ToCBOR()
	assert.NoError(t, err)
	assertCoRIMEq(t, testUnsignedCorimWithExtensionsCBOR, out)

	out, err = cmd.ToCBOR()
	assert.NoError(t, err)
	// the first 3 bytes in Tags[0] is the tag indicating CoRIM
	assertCBOREq(t, c.Tags[0][3:], out)

	c, err = UnmarshalUnsignedCorimFromJSON(testUnsignedCorimJSON)
	assert.NoError(t, err)
	assert.Nil(t, c.Profile)

	c, err = UnmarshalUnsignedCorimFromJSON(testUnsignedCorimWithExtensionsJSON)
	assert.NoError(t, err)

	assert.Equal(t, profID, c.Profile)
	assert.Equal(t, "foo", c.Extensions.MustGetString("Extension1"))

	cmd = profileManifest.GetComid()
	err = cmd.FromJSON(testComidJSON)
	assert.NoError(t, err)

	cmd = profileManifest.GetComid()
	err = cmd.FromJSON(testComidWithExtensionsJSON)
	assert.NoError(t, err)

	address = cmd.Entities.Values[0].Extensions.MustGetString("Address")
	assert.Equal(t, "123 Fake Street", address)

	ts = cmd.Triples.ReferenceValues.Values[0].Measurements.Values[0].
		Val.Extensions.MustGetInt("timestamp")
	assert.Equal(t, 1720782190, ts)

	s, err := UnmarshalSignedCorimFromCBOR(testGoodSignedCorimCBOR)
	assert.NoError(t, err)
	assert.Nil(t, s.UnsignedCorim.Profile)

	s, err = UnmarshalSignedCorimFromCBOR(testSignedCorimWithExtensionsCBOR)
	assert.NoError(t, err)

	assert.Equal(t, profID, s.UnsignedCorim.Profile)
	assert.Equal(t, "foo", s.UnsignedCorim.Extensions.MustGetString("Extension1"))

	UnregisterProfile(profID)
}
