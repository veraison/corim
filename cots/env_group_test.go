// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cots

import (
	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvironmentGroup(t *testing.T) {
	tv := EnvironmentGroup{}
	tv.SetNamedTaStore("Some TA Store")
	assert.NotNil(t, tv)
	tvs := EnvironmentGroups{}
	tvs.AddEnvironmentGroup(tv)
	jTvs, _ := tvs.ToJSON()
	assert.NotNil(t, jTvs)

	egShared := EnvironmentGroup{}
	egShared.Environment = &comid.Environment{}
	egShared.Environment.Class = comid.NewClassOID("1.2.3.4.5")
	tvs1 := EnvironmentGroups{}
	tvs1.AddEnvironmentGroup(egShared)
	jEgShared, _ := tvs1.ToJSON()
	assert.NotNil(t, jEgShared)

	egZesty := NewEnvironmentGroup()
	egZesty.SwidTag = &AbbreviatedSwidTag{}
	egZesty.SwidTag.Entities = swid.Entities{}
	e_zesty := swid.Entity{EntityName: "Zesty Hands, Inc."}
	e_zesty.SetRoles(swid.RoleSoftwareCreator)
	egZesty.SwidTag.Entities = append(egZesty.SwidTag.Entities, e_zesty)
	tvs2 := EnvironmentGroups{}
	tvs2.AddEnvironmentGroup(*egZesty)
	jEgZesty, _ := tvs2.ToJSON()
	assert.NotNil(t, jEgZesty)

	egVendorModel := EnvironmentGroup{}
	egVendorModel.Environment = &comid.Environment{}
	egVendorModel.Environment.Class = &comid.Class{}
	egVendorModel.Environment.Class.SetVendor("Zesty Hands, Inc.")
	tvs3 := EnvironmentGroups{}
	tvs3.AddEnvironmentGroup(egVendorModel)
	jEgVendorModel, _ := tvs3.ToJSON()
	assert.NotNil(t, jEgVendorModel)

	egVendorModel1 := EnvironmentGroup{}
	egVendorModel1.Environment = &comid.Environment{}
	egVendorModel1.Environment.Class = &comid.Class{}
	egVendorModel1.Environment.Class.SetVendor("Snobbish Apparel, Inc.")
	tvs4 := EnvironmentGroups{}
	tvs4.AddEnvironmentGroup(egVendorModel)
	tvs4.AddEnvironmentGroup(egVendorModel1)
	jEgVendorsModel, _ := tvs4.ToJSON()
	assert.NotNil(t, jEgVendorsModel)
}
