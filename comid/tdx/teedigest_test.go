// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
)

func getNewDigests() Digests {
	d := comid.NewDigests()
	d.AddDigest(swid.Sha256, comid.MustHexDecode(nil, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"))
	return *d
}
func TestTeeDigest_NewTeeDigest_OK(t *testing.T) {
	d := getNewDigests()
	_, err := NewTeeDigest(d)
	require.Nil(t, err)
}

func TestTeeDigest_NewTeeDigestExpr_OK(t *testing.T) {
	d := getNewDigests()
	_, err := NewTeeDigestExpr(NMEM, d)
	require.Nil(t, err)
}

func TestTeeDigest_NewTeeDigestExpr_NOK(t *testing.T) {
	expectedErr := "invalid operator : 1"
	d := getNewDigests()
	_, err := NewTeeDigestExpr(GT, d)
	assert.EqualError(t, err, expectedErr)
}

func TestTeeDigest_GetTeeDigest_OK(t *testing.T) {
	d := getNewDigests()
	dg := TeeDigest{val: d}
	d1, err := dg.GetDigest()
	require.NoError(t, err)
	b := d.Equal(d1)
	require.True(t, b)
}

func TestTeeDigest_AddTeeDigest_OK(t *testing.T) {
	d := getNewDigests()
	dg := TeeDigest{val: d}
	_, err := dg.AddTeeDigest(NOP, d)
	require.NoError(t, err)
}

func TestTeeDigest_AddTeeDigest_NOK(t *testing.T) {
	expectedErr := "operator mis-match TeeDigest Op: 6, Input Op: 5"
	d := getNewDigests()
	dg := TeeDigest{TaggedSetDigestExpression{SetOperator: MEM, SetDigest: SetDigest(d)}}
	_, err := dg.AddTeeDigest(NOP, d)
	assert.EqualError(t, err, expectedErr)
}

func TestTeeDigest_Marshal_UnMarshal_OK(t *testing.T) {
	d := getNewDigests()
	dg := TeeDigest{d}
	b, err := dg.MarshalCBOR()
	require.Nil(t, err)
	x := &TeeDigest{}
	err = x.UnmarshalCBOR(b)
	require.Nil(t, err)
	err = x.Valid()
	require.Nil(t, err)
}
