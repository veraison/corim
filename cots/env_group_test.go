// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cots

import (
	"testing"

	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentGroup_JSON_full_Roundtrip_ok(t *testing.T) {
	expected := `{"swidtag":{"entity":[{"entity-name":"Round Tripper, Inc.","role":"softwareCreator"}]}}`
	expectedSwidTagEntityName := "Round Tripper, Inc."
	expectedSwidTagEntityRole := "softwareCreator"

	tv := EnvironmentGroup{}

	// Add tag
	swidTag := AbbreviatedSwidTag{}
	swidTag.Entities = swid.Entities{
		swid.Entity{EntityName: "Round Tripper, Inc."},
	}
	err := swidTag.Entities[0].SetRoles(swid.RoleSoftwareCreator)
	assert.Nil(t, err)
	tv.SwidTag = &swidTag

	// Roundtrip
	jTv, err := tv.ToJSON()
	t.Logf("JSON: '%s'", string(jTv))
	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(jTv))

	afterRoundTrip := EnvironmentGroup{}
	err = afterRoundTrip.FromJSON(jTv)
	assert.Nil(t, err)

	assert.EqualValues(t, expectedSwidTagEntityName, afterRoundTrip.SwidTag.Entities[0].EntityName)
	assert.EqualValues(t, expectedSwidTagEntityRole, afterRoundTrip.SwidTag.Entities[0].Roles.String())
}

func TestEnvironmentGroup_CBOR_Roundtrip_ok(t *testing.T) {
	expected := []byte{0xa1, 0x1, 0xa1, 0x0, 0xa1, 0x0, 0xd8, 0x6f, 0x44, 0x0, 0x0, 0x0, 0x0}
	expectedEnvClassOIDValue := comid.NewClassOID("0.0.0.0.0")

	tv := EnvironmentGroup{}

	// Add environment
	env := &comid.Environment{}
	env.Class = comid.NewClassOID("0.0.0.0.0")
	tv.Environment = env

	// Roundtrip
	cTv, err := tv.ToCBOR()
	t.Logf("CBOR: %x", cTv)
	assert.Nil(t, err)
	assert.Equal(t, expected, cTv)

	afterRoundTrip := EnvironmentGroup{}
	err = afterRoundTrip.FromCBOR(cTv)
	assert.Nil(t, err)

	assert.EqualValues(t, expectedEnvClassOIDValue, afterRoundTrip.Environment.Class)
}

func TestEnvironmentGroup_Valid_bad_abbreviated_swidtag(t *testing.T) {
	tv := EnvironmentGroup{}

	tv.SwidTag = &AbbreviatedSwidTag{}
	err := tv.Valid()

	assert.EqualError(t, err, "abbreviated swid tag validation failed: no entities present, must have at least 1 entity")
}

func TestEnvironmentGroup_Valid_empty_environment(t *testing.T) {
	tv := EnvironmentGroup{}

	tv.Environment = &comid.Environment{}
	err := tv.Valid()

	assert.EqualError(t, err, "environment group validation failed: environment must not be empty")
}

func TestEnvironmentGroups_Valid_bad_environmentgroup_in_list(t *testing.T) {
	// In this case using bad environment entry
	envG := EnvironmentGroup{}
	envG.Environment = &comid.Environment{}

	tv := EnvironmentGroups{envG}
	err := tv.Valid()

	assert.EqualError(t, err, "bad environment group at index 0: environment group validation failed: environment must not be empty")
}

func TestEnvironmentGroups_JSON_Roundtrip_ok(t *testing.T) {
	expected := `[
		{
			"environment": {
				"class": {
					"vendor":"Zesty Hands, Inc."
				}
			}
		},
		{
			"environment": {
				"class": {
					"vendor":"Snobbish Apparel, Inc."
				}
			}
		}		
	]`
	tvs := EnvironmentGroups{}

	egVendorModel := EnvironmentGroup{}
	egVendorModel.Environment = &comid.Environment{}
	egVendorModel.Environment.Class = &comid.Class{}
	egVendorModel.Environment.Class.SetVendor("Zesty Hands, Inc.")
	tvs.AddEnvironmentGroup(egVendorModel)

	egVendorModel1 := EnvironmentGroup{}
	egVendorModel1.Environment = &comid.Environment{}
	egVendorModel1.Environment.Class = &comid.Class{}
	egVendorModel1.Environment.Class.SetVendor("Snobbish Apparel, Inc.")
	tvs.AddEnvironmentGroup(egVendorModel1)

	VendorModelsJSON, err := tvs.ToJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(VendorModelsJSON))

	t.Logf("JSON: '%s'", string(VendorModelsJSON))

	actual := EnvironmentGroups{}
	err = actual.FromJSON(VendorModelsJSON)
	assert.Nil(t, err)
	assert.Equal(t, tvs, actual)
}
