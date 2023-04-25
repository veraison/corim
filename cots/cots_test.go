// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/corim/comid"
)

func TestConciseTaStore_Valid_no_environment_groups(t *testing.T) {
	cots := ConciseTaStore{}
	assert.EqualError(t, cots.Valid(), "environmentGroups must be present")
}

func TestConciseTaStore_Valid_invalid_environment_groups(t *testing.T) {
	cots := ConciseTaStore{}
	cots.Environments = EnvironmentGroups{
		EnvironmentGroup{
			Environment: &comid.Environment{},
		},
	}

	assert.EqualError(t, cots.Valid(), "invalid environmentGroups: bad environment group at index 0: environment group validation failed: environment must not be empty")

}

func TestConciseTaStore_Valid_empty_keys(t *testing.T) {
	cots := ConciseTaStore{}
	cots.Environments = EnvironmentGroups{}
	cots.Keys = &TasAndCas{}

	assert.EqualError(t, cots.Valid(), "empty Keys")
}

func TestConciseTaStore_Valid_invalid_tag_identity(t *testing.T) {
	cots := ConciseTaStore{}
	cots.Environments = EnvironmentGroups{}
	cots.Keys = NewTasAndCas().AddTaCert(ta)
	cots.TagIdentity = &comid.TagIdentity{}

	assert.EqualError(t, cots.Valid(), "invalid TagIdentity: empty tag-id")
}

func TestConciseTaStores_Valid_empty_stores(t *testing.T) {
	cotsList := ConciseTaStores{}
	assert.EqualError(t, cotsList.Valid(), "empty concise-ta-stores")
}

func TestConciseTaStores_Valid_bad_cots(t *testing.T) {
	cotsList := ConciseTaStores{*NewConciseTaStore().AddPurpose("cots")}
	assert.EqualError(t, cotsList.Valid(), "bad ConciseTaStore group at index 0: environmentGroups must be present")
}
