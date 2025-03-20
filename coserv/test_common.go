// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
)

var (
	testExpiry, _ = time.Parse("2006-01-02T15:04:05Z", "2030-12-13T18:30:02Z")
)

func exampleClassSelector(t *testing.T) *EnvironmentSelector {
	class0 := comid.NewClassBytes(comid.TestBytes).
		SetVendor("Example Vendor").
		SetModel("Example Model")
	require.NotNil(t, class0)

	class1 := comid.NewClassUUID(comid.TestUUID)
	require.NotNil(t, class1)

	selector := NewEnvironmentSelector().
		AddClass(*class0).
		AddClass(*class1)
	require.NotNil(t, selector)

	return selector
}

func exampleInstanceSelector(t *testing.T) *EnvironmentSelector {
	instance0, err := comid.NewUEIDInstance(comid.TestUEID)
	require.NoError(t, err)

	instance1, err := comid.NewBytesInstance(comid.TestBytes)
	require.NoError(t, err)

	selector := NewEnvironmentSelector().
		AddInstance(*instance0).
		AddInstance(*instance1)
	require.NotNil(t, selector)

	return selector
}

func exampleGroupSelector(t *testing.T) *EnvironmentSelector {
	group0, err := comid.NewBytesGroup(comid.TestBytes)
	require.NoError(t, err)

	group1, err := comid.NewUUIDGroup(comid.TestUUID)
	require.NoError(t, err)

	selector := NewEnvironmentSelector().
		AddGroup(*group0).
		AddGroup(*group1)
	require.NotNil(t, selector)

	return selector
}

func badExampleMixedSelector(t *testing.T) *EnvironmentSelector {
	group0, err := comid.NewBytesGroup(comid.TestBytes)
	require.NoError(t, err)

	instance0, err := comid.NewUEIDInstance(comid.TestUEID)
	require.NoError(t, err)

	class0 := comid.NewClassUUID(comid.TestUUID)
	require.NotNil(t, class0)

	selector := NewEnvironmentSelector().
		AddGroup(*group0).
		AddInstance(*instance0).
		AddGroup(*group0)
	require.NotNil(t, selector)

	return selector
}

func badExampleEmptySelector(t *testing.T) *EnvironmentSelector {
	es := NewEnvironmentSelector()
	require.NotNil(t, es)
	return es
}

func exampleClassQuery(t *testing.T) *Query {
	qry, err := NewQuery(ArtifactTypeReferenceValues, *exampleClassSelector(t))
	require.NoError(t, err)
	return qry
}

func exampleInstanceQuery(t *testing.T) *Query {
	qry, err := NewQuery(ArtifactTypeEndorsedValues, *exampleInstanceSelector(t))
	require.NoError(t, err)
	return qry
}

func exampleGroupQuery(t *testing.T) *Query {
	qry, err := NewQuery(ArtifactTypeTrustAnchors, *exampleGroupSelector(t))
	require.NoError(t, err)
	return qry
}

func exampleReferenceValuesResultSet(t *testing.T) *ResultSet {
	e0 := comid.Environment{
		Class: comid.NewClassBytes(comid.TestBytes),
	}

	m00, err := comid.NewUUIDMeasurement(comid.TestUUID)
	require.NoError(t, err)
	m00.SetVersion("1.2.3", swid.VersionSchemeSemVer).SetMinSVN(2)

	m0 := comid.NewMeasurements().Add(m00)

	rv0 := comid.ValueTriple{
		Environment:  e0,
		Measurements: *m0,
	}

	require.NoError(t, rv0.Valid())

	rset := NewResultSet().
		SetExpiry(testExpiry).
		AddReferenceValues(rv0)

	return rset
}
