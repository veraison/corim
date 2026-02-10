// Copyright 2021-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/cots"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/swid"
)

func TestUnsignedCorim_id_string(t *testing.T) {
	testIDString := "test string"

	tv := NewUnsignedCorim().SetID(testIDString)
	require.NotNil(t, tv)

	actual := tv.GetID()
	assert.Equal(t, testIDString, actual)
}

func TestUnsignedCorim_id_string_empty(t *testing.T) {
	emptyString := ""

	tv := NewUnsignedCorim()
	require.NotNil(t, tv)

	assert.Nil(t, tv.SetID(emptyString))
}

func TestUnsignedCorim_id_uuid(t *testing.T) {
	tv := NewUnsignedCorim().SetID(comid.TestUUIDString)
	require.NotNil(t, tv)

	actual := tv.GetID()
	assert.Equal(t, comid.TestUUIDString, actual)
}

func TestUnsignedCorim_id_uuid_empty(t *testing.T) {
	emptyUUID := []byte{}

	tv := NewUnsignedCorim()
	require.NotNil(t, tv)

	assert.Nil(t, tv.SetID(emptyUUID))
}

func TestUnsignedCorim_AddComid_and_marshal(t *testing.T) {
	// Build a CoMID with generic types (UUID for class-id and measurement key)
	c := comid.NewComid().
		SetLanguage("en-GB").
		SetTagIdentity("43BBE37F-2E61-4B33-AED3-53CFF1428B16", 0).
		AddEntity("ACME Ltd.", &comid.TestRegID, comid.RoleCreator, comid.RoleTagCreator, comid.RoleMaintainer).
		AddReferenceValue(
			&comid.ValueTriple{
				Environment: comid.Environment{
					Class: comid.NewClassUUID(comid.TestUUID).
						SetVendor("ACME").
						SetModel("RoadRunner"),
				},
				Measurements: *comid.NewMeasurements().
					Add(
						comid.MustNewUUIDMeasurement(comid.TestUUID).
							AddDigest(1, comid.MustHexDecode(t, "87428fc522803d31065e7bce3cf03fe475096631e5e07bbd7a0fde60c4cf25c7")),
					),
			},
		)
	require.NotNil(t, c)

	tv := NewUnsignedCorim().SetID("test corim id")
	require.NotNil(t, tv)

	assert.NotNil(t, tv.AddComid(c))

	actual, err := tv.ToCBOR()
	assert.Nil(t, err)

	fmt.Printf("CBOR: %x\n", actual)

	// Verify the result can be decoded back
	var decoded UnsignedCorim
	err = decoded.FromCBOR(actual)
	assert.Nil(t, err)
	assert.Equal(t, "test corim id", decoded.GetID())
	assert.Equal(t, 1, len(decoded.Tags))
}

func TestUnsignedCorim_AddCots_and_marshal(t *testing.T) {
	tv := NewUnsignedCorim().SetID("test corim id with CoTS")
	require.NotNil(t, tv)

	c := &cots.ConciseTaStore{}

	err := c.FromJSON([]byte(cots.ConciseTaStoreTemplateSingleOrg))
	require.Nil(t, err)
	assert.NotNil(t, tv.AddCots(c))

	actual, err := tv.ToCBOR()
	assert.Nil(t, err)

	fmt.Printf("CBOR: %x", actual)

	expected := comid.MustHexDecode(t, "d901f5a200777465737420636f72696d206964207769746820436f54530181d901fb5896a301a20050ab0f44b1bfdc4604ab4a30f80407ebcc01050281a101a100a10173576f7274686c657373205365612c20496e632e06a100818202585b3059301306072a8648ce3d020106082a8648ce3d03010703420004ad8a0c01da9eda0253dc2bc27227d9c7213df8df13e89cb9cdb7a8e4b62d9ce8a99a2d705c0f7f80db65c006d1091422b47fc611cbd46869733d9c483884d5fe")

	assert.Equal(t, expected, actual)
}

