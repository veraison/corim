// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdvisoryIDs_NewTeeAvisoryIDsExpr_OK(t *testing.T) {
	a := TestAdvisoryIDs
	_, err := NewTeeAdvisoryIDsExpr(MEM, a)
	require.Nil(t, err)
}

func TestAdvisoryIDs_NewTeeAvisoryIDsString_OK(t *testing.T) {
	a := TestAdvisoryIDs
	_, err := NewTeeAdvisoryIDsString(a)
	require.Nil(t, err)
}

func TestAdvisoryIDs_NewTeeAvisoryIDsExpr_NOK(t *testing.T) {
	expectedErr := "invalid operator : 5"
	a := make([]string, len(TestAdvisoryIDs))
	_, err := NewTeeAdvisoryIDsExpr(NOP, a)
	assert.EqualError(t, err, expectedErr)
}

func TestAdvisoryIDs_AddAdvisoryIDsString_OK(t *testing.T) {
	a := TestAdvisoryIDs
	adv := &TeeAdvisoryIDs{val: a}
	_, err := adv.AddTeeAdvisoryIDs(NOP, []string{"abcd"})
	require.NoError(t, err)
}

func TestAdvisoryIDs_AddAdvisoryIDsExpr_OK(t *testing.T) {
	a := TestAdvisoryIDs
	adv := TeeAdvisoryIDs{val: TaggedSetStringExpression{SetOperator: MEM, SetString: a}}
	_, err := adv.AddTeeAdvisoryIDs(MEM, []string{"abcd"})
	require.NoError(t, err)
}

func TestAdvisoryIDs_Valid_OK(t *testing.T) {
	a := TestAdvisoryIDs
	ta := TeeAdvisoryIDs{val: a}
	err := ta.Valid()
	require.NoError(t, err)
	ta = TeeAdvisoryIDs{val: TaggedSetStringExpression{SetOperator: MEM, SetString: a}}
	err = ta.Valid()
	require.NoError(t, err)
}

func TestAdvisoryIDs_Valid_NOK(t *testing.T) {
	expectedErr := "TeeAdvisoryID not set"
	a := []string{}
	ta := TeeAdvisoryIDs{val: a}
	err := ta.Valid()
	assert.EqualError(t, err, expectedErr)
	expectedErr = "unknown type tdx.SetStringExpression for TeeAdvisoryIDs"
	ta = TeeAdvisoryIDs{val: SetStringExpression{SetOperator: NOP, SetString: []string{"abc"}}}
	err = ta.Valid()
	assert.EqualError(t, err, expectedErr)
}

func TestAdvisoryIDs_JSON(t *testing.T) {
	a := TestAdvisoryIDs

	for _, tv := range []struct {
		input         []string
		ExpectedBytes []byte
	}{
		{
			input:         []string{"SA-123"},
			ExpectedBytes: []byte(`{"type":"string","value":["SA-123"]}`),
		},
		{
			input:         a,
			ExpectedBytes: []byte(`{"type":"string","value":["SA-00078","SA-00077","SA-00079"]}`),
		},
	} {

		t.Run("test", func(t *testing.T) {
			ta, err := NewTeeAdvisoryIDsString(tv.input)
			require.NoError(t, err)

			data, err := ta.MarshalJSON()
			require.NoError(t, err)
			fmt.Printf("received string %s", string(data))
			assert.Equal(t, tv.ExpectedBytes, data)

			out := &TeeAdvisoryIDs{}
			err = out.UnmarshalJSON(data)
			require.NoError(t, err)
			assert.Equal(t, tv.input, out.val)
		})
	}
}
