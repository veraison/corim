// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvidenceID_SetUUID_OK(t *testing.T) {
	ev := &EvidenceID{}
	testUUID, err := uuid.Parse(TestUUIDString)
	require.NoError(t, err)
	i := ev.SetUUID(testUUID)
	require.NotNil(t, i)
}

func TestEvidenceID_GetUUID_OK(t *testing.T) {
	ev := MustNewUUIDEvidenceID(TestUUID)
	require.NotNil(t, ev)
	u, err := ev.GetUUID()
	assert.Nil(t, err)
	assert.Equal(t, u, TestUUID)
}

func TestEvidence_GetUUID_NOK(t *testing.T) {
	ev := &EvidenceID{}
	expectedErr := "evidence-id type is: <nil>"
	_, err := ev.GetUUID()
	assert.EqualError(t, err, expectedErr)
}