func TestUnsignedCorim_AddCoswid_and_marshal(t *testing.T) {
	tv := NewUnsignedCorim().SetID("test corim id with CoSWID")
	require.NotNil(t, tv)

	var c swid.SoftwareIdentity

	data := []byte(`<SoftwareIdentity xmlns="http://standards.iso.org/iso/19770/-2/2015/schema.xsd" tagId="com.acme.rrd2013-ce-sp1-v4-1-5-0" name="ACME Roadrunner Detector 2013 Coyote Edition SP1" version="4.1.5"><Meta activationStatus="trial" colloquialVersion="2013" edition="coyote" product="Roadrunner Detector" revision="sp1"></Meta><Entity name="The ACME Corporation" regid="acme.com" role="tagCreator softwareCreator"></Entity><Entity name="Coyote Services, Inc." regid="mycoyote.com" role="distributor"></Entity><Link href="www.gnu.org/licenses/gpl.txt" rel="license"></Link><Payload><Directory name="rrdetector" root="%programdata%"><File name="rrdetector.exe" size="532712" hash="sha-256:oxT8LcZjrnpra8Z4dZQFc5bms/VpzVD9XdtNG7r9K2o="></File></Directory></Payload></SoftwareIdentity>`)
	err := c.FromXML(data)
	require.Nil(t, err)

	assert.NotNil(t, tv.AddCoswid(&c))

	actual, err := tv.ToCBOR()
	assert.Nil(t, err)

	fmt.Printf("CBOR: %x", actual)

	expected := comid.MustHexDecode(t, "d901f5a20078197465737420636f72696d206964207769746820436f535749440181d901f9590179a8007820636f6d2e61636d652e727264323031332d63652d7370312d76342d312d352d300c0001783041434d4520526f616472756e6e6572204465746563746f72203230313320436f796f74652045646974696f6e205350310d65342e312e3505a5182b65747269616c182d6432303133182f66636f796f7465183473526f616472756e6e6572204465746563746f721836637370310282a3181f745468652041434d4520436f72706f726174696f6e18206861636d652e636f6d1821820102a3181f75436f796f74652053657276696365732c20496e632e18206c6d79636f796f74652e636f6d18210404a21826781c7777772e676e752e6f72672f6c6963656e7365732f67706c2e7478741828676c6963656e736506a110a318186a72726465746563746f7218196d2570726f6772616d6461746125181aa111a318186e72726465746563746f722e657865141a000820e80782015820a314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6a")

	assert.Equal(t, expected, actual)
}

func TestUnsignedCorim_unmarshal(t *testing.T) {
	tv := testGoodUnsignedCorimCBOR

	var unsignedCorim UnsignedCorim

	err := unsignedCorim.FromCBOR(tv)
	assert.Nil(t, err)

	assert.Nil(t, unsignedCorim.Valid())

	expectedID := unsignedCorim.GetID()
	assert.Equal(t, expectedID, "test corim id")

	assert.NotNil(t, unsignedCorim.Tags)
	assert.Equal(t, 1, len(unsignedCorim.Tags))

	var c comid.Comid
	err = c.FromCBOR(unsignedCorim.Tags[0].Content)
	assert.Nil(t, err)
}

func TestUnsignedCorim_Valid_no_id(t *testing.T) {
	tv := NewUnsignedCorim()
	require.NotNil(t, tv)

	expectedError := "empty id"

	err := tv.Valid()

	assert.EqualError(t, err, expectedError)
}

func TestUnsignedCorim_Valid_no_tags(t *testing.T) {
	tv := NewUnsignedCorim().SetID("no.tags.corim")
	require.NotNil(t, tv)

	expectedError := "tags validation failed: no tags"

	err := tv.Valid()

	assert.EqualError(t, err, expectedError)
}

func TestUnsignedCorim_Valid_invalid_tags(t *testing.T) {
	tv := NewUnsignedCorim().SetID("invalid.tags.corim")
	require.NotNil(t, tv)

	tv.Tags = append(tv.Tags, Tag{Number: 1337, Content: []byte{}})

	expectedError := "tag validation failed at pos 0: empty tag"

	err := tv.Valid()

	assert.EqualError(t, err, expectedError)
}

