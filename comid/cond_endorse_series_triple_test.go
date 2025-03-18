// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/extensions"
)

func Test_CondEndorseSeries_NewCondEndorseSeriesTriples_OK(t *testing.T) {
	c := NewCondEndorseSeriesTriples()
	require.NotNil(t, c)
}

func Test_CondEndorseSeries_Valid_NOK(t *testing.T) {
	c := NewCondEndorseSeriesTriples()
	expectedErr := "error at index 0: stateful environment validation failed: environment validation failed: environment must not be empty"
	series := &CondEndorseSeriesTriple{}
	c.Add(series)
	err := c.Valid()
	assert.EqualError(t, err, expectedErr)
}

type testExtensions struct {
	TestSVN uint `cbor:"-72,keyasint,omitempty" json:"testsvn,omitempty"`
}

func Test_CondEndorseSeries_RegisterExtensions(t *testing.T) {
	extMap := extensions.NewMap().
		Add(ExtMval, &testExtensions{})
	series := &CondEndorseSeriesTriple{}
	err := series.RegisterExtensions(extMap)
	require.NoError(t, err)
}

func Test_CondEndorseSeries_RegisterExtensions_NOK(t *testing.T) {
	expectedErr := `condition: unexpected extension point: "ReferenceValue"`
	extMap := extensions.NewMap().
		Add(ExtReferenceValue, &testExtensions{})
	series := &CondEndorseSeriesTriple{}
	err := series.RegisterExtensions(extMap)
	assert.EqualError(t, err, expectedErr)
}
