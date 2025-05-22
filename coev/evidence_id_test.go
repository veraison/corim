// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvidenceID_NewEvidenceID(t *testing.T) {
	ev, err := NewEvidenceID(TestUUIDString, "uuid")
	require.NoError(t, err)
	require.NotNil(t, ev)
}

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

func TestEvidenceID_UnmarshalJSON_OK(t *testing.T) {
	for _, tv := range []struct {
		Name  string
		Input string
	}{
		{
			Name:  "valid input test 1",
			Input: `{ "type": "uuid", "value": "31fb5abf-023e-4992-aa4e-95f9c1503bfa" }`,
		},
		{
			Name:  "valid input test 2",
			Input: `{ "type": "uuid", "value": "31fb5abf-023e-4992-aa4e-95f9c1503bfb"}`,
		},
	} {
		t.Run(tv.Name, func(t *testing.T) {
			var actual EvidenceID
			err := actual.UnmarshalJSON([]byte(tv.Input))
			require.NoError(t, err)
		})
	}
}