func TestUnsignedCorim_Valid_ok(t *testing.T) {
	// minimalist CoMID
	c := comid.NewComid().
		SetTagIdentity("vendor.example/prod/1", 0).
		AddAttestVerifKey(
			&comid.KeyTriple{
				Environment: comid.Environment{
					Instance: comid.MustNewUUIDInstance(comid.TestUUID),
				},
				VerifKeys: *comid.NewCryptoKeys().
					Add(
						comid.MustNewPKIXBase64Key(comid.TestECPubKey),
					),
			},
		)
	require.NotNil(t, c)

	tv := NewUnsignedCorim().
		SetID("invalid.tags.corim").
		AddDependentRim("http://endorser.example/addon.corim", nil).
		SetProfile("https://arm.com/psa/iot/2.0.0").
		AddComid(c).
		SetRimValidity(time.Now().Add(time.Hour), nil).
		AddEntity("ACME Ltd.", nil, RoleManifestCreator)

	require.NotNil(t, tv)

	err := tv.Valid()

	assert.Nil(t, err)
}

func TestUnsignedCorim_SetRimValidity_invalid(t *testing.T) {
	notBefore := time.Now().Add(time.Hour)
	notAfter := time.Now()

	tv := NewUnsignedCorim().
		SetRimValidity(notAfter, &notBefore)

	assert.Nil(t, tv)
}

func TestUnsignedCorim_SetRimValidity_full(t *testing.T) {
	notBefore := time.Now()
	notAfter := time.Now().Add(time.Hour)

	actual := NewUnsignedCorim().
		SetRimValidity(notAfter, &notBefore)

	expected := UnsignedCorim{
		RimValidity: &Validity{
			NotBefore: &notBefore,
			NotAfter:  notAfter,
		},
	}

	assert.NotNil(t, actual)
	assert.Equal(t, expected, *actual)
}

func TestUnsignedCorim_SetRimValidity_no_optional_not_before(t *testing.T) {
	notAfter := time.Now().Add(time.Hour)

	actual := NewUnsignedCorim().
		SetRimValidity(notAfter, nil)

	expected := UnsignedCorim{
		RimValidity: &Validity{
			NotBefore: nil,
			NotAfter:  notAfter,
		},
	}

	assert.NotNil(t, actual)
	assert.Equal(t, expected, *actual)
}

func TestUnsignedCorim_AddEntity_full(t *testing.T) {
	name := "ACME Ltd."
	role := RoleManifestCreator
	regID := "https://acme.example"
	taggedRegID := comid.TaggedURI(regID)

	actual := NewUnsignedCorim().
		AddEntity(name, &regID, role)

	expected := UnsignedCorim{
		Entities: NewEntities().Add(&Entity{
			Name:  MustNewStringEntityName(name),
			Roles: Roles{role},
			RegID: &taggedRegID,
		}),
	}

	assert.NotNil(t, actual)
	assert.Equal(t, expected, *actual)
}

func TestUnsignedCorim_AddEntity_unknown_role(t *testing.T) {
	tv := NewUnsignedCorim().
		AddEntity("ACME Ltd.", nil, Role(666))

	assert.Nil(t, tv)
}

func TestUnsignedCorim_AddEntity_empty_entity_name(t *testing.T) {
	anEmptyName := ""

	tv := NewUnsignedCorim().
		AddEntity(anEmptyName, nil, RoleManifestCreator)

	assert.Nil(t, tv)
}

func TestUnsignedCorim_AddEntity_non_nil_empty_URI(t *testing.T) {
	anEmptyURI := ""

	tv := NewUnsignedCorim().
		AddEntity("ACME Ltd.", &anEmptyURI, RoleManifestCreator)

	assert.Nil(t, tv)
}

