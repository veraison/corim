// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeeTcbEvalNumber_NewTeeTcbEvalNumberNumeric_OK(t *testing.T) {
	_, err := NewTeeTcbEvalNumberNumeric(TestTCBEvalNum)
	require.NoError(t, err)
}

func TestTeeTcbEvalNumber_NewTeeTcbEvalNumberUint_OK(t *testing.T) {
	_, err := NewTeeTcbEvalNumberUint(TestTCBEvalNum)
	require.NoError(t, err)
}

func TestTeeTcbEvalNumber_Valid_OK(t *testing.T) {
	tcb := TeeTcbEvalNumber{val: TestTCBEvalNum}
	err := tcb.Valid()
	require.NoError(t, err)
}

func TestTeeTcbEvalNumber_Valid_NOK(t *testing.T) {
	expectedErr := "unknown type int for TeeTcbEvalNumber"
	en := TeeTcbEvalNumber{val: -10}
	err := en.Valid()
	assert.EqualError(t, err, expectedErr)
}

func TestTeeTcbEvalNumber_GetUint_OK(t *testing.T) {
	tcb := TeeTcbEvalNumber{val: TestTCBEvalNum}
	val, err := tcb.GetUint()
	require.NoError(t, err)
	require.Equal(t, val, TestTCBEvalNum)
}

func TestTeeTcbEvalNumber_IsUint_OK(t *testing.T) {
	tcb := TeeTcbEvalNumber{val: TestTCBEvalNum}
	b := tcb.IsUint()
	require.True(t, b)
}

func TestTeeTcbEvalNumber_JSON1(t *testing.T) {
	expectedErr := "unknown type int for TeeTcbEvalNumber"
	en := TeeTcbEvalNumber{val: -10}
	err := en.Valid()
	assert.EqualError(t, err, expectedErr)
}

func TestTeeTcbEvalNumber_JSON(t *testing.T) {

	for _, tv := range []struct {
		input         interface{}
		ExpectedBytes []byte
	}{
		{
			input:         uint(10),
			ExpectedBytes: []byte(`{"type":"uint","value":10}`),
		},

		{
			input:         TaggedNumericExpression{NumericOperator: GE, NumericType: NumericType{val: uint(100)}},
			ExpectedBytes: []byte(`{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":100}}}`),
		},
	} {

		t.Run("test", func(t *testing.T) {
			tcb := &TeeTcbEvalNumber{val: tv.input}

			data, err := tcb.MarshalJSON()
			require.NoError(t, err)
			fmt.Printf("received string %s", string(data))
			assert.Equal(t, tv.ExpectedBytes, data)

			out := &TeeTcbEvalNumber{}
			err = out.UnmarshalJSON(data)
			require.NoError(t, err)
		})
	}
}
