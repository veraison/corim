// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
)

func TestCoSWIDTriple_NewCoSWIDTriple(t *testing.T) {
	s := NewCoSWIDTriple()
	require.NotNil(t, s)
}

func TestCoSWIDTriple_AddEnvironment(t *testing.T) {
	s := &CoSWIDTriple{}
	tv := &comid.Environment{
		Class: comid.NewClassUUID(TestUUID),
	}

	err := s.AddEnvironment(tv)
	require.Nil(t, err)
}

func TestCoSWIDTriple_AddEnvironment_NOK(t *testing.T) {
	expectedErr := "no environment to add"
	s := &CoSWIDTriple{}
	var tv *comid.Environment
	err := s.AddEnvironment(tv)
	assert.EqualError(t, err, expectedErr)
	expectedErr = "environment is not valid: environment must not be empty"
	tv = &comid.Environment{}
	err = s.AddEnvironment(tv)
	assert.EqualError(t, err, expectedErr)
}

func TestCoSWIDTriple_AddEvidence(t *testing.T) {
	s := &CoSWIDTriple{}
	tv := &CoSWIDEvidenceMap{
		TagID:    swid.NewTagID(TestTag),
		Evidence: swid.Evidence{Date: TestDate, DeviceID: TestDeviceID},
	}
	err := s.AddEvidence(tv)
	require.Nil(t, err)
}

func TestCoSWIDTriple_AddEvidence_NOK(t *testing.T) {
	expectedErr := "no evidence map to add"
	s := &CoSWIDTriple{}
	var tv *CoSWIDEvidenceMap
	err := s.AddEvidence(tv)
	assert.EqualError(t, err, expectedErr)
}
