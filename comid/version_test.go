// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/swid"
)

func TestVersion_Valid_OK(t *testing.T) {
	v := NewVersion()

	v.SetVersion("1.55.22")
	v.SetScheme(swid.VersionSchemeSemVer)

	assert.Nil(t, v.Valid())
}

func TestVersion_Valid_NOK(t *testing.T) {
	v := NewVersion()

	assert.EqualError(t, v.Valid(), "empty version")
}

func TestVersion_Equal_True(t *testing.T) {
	claim := NewVersion()
	claim.SetVersion("1.55.22")
	claim.SetScheme(swid.VersionSchemeSemVer)

	ref := NewVersion()
	ref.SetVersion("1.55.22")
	ref.SetScheme(swid.VersionSchemeSemVer)

	assert.True(t, claim.Equal(*ref))
}

func TestVersion_Equal_False(t *testing.T) {
	claim := NewVersion()
	claim.SetVersion("1.55.22")
	claim.SetScheme(swid.VersionSchemeSemVer)

	ref := NewVersion()
	ref.SetVersion("1.55.40")
	ref.SetScheme(swid.VersionSchemeSemVer)

	assert.False(t, claim.Equal(*ref))
}
