// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initAdvisoryIDs() []any {
	s := make([]any, len(TestAdvisoryIDs))
	for i := range TestAdvisoryIDs {
		s[i] = TestAdvisoryIDs[i]
	}
	return s
}

func TestAdvisoryIDs_NewTeeAvisoryIDs_OK(t *testing.T) {
	a := initAdvisoryIDs()
	_, err := NewTeeAvisoryIDs(a)
	require.Nil(t, err)
}

func TestAdvisoryIDs_NewTeeAvisoryIDs_NOK(t *testing.T) {
	expectedErr := "invalid type: int for AdvisoryIDs at index: 0"
	a := make([]any, len(TestAdvisoryIDs))
	for i := range TestAdvisoryIDs {
		a[i] = i
	}
	_, err := NewTeeAvisoryIDs(a)
	assert.EqualError(t, err, expectedErr)
}

func TestAdvisoryIDs_AddAdvisoryIDs_OK(t *testing.T) {
	a := initAdvisoryIDs()
	adv := TeeAdvisoryIDs{}
	err := adv.AddTeeAdvisoryIDs(a)
	require.NoError(t, err)
}

func TestAdvisoryIDs_AddAdvisoryIDs_NOK(t *testing.T) {
	expectedErr := "invalid type: float64 for AdvisoryIDs at index: 0"
	s := make([]any, len(TestInvalidAdvisoryIDs))
	for i := range TestInvalidAdvisoryIDs {
		s[i] = TestInvalidAdvisoryIDs[i]
	}
	adv := TeeAdvisoryIDs{}
	err := adv.AddTeeAdvisoryIDs(s)
	assert.EqualError(t, err, expectedErr)
}

func TestAdvisoryIDs_Valid_OK(t *testing.T) {
	a := initAdvisoryIDs()
	adv, err := NewTeeAvisoryIDs(a)
	require.NoError(t, err)
	err = adv.Valid()
	require.NoError(t, err)
}

func TestAdvisoryIDs_Valid_NOK(t *testing.T) {
	expectedErr := "empty AdvisoryIDs"
	adv := TeeAdvisoryIDs{}
	err := adv.Valid()
	assert.EqualError(t, err, expectedErr)

	expectedErr = "invalid type: float64 for AdvisoryIDs at index: 0"
	s := make([]any, len(TestInvalidAdvisoryIDs))
	for i := range TestInvalidAdvisoryIDs {
		s[i] = TestInvalidAdvisoryIDs[i]
	}
	adv = TeeAdvisoryIDs(s)
	err = adv.Valid()
	assert.EqualError(t, err, expectedErr)

}
