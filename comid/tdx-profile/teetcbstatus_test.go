// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTcbStatus_NewTeeTcbStatusString_OK(t *testing.T) {
	s := TestTCBStatus
	_, err := NewTeeTcbStatusString(s)
	require.NoError(t, err)
}

func TestTcbStatus_NewTeeTcbStatusExpr_OK(t *testing.T) {
	s := TestTCBStatus
	_, err := NewTcbStatusExpr(MEM, s)
	require.NoError(t, err)
}

func TestTcbStatus_NewTeeTcbStatusString_NOK(t *testing.T) {
	expectedErr := "zero len value for TeeTcbStatus"
	s := []string{}
	_, err := NewTeeTcbStatusString(s)
	assert.EqualError(t, err, expectedErr)
}

func TestTcbStatus_NewTeeTcbStatusExpr_NOK(t *testing.T) {
	expectedErr := "invalid operator : 5"
	s := TestTCBStatus
	_, err := NewTcbStatusExpr(NOP, s)
	assert.EqualError(t, err, expectedErr)
}

func TestTcbStatus_AddTcbStatus_OK(t *testing.T) {
	s := TestTCBStatus
	status := TeeTcbStatus{val: []string{"abcd"}}
	_, err := status.AddTeeTcbStatus(NOP, s)
	require.Nil(t, err)
}

func TestTcbStatus_AddTcbStatus_NOK(t *testing.T) {
	expectedErr := "operator mis-match TeeTcbStatus Op: 6, Input Op: 2"
	s := TestTCBStatus
	status, err := NewTcbStatusExpr(MEM, s)
	require.Nil(t, err)
	_, err = status.AddTeeTcbStatus(GE, []string{"abcd"})
	assert.EqualError(t, err, expectedErr)
}

func TestTcbStatus_Valid_OK(t *testing.T) {
	s := TestTCBStatus
	status, err := NewTeeTcbStatusString(s)
	require.Nil(t, err)
	err = status.Valid()
	require.Nil(t, err)
}

func TestTcbStatus_Valid_NOK(t *testing.T) {
	expectedErr := "invalid operator in a TeeTcbStatus: 2"
	status := TeeTcbStatus{val: TaggedSetStringExpression(SetStringExpression{SetOperator: 2, SetString: []string{"valid"}})}
	err := status.Valid()
	assert.EqualError(t, err, expectedErr)
	expectedErr = "unknown type []int for TeeTcbStatus"
	status = TeeTcbStatus{val: TestInvalidTCBStatus}
	err = status.Valid()
	assert.EqualError(t, err, expectedErr)
}
