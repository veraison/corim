// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestInt = 10
const TestUint = uint(10)
const TestFloat = 10.123

func TestNumericType_NewNumericType(t *testing.T) {
	_, err := NewNumericType(uint(10))
	require.NoError(t, err)
	_, err = NewNumericType(10)
	require.NoError(t, err)
	_, err = NewNumericType(10.123)
	require.NoError(t, err)
}

func TestNumericType_SetNumericType(t *testing.T) {
	o := &NumericType{}
	err := o.SetNumericType(uint(10))
	require.NoError(t, err)
	err = o.SetNumericType(10)
	require.NoError(t, err)
	err = o.SetNumericType(10.123)
	require.NoError(t, err)
}

func TestNumericType_IsType(t *testing.T) {
	o := &NumericType{val: 10}
	b := o.IsInt()
	require.True(t, b)
	o = &NumericType{val: uint(10)}
	b = o.IsUint()
	require.True(t, b)
	o = &NumericType{val: 10.00}
	require.True(t, b)
}

func TestNumericType_GetType(t *testing.T) {
	o := &NumericType{val: 10}
	i, err := o.GetInt()
	require.NoError(t, err)
	assert.Equal(t, TestInt, i)

	o = &NumericType{val: TestUint}
	k, err := o.GetUint()
	require.NoError(t, err)
	assert.Equal(t, TestUint, k)

	o = &NumericType{val: TestFloat}
	f, err := o.GetFloat()
	require.NoError(t, err)
	assert.Equal(t, TestFloat, f)
}

func TestNumericType_GetType_NOK(t *testing.T) {
	o := &NumericType{val: 10}
	i, err := o.GetInt()
	require.NoError(t, err)
	assert.Equal(t, TestInt, i)

	o = &NumericType{val: TestUint}
	k, err := o.GetUint()
	require.NoError(t, err)
	assert.Equal(t, TestUint, k)

	o = &NumericType{val: TestFloat}
	f, err := o.GetFloat()
	require.NoError(t, err)
	assert.Equal(t, TestFloat, f)
}

func TestNumericType_Valid_OK(t *testing.T) {
	o := NumericType{val: 10}
	err := o.Valid()
	require.NoError(t, err)
}

func TestNumericType_Valid_NOK(t *testing.T) {
	expectedErr := "unsupported NumericType type: string"
	o := NumericType{val: "test"}
	err := o.Valid()
	assert.EqualError(t, err, expectedErr)
}
