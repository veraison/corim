// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCCAPlatformConfigID_Set_ok(t *testing.T) {
	var cca CCAPlatformConfigID

	err := cca.Set(TestCCALabel)
	assert.NoError(t, err)
}

func TestCCAPlatformConfigID_Set_nok(t *testing.T) {
	var cca CCAPlatformConfigID
	expectedErr := "empty input string"
	err := cca.Set("")
	assert.EqualError(t, err, expectedErr)
}

func TestCCAPlatformConfigID_Get_nok(t *testing.T) {
	var cca CCAPlatformConfigID
	expectedErr := "empty CCA platform config ID"
	_, err := cca.Get()
	assert.EqualError(t, err, expectedErr)
}
