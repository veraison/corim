// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cots

import (
	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConciseTaStore(t *testing.T) {
	// TODO add expected values or remove test cases
	shared_ta, _ := ioutil.ReadFile("../cocli/data/cots/shared_ta.der")
	shared_ca, _ := ioutil.ReadFile("../cocli/data/cots/shared_ca.der")
	snob_ta, _ := ioutil.ReadFile("../cocli/data/cots/Snobbish Apparel_ta.der")
	zesty_ta, _ := ioutil.ReadFile("../cocli/data/cots/Zesty Hands_ta.der")

	eg_dice := NewEnvironmentGroup()
	eg_dice.SetNamedTaStore("DICE Trust Anchors")

	cots_dice := ConciseTaStore{}
	cots_dice.Keys = NewTasAndCas()
	cots_dice.Keys.AddTaCert(shared_ta)
	cots_dice.Keys.AddTaCert(snob_ta)
	cots_dice.Keys.AddTaCert(zesty_ta)
	cots_dice.Environments = *NewEnvironmentGroups()
	cots_dice.Environments.AddEnvironmentGroup(*eg_dice)

	cots_dice_cbor, _ := cots_dice.ToCBOR()
	assert.NotNil(t, cots_dice_cbor)
	cots_dice_json, _ := cots_dice.ToJSON()
	assert.NotNil(t, cots_dice_json)

	eg_shared := NewEnvironmentGroup()
	eg_shared.Environment = &comid.Environment{}
	eg_shared.Environment.Class = comid.NewClassOID("1.2.3.4.5")

	cots_shared := ConciseTaStore{}
	cots_shared.Keys = NewTasAndCas()
	cots_shared.Keys.AddTaCert(shared_ta)
	cots_shared.Keys.AddCaCert(shared_ca)
	cots_shared.Environments = make([]EnvironmentGroup, 1)
	cots_shared.Environments[0] = *eg_shared

	cots_shared_cbor, _ := cots_shared.ToCBOR()
	assert.NotNil(t, cots_shared_cbor)

	eg_zesty := NewEnvironmentGroup()
	eg_zesty.SwidTag = &AbbreviatedSwidTag{}
	eg_zesty.SwidTag.Entities = swid.Entities{}
	e_zesty := swid.Entity{EntityName: "Zesty Hands, Inc."}
	e_zesty.SetRoles(swid.RoleSoftwareCreator)
	eg_zesty.SwidTag.Entities = append(eg_zesty.SwidTag.Entities, e_zesty)

	cots_zesty := ConciseTaStore{}
	cots_zesty.Keys = NewTasAndCas()
	cots_zesty.Keys.AddTaCert(zesty_ta)
	cots_zesty.Environments = make([]EnvironmentGroup, 1)
	cots_zesty.Environments[0] = *eg_zesty

	cots_zesty_cbor, _ := cots_zesty.ToCBOR()
	assert.NotNil(t, cots_zesty_cbor)

	cots_zesty_perm_claim := ConciseTaStore{}
	cots_zesty_perm_claim.Keys = NewTasAndCas()
	cots_zesty_perm_claim.Keys.AddTaCert(zesty_ta)
	cots_zesty_perm_claim.Environments = make([]EnvironmentGroup, 1)
	cots_zesty_perm_claim.Environments[0] = *eg_zesty


	perm_name := "Bitter Paper"
	perm_claims1 := EatCWTClaim{SoftwareNameLabel: &perm_name}
	cots_zesty_perm_claim.PermClaims = append(cots_zesty_perm_claim.PermClaims, perm_claims1)

	cots_zesty_perm_claim_cbor, _ := cots_zesty_perm_claim.ToCBOR()
	assert.NotNil(t, cots_zesty_perm_claim_cbor)

	eg_snob := NewEnvironmentGroup()
	eg_snob.SwidTag = &AbbreviatedSwidTag{}
	eg_snob.SwidTag.Entities = swid.Entities{}
	e_snob := swid.Entity{EntityName: "Snobbish Apparel, Inc."}
	e_snob.SetRoles(swid.RoleSoftwareCreator)
	eg_snob.SwidTag.Entities = append(eg_snob.SwidTag.Entities, e_snob)

	cots_snob_excl_claim := ConciseTaStore{}
	cots_snob_excl_claim.Keys = NewTasAndCas()
	cots_snob_excl_claim.Keys.AddTaCert(snob_ta)
	cots_snob_excl_claim.Environments = make([]EnvironmentGroup, 1)
	cots_snob_excl_claim.Environments[0] = *eg_snob

	excl_name := "Legal Lawyer"
	excl_claims1 := EatCWTClaim{SoftwareNameLabel: &excl_name}
	cots_snob_excl_claim.ExclClaims = append(cots_snob_excl_claim.ExclClaims, excl_claims1)

	cots_snob_excl_claim_cbor, _ := cots_snob_excl_claim.ToCBOR()
	assert.NotNil(t, cots_snob_excl_claim_cbor)

	cts := NewConciseTaStores()
	cts.AddConciseTaStores(cots_dice)
	cts.AddConciseTaStores(cots_shared)
	cts.AddConciseTaStores(cots_zesty)
	cts.AddConciseTaStores(cots_zesty_perm_claim)
	cts.AddConciseTaStores(cots_snob_excl_claim)

	cts_cbor2, _ := cts.ToCBOR()
	assert.NotNil(t, cts_cbor2)

	cts_cbor3 := append(CotsTag, cts_cbor2...)
	assert.NotNil(t, cts_cbor3)

	var roundtrip ConciseTaStores
	errrt := roundtrip.FromCBOR(cts_cbor3)
	assert.Nil(t, errrt)
	assert.Nil(t, roundtrip.Valid())
}
