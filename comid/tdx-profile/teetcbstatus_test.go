// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initTcbStatus() []any {
	s := make([]any, len(TestTCBStatus))
	for i := range TestTCBStatus {
		s[i] = TestTCBStatus[i]
	}
	return s
}

func TestTcbStatus_NewTeeTcbStatus_OK(t *testing.T) {
	s := initTcbStatus()
	_, err := NewTeeTcbStatus(s)
	require.NoError(t, err)
}

func TestTcbStatus_NewTeeTcbStatus_NOK(t *testing.T) {
	s := make([]any, len(TestTCBStatus))
	for i := range TestTCBStatus {
		s[i] = i
	}
	expectedErr := "invalid type: int for tcb status at index: 0"
	_, err := NewTeeTcbStatus(s)
	assert.EqualError(t, err, expectedErr)
	var m []any
	expectedErr = "nil value argument"
	_, err = NewTeeTcbStatus(m)
	assert.EqualError(t, err, expectedErr)
}

func TestTcbStatus_AddTcbStatus_OK(t *testing.T) {
	s := initTcbStatus()
	status := TeeTcbStatus{}
	err := status.AddTeeTcbStatus(s)
	require.Nil(t, err)
}

func TestTcbStatus_AddTcbStatus_NOK(t *testing.T) {
	expectedErr := "invalid type: int for tcb status at index: 0"
	s := make([]any, len(TestInvalidTCBStatus))
	for i := range TestInvalidTCBStatus {
		s[i] = TestInvalidTCBStatus[i]
	}
	status := TeeTcbStatus{}
	err := status.AddTeeTcbStatus(s)
	assert.EqualError(t, err, expectedErr)
}

func TestTcbStatus_Valid_OK(t *testing.T) {
	s := initTcbStatus()
	status, err := NewTeeTcbStatus(s)
	require.Nil(t, err)
	err = status.Valid()
	require.Nil(t, err)
}

func TestTcbStatus_Valid_NOK(t *testing.T) {
	expectedErr := "empty tcb status"
	status := TeeTcbStatus{}
	err := status.Valid()
	assert.EqualError(t, err, expectedErr)
	expectedErr = "invalid type: int for tcb status at index: 0"
	s := make([]any, len(TestInvalidTCBStatus))
	for i := range TestInvalidTCBStatus {
		s[i] = TestInvalidTCBStatus[i]
	}
	status = TeeTcbStatus(s)
	err = status.Valid()
	assert.EqualError(t, err, expectedErr)
}
