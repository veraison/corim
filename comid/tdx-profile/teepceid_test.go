// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPCEID_NewTeePCEID_OK(t *testing.T) {
	_, err := NewTeePCEID(TestPCEID)
	require.NoError(t, err)
}

func TestPCEID_NewTeePCEID_NOK(t *testing.T) {
	expectedErr := "null string for TeePCEID"
	_, err := NewTeePCEID("")
	assert.EqualError(t, err, expectedErr)
}

func TestPCEID_Valid_OK(t *testing.T) {
	pceID, err := NewTeePCEID(TestPCEID)
	require.NoError(t, err)
	err = pceID.Valid()
	require.NoError(t, err)
}

func TestPCEID_Valid_NOK(t *testing.T) {
	pceID := TeePCEID("")
	expectedErr := "nil TeePCEID"
	err := pceID.Valid()
	assert.EqualError(t, err, expectedErr)
}