func TestUnsignedCorim_FromJSON(t *testing.T) {
	data := []byte(`{"corim-id": "5c57e8f4-46cd-421b-91c9-08cf93e13cfc"}`)

	err := NewUnsignedCorim().FromJSON(data)

	assert.NoError(t, err)
}

func TestUnsignedCorim_ToJSON(t *testing.T) {
	c := comid.NewComid().
		SetTagIdentity("vendor.example/prod/1", 0).
		AddAttestVerifKey(
			&comid.KeyTriple{
				Environment: comid.Environment{
					Instance: comid.MustNewUUIDInstance(comid.TestUUID),
				},
				VerifKeys: *comid.NewCryptoKeys().
					Add(
						comid.MustNewPKIXBase64Key(comid.TestECPubKey),
					),
			},
		)
	require.NotNil(t, c)

	tv := NewUnsignedCorim().
		SetID("invalid.tags.corim").
		AddDependentRim("http://endorser.example/addon.corim", nil).
		SetProfile("https://arm.com/psa/iot/2.0.0").
		AddComid(c)

	require.NotNil(t, tv)

	extMap := extensions.NewMap().Add(ExtEntity, &struct{}{})
	err := tv.RegisterExtensions(extMap)
	assert.NoError(t, err)

	buf, err := tv.ToJSON()

	assert.NoError(t, err)
	expectedJSON := `
	{
		"corim-id":"invalid.tags.corim",
		"tags":["2QH6WOuiAaEAdXZlbmRvci5leGFtcGxlL3Byb2QvMQShA4GCoQHYJVAx+1q/Aj5JkqpOlfnBUDv6gdkCKnixLS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFVzFCdnFGKy9yeThCV2E3WkVNVTF4WVlIRVE4QgpsTFQ0TUZIT2FPK0lDVHRJdnJFZUVwci9zZlRBUDY2SDJoQ0hkYjVIRVhLdFJLb2Q2UUxjT0xQQTFRPT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0t"],
		"dependent-rims":[{"href":"http://endorser.example/addon.corim"}],
		"profile":"https://arm.com/psa/iot/2.0.0"
	}
	`

	assert.JSONEq(t, expectedJSON, string(buf))
}

func TestUnsignedCorim_ToCBOR(t *testing.T) {
	c := comid.NewComid().
		SetTagIdentity("vendor.example/prod/1", 0).
		AddAttestVerifKey(
			&comid.KeyTriple{
				Environment: comid.Environment{
					Instance: comid.MustNewUUIDInstance(comid.TestUUID),
				},
				VerifKeys: *comid.NewCryptoKeys().
					Add(
						comid.MustNewPKIXBase64Key(comid.TestECPubKey),
					),
			},
		)
	require.NotNil(t, c)

	tv := NewUnsignedCorim().
		SetID("invalid.tags.corim").
		AddDependentRim("http://endorser.example/addon.corim", nil).
		SetProfile("https://arm.com/psa/iot/2.0.0").
		AddComid(c)

	require.NotNil(t, tv)

	extMap := extensions.NewMap().Add(ExtEntity, &struct{}{})
	err := tv.RegisterExtensions(extMap)
	assert.NoError(t, err)

	buf, err := tv.ToCBOR()
	assert.NoError(t, err)

	other := NewUnsignedCorim()
	err = other.FromCBOR(buf)
	assert.NoError(t, err)

	assert.EqualValues(t, tv.Profile, other.Profile)
	assert.EqualValues(t, tv.ID, other.ID)
	assert.Nil(t, other.Entities)
}

func TestUnsignedCorim_extensions(t *testing.T) {
	c := NewUnsignedCorim()
	corimExts := struct{}{}
	extMap := extensions.NewMap().
		Add(ExtUnsignedCorim, &corimExts).
		Add(ExtEntity, &struct{}{})

	err := c.RegisterExtensions(extMap)
	assert.NoError(t, err)
	assert.Equal(t, &corimExts, c.GetExtensions())

	badMap := extensions.NewMap().Add(extensions.Point("test"), &struct{}{})
	err = c.RegisterExtensions(badMap)
	assert.EqualError(t, err, `unexpected extension point: "test"`)
}

