// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cots

import (
	"os"
	"testing"

	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"

	"github.com/stretchr/testify/assert"
)

func TestNewConciseTaStore(t *testing.T) {
	// TODO add expected values or remove test cases
	sharedTa, _ := os.ReadFile("../cocli/data/cots/shared_ta.der")
	sharedCa, _ := os.ReadFile("../cocli/data/cots/shared_ca.der")
	snobTa, _ := os.ReadFile("../cocli/data/cots/Snobbish Apparel_ta.der")
	zestyTa, _ := os.ReadFile("../cocli/data/cots/Zesty Hands_ta.der")

	egDice := NewEnvironmentGroup()
	egDice.SetNamedTaStore("DICE Trust Anchors")

	cotsDice := ConciseTaStore{}
	cotsDice.Keys = NewTasAndCas()
	cotsDice.Keys.AddTaCert(sharedTa)
	cotsDice.Keys.AddTaCert(snobTa)
	cotsDice.Keys.AddTaCert(zestyTa)
	cotsDice.Environments = *NewEnvironmentGroups()
	cotsDice.Environments.AddEnvironmentGroup(*egDice)

	cotsDiceCbor, _ := cotsDice.ToCBOR()
	assert.NotNil(t, cotsDiceCbor)
	cotsDiceJSON, _ := cotsDice.ToJSON()
	assert.NotNil(t, cotsDiceJSON)

	egShared := NewEnvironmentGroup()
	egShared.Environment = &comid.Environment{}
	egShared.Environment.Class = comid.NewClassOID("1.2.3.4.5")

	cotsShared := ConciseTaStore{}
	cotsShared.Keys = NewTasAndCas()
	cotsShared.Keys.AddTaCert(sharedTa)
	cotsShared.Keys.AddCaCert(sharedCa)
	cotsShared.Environments = make([]EnvironmentGroup, 1)
	cotsShared.Environments[0] = *egShared

	cotsSharedCbor, _ := cotsShared.ToCBOR()
	assert.NotNil(t, cotsSharedCbor)

	egZesty := NewEnvironmentGroup()
	egZesty.SwidTag = &AbbreviatedSwidTag{}
	egZesty.SwidTag.Entities = swid.Entities{}
	eZesty := swid.Entity{EntityName: "Zesty Hands, Inc."}
	err := eZesty.SetRoles(swid.RoleSoftwareCreator)
	assert.Nil(t, err)
	egZesty.SwidTag.Entities = append(egZesty.SwidTag.Entities, eZesty)

	cotsZesty := ConciseTaStore{}
	cotsZesty.Keys = NewTasAndCas()
	cotsZesty.Keys.AddTaCert(zestyTa)
	cotsZesty.Environments = make([]EnvironmentGroup, 1)
	cotsZesty.Environments[0] = *egZesty

	cotsZestyCbor, _ := cotsZesty.ToCBOR()
	assert.NotNil(t, cotsZestyCbor)

	cotsZestyPermClaim := ConciseTaStore{}
	cotsZestyPermClaim.Keys = NewTasAndCas()
	cotsZestyPermClaim.Keys.AddTaCert(zestyTa)
	cotsZestyPermClaim.Environments = make([]EnvironmentGroup, 1)
	cotsZestyPermClaim.Environments[0] = *egZesty

	permName := "Bitter Paper"
	permClaims1 := EatCWTClaim{SoftwareNameLabel: &permName}
	cotsZestyPermClaim.PermClaims = append(cotsZestyPermClaim.PermClaims, permClaims1)

	cotsZestyPermClaimCbor, _ := cotsZestyPermClaim.ToCBOR()
	assert.NotNil(t, cotsZestyPermClaimCbor)

	egSnob := NewEnvironmentGroup()
	egSnob.SwidTag = &AbbreviatedSwidTag{}
	egSnob.SwidTag.Entities = swid.Entities{}
	eSnob := swid.Entity{EntityName: "Snobbish Apparel, Inc."}
	err = eSnob.SetRoles(swid.RoleSoftwareCreator)
	assert.Nil(t, err)
	egSnob.SwidTag.Entities = append(egSnob.SwidTag.Entities, eSnob)

	cotsSnobExclClaim := ConciseTaStore{}
	cotsSnobExclClaim.Keys = NewTasAndCas()
	cotsSnobExclClaim.Keys.AddTaCert(snobTa)
	cotsSnobExclClaim.Environments = make([]EnvironmentGroup, 1)
	cotsSnobExclClaim.Environments[0] = *egSnob

	exclName := "Legal Lawyer"
	exclClaims1 := EatCWTClaim{SoftwareNameLabel: &exclName}
	cotsSnobExclClaim.ExclClaims = append(cotsSnobExclClaim.ExclClaims, exclClaims1)

	cotsSnobExclClaimCbor, _ := cotsSnobExclClaim.ToCBOR()
	assert.NotNil(t, cotsSnobExclClaimCbor)

	cts := NewConciseTaStores()
	cts.AddConciseTaStores(cotsDice)
	cts.AddConciseTaStores(cotsShared)
	cts.AddConciseTaStores(cotsZesty)
	cts.AddConciseTaStores(cotsZestyPermClaim)
	cts.AddConciseTaStores(cotsSnobExclClaim)

	ctsCbor2, _ := cts.ToCBOR()
	assert.NotNil(t, ctsCbor2)

	ctsCbor3 := append(CotsTag, ctsCbor2...)
	assert.NotNil(t, ctsCbor3)

	var roundtrip ConciseTaStores
	errrt := roundtrip.FromCBOR(ctsCbor3)
	assert.Nil(t, errrt)
	assert.Nil(t, roundtrip.Valid())
}
