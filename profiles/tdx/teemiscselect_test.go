// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeeMiscSelect_NewTeeMiscSelect_OK(t *testing.T) {
	_, err := NewTeeMiscSelect(TestTeeMiscSelect)
	require.Nil(t, err)
}

func TestTeeMiscSelect_NewTeeMiscSelect_NOK(t *testing.T) {
	expectedErr := "nil value for TeeMiscSelect"
	_, err := NewTeeMiscSelect(nil)
	assert.EqualError(t, err, expectedErr)
}

func TestNewTeeMiscSelect_Valid_OK(t *testing.T) {
	tA := TeeMiscSelect(TestTeeMiscSelect)
	err := tA.Valid()
	require.Nil(t, err)
}

func TestTeeMiscSelect_Valid_NOK(t *testing.T) {
	tA := TeeMiscSelect{}
	expectedErr := "zero len TeeMiscSelect"
	err := tA.Valid()
	assert.EqualError(t, err, expectedErr)
}
