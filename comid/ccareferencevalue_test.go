// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCCARefVal_SetLabel_ok(t *testing.T) {
	ccrefval := &CCARefValID{}
	err := ccrefval.SetLabel(TestCCALabel)
	assert.NoError(t, err)
}

func TestCCARefVal_SetLabel_nok(t *testing.T) {
	ccrefval := &CCARefValID{}
	err := ccrefval.SetLabel("")
	assert.EqualError(t, err, "no label supplied")
}
