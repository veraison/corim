// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTeeAttributes_NewTeeAttributes_OK(t *testing.T) {
	_, err := NewTeeAttributes(TestTeeAttributes)
	require.Nil(t, err)
}

func TestNewTeeAttributes_NewTeeAttributes_NOK(t *testing.T) {
	expectedErr := "nil TeeAttributes"
	_, err := NewTeeAttributes(nil)
	assert.EqualError(t, err, expectedErr)
}

func TestNewTeeAttributes_Valid_OK(t *testing.T) {
	tA := TeeAttributes(TestTeeAttributes)
	err := tA.Valid()
	require.Nil(t, err)
}

func TestNewTeeAttributes_Valid_NOK(t *testing.T) {
	tA := TeeAttributes{}
	expectedErr := "zero len TeeAttributes"
	err := tA.Valid()
	assert.EqualError(t, err, expectedErr)
}