func TestUnsignedCorim_truncated(t *testing.T) {
	c := NewUnsignedCorim()
	err := c.FromCBOR([]byte{0x01, 0x02})
	assert.EqualError(t, err, "input too short")
}

func TestLocator_Valid(t *testing.T) {
	l := Locator{}
	assert.EqualError(t, l.Valid(), "empty href")

	l.Href = comid.TaggedURI("https://example.com")
	assert.NoError(t, l.Valid())

	l.Thumbprint = &swid.HashEntry{}
	assert.EqualError(t, l.Valid(), "invalid locator thumbprint: unknown hash algorithm 0")

}

func TestTag_CBOR_marshal(t *testing.T) {
	tag := Tag{Number: 1337, Content: []byte{0x01, 0x02, 0x03}}

	actual, err := tag.MarshalCBOR()
	assert.NoError(t, err)

	expected := []byte{
		0xd9, 0x05, 0x39, // tag(1337)
		0x43, // bstr(3)
		0x01, 0x02, 0x03,
	}

	assert.Equal(t, expected, actual)
}

func TestUnsignedCorim_ToJSON_FromJSON_roundtrip(t *testing.T) {
	tv := NewUnsignedCorim().
		SetID("test-corim-id")
	require.NotNil(t, tv)

	// Use testGoodUnsignedCorimCBOR which is known to be valid
	var valid UnsignedCorim
	err := valid.FromCBOR(testGoodUnsignedCorimCBOR)
	require.NoError(t, err)

	jsonData, err := valid.ToJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded UnsignedCorim
	err = decoded.FromJSON(jsonData)
	assert.NoError(t, err)

	assert.EqualValues(t, valid.ID, decoded.ID)
	assert.EqualValues(t, valid.Profile, decoded.Profile)
	assert.EqualValues(t, valid.DependentRims, decoded.DependentRims)
	assert.EqualValues(t, valid.RimValidity, decoded.RimValidity)
	assert.EqualValues(t, valid.Entities, decoded.Entities)
	assert.Equal(t, len(valid.Tags), len(decoded.Tags))
	for i, expectedTag := range valid.Tags {
		actualTag := decoded.Tags[i]
		assert.Equal(t, expectedTag.Number, actualTag.Number)
		assertCBOREq(t, expectedTag.Content, actualTag.Content)
	}
}

func TestUnsignedCorim_FromJSON_Invalid(t *testing.T) {
	var c UnsignedCorim
	err := c.FromJSON([]byte(`{invalid}`))
	assert.Error(t, err)
}

func TestComid_iterators(t *testing.T) {
	cm := comid.NewTestComid(t)
	c := NewUnsignedCorim()
	c.AddComid(cm)

	keySeq, errFunc := c.IterAttestVerifKeys()
	assert.Equal(t,
		slices.Collect(keySeq)[0].VerifKeys[0].String(),
		(*cm.Triples.AttestVerifKeys)[0].VerifKeys[0].String(),
	)
	assert.NoError(t, errFunc())

	keySeq, errFunc = c.IterDevIdentityKeys()
	assert.Equal(t,
		slices.Collect(keySeq)[0].VerifKeys[0].String(),
		(*cm.Triples.DevIdentityKeys)[0].VerifKeys[0].String(),
	)
	assert.NoError(t, errFunc())

	valSeq, errFunc := c.IterRefVals()
	assert.Equal(t,
		slices.Collect(valSeq)[0].Measurements.Values[0].Val.RawValue,
		cm.Triples.ReferenceValues.Values[0].Measurements.Values[0].Val.RawValue,
	)
	assert.NoError(t, errFunc())

	valSeq, errFunc = c.IterEndVals()
	assert.Equal(t,
		slices.Collect(valSeq)[0].Measurements.Values[0].Val.RawValue,
		cm.Triples.EndorsedValues.Values[0].Measurements.Values[0].Val.RawValue,
	)
	assert.NoError(t, errFunc())
}
