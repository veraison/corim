// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/extensions"
)

func Test_CondEndSeries_NewCondEndSeriesTriples_OK(t *testing.T) {
	c := NewCondEndSeriesTriples()
	require.NotNil(t, c)
}

func Test_CondEndSeries_Valid_NOK(t *testing.T) {
	expectedErr := "Empty Collection"
	c := NewCondEndSeriesTriples()
	err := c.Valid()
	assert.EqualError(t, err, expectedErr)
}

func Test_CondEndSeries_Valid1_NOK(t *testing.T) {
	expectedErr := "error at index 0: stateful environment validation failed: environment validation failed: environment must not be empty"
	c := NewCondEndSeriesTriples()
	series := &CondEndSeriesTriple{}
	c.Add(series)
	err := c.Valid()
	assert.EqualError(t, err, expectedErr)
}

type testExtensions struct {
	TestSVN uint `cbor:"-72,keyasint,omitempty" json:"testsvn,omitempty"`
}

func Test_CondEndSeries_RegisterExtensions(t *testing.T) {
	extMap := extensions.NewMap().
		Add(ExtMval, &testExtensions{})
	series := &CondEndSeriesTriple{}
	err := series.RegisterExtensions(extMap)
	require.NoError(t, err)
}
