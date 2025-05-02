// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeeTcbCompSvn_NewTeeTcbCompSvnUint_OK(t *testing.T) {
	_, err := NewTeeTcbCompSvnUint(TestCompSvn)
	require.NoError(t, err)
}

func TestTeeTcbCompSvn_NewTeeTcbCompSvnUint_NOK(t *testing.T) {
	expectedErr := "no value supplied for TeeTcbCompSVN"
	var val []uint
	_, err := NewTeeTcbCompSvnUint(val)
	assert.EqualError(t, err, expectedErr)
	expectedErr = "invalid length 18 for TeeTcbCompSVN"
	v := make([]uint, 18)
	_, err = NewTeeTcbCompSvnUint(v)
	assert.EqualError(t, err, expectedErr)
}

func TestTeeTcbCompSvn_NewTeeTcbCompSvnExpression_OK(t *testing.T) {
	_, err := NewTeeTcbCompSvnExpression(TestCompSvn)
	require.NoError(t, err)
}

func TestTeeTcbCompSvn_NewTeeTcbCompSvnExpression_NOK(t *testing.T) {
	expectedErr := "no value supplied for TeeTcbCompSVN"
	var val []uint
	_, err := NewTeeTcbCompSvnExpression(val)
	assert.EqualError(t, err, expectedErr)
	expectedErr = "invalid length 18 for TeeTcbCompSVN"
	v := make([]uint, 18)
	_, err = NewTeeTcbCompSvnExpression(v)
	assert.EqualError(t, err, expectedErr)
}

func TestTeeTcbCompSvn_Valid_OK(t *testing.T) {
	tc, err := NewTeeTcbCompSvnExpression(TestCompSvn)
	require.NoError(t, err)
	err = tc.Valid()
	require.NoError(t, err)
	tc, err = NewTeeTcbCompSvnUint(TestCompSvn)
	require.NoError(t, err)
	err = tc.Valid()
	require.NoError(t, err)
}

func TestTeeTcbCompSvn_Valid_NOK(t *testing.T) {
	expectedErr := "invalid TeeSvn at index 0, unknown type <nil> for TeeSVN"
	var tc TeeTcbCompSvn
	err := tc.Valid()
	assert.EqualError(t, err, expectedErr)
}
