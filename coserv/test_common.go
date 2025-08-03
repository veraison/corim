// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
)

var (
	testExpiry, _    = time.Parse("2006-01-02T15:04:05Z", "2030-12-13T18:30:02Z")
	testTimestamp, _ = time.Parse("2006-01-02T15:04:05Z", "2030-12-01T18:30:01Z")
	testAuthority    = []byte{0xab, 0xcd, 0xef}
	testBytes        = []byte{0x00, 0x11, 0x22, 0x33}
)

func readTestVectorSlice(t *testing.T, fname string) []byte {
	b, err := os.ReadFile(path.Join("testvectors", fname)) // nolint:gosec
	require.NoError(t, err)
	return b
}

func readTestVectorString(t *testing.T, fname string) string {
	b, err := os.ReadFile(path.Join("testvectors", fname)) // nolint:gosec
	require.NoError(t, err)
	return string(b)
}

func exampleClassSelector(t *testing.T) *EnvironmentSelector {
	class0 := comid.NewClassBytes(testBytes).
		SetVendor("Example Vendor").
		SetModel("Example Model")
	require.NotNil(t, class0)

	class1 := comid.NewClassUUID(comid.TestUUID)
	require.NotNil(t, class1)

	selector := NewEnvironmentSelector().
		AddClass(StatefulClass{Class: class0}).
		AddClass(StatefulClass{Class: class1})
	require.NotNil(t, selector)

	return selector
}

func exampleClassSelector2(t *testing.T) *EnvironmentSelector {
	class0 := comid.NewClassBytes(testBytes).
		SetVendor("Example Vendor").
		SetModel("Example Model")
	require.NotNil(t, class0)

	selector := NewEnvironmentSelector().
		AddClass(StatefulClass{Class: class0})
	require.NotNil(t, selector)

	return selector
}

func exampleInstanceSelector(t *testing.T) *EnvironmentSelector {
	instance0, err := comid.NewUEIDInstance(comid.TestUEID)
	require.NoError(t, err)

	instance1, err := comid.NewBytesInstance(comid.TestBytes)
	require.NoError(t, err)

	selector := NewEnvironmentSelector().
		AddInstance(StatefulInstance{Instance: instance0}).
		AddInstance(StatefulInstance{Instance: instance1})
	require.NotNil(t, selector)

	return selector
}

func exampleGroupSelector(t *testing.T) *EnvironmentSelector {
	group0, err := comid.NewBytesGroup(comid.TestBytes)
	require.NoError(t, err)

	group1, err := comid.NewUUIDGroup(comid.TestUUID)
	require.NoError(t, err)

	selector := NewEnvironmentSelector().
		AddGroup(StatefulGroup{Group: group0}).
		AddGroup(StatefulGroup{Group: group1})
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
		AddGroup(StatefulGroup{Group: group0}).
		AddInstance(StatefulInstance{Instance: instance0}).
		AddGroup(StatefulGroup{Group: group0})
	require.NotNil(t, selector)

	return selector
}

func badExampleEmptySelector(t *testing.T) *EnvironmentSelector {
	es := NewEnvironmentSelector()
	require.NotNil(t, es)
	return es
}

func exampleClassQuery(t *testing.T) *Query {
	qry, err := NewQuery(ArtifactTypeReferenceValues, *exampleClassSelector(t), ResultTypeCollectedArtifacts)
	require.NoError(t, err)
	return qry
}

func exampleReferenceValuesResultSet(t *testing.T) *ResultSet {
	env := comid.Environment{
		Class: comid.NewClassBytes(comid.TestBytes),
	}

	measurement, err := comid.NewUUIDMeasurement(comid.TestUUID)
	require.NoError(t, err)
	measurement.SetVersion("1.2.3", swid.VersionSchemeSemVer).SetMinSVN(2)

	measurements := comid.NewMeasurements().Add(measurement)

	refval := comid.ValueTriple{
		Environment:  env,
		Measurements: *measurements,
	}

	require.NoError(t, refval.Valid())

	authority, err := comid.NewCryptoKeyTaggedBytes(testAuthority)
	require.NoError(t, err)

	rvq := RefValQuad{
		Authorities: &[]comid.CryptoKey{*authority},
		RVTriple:    &refval,
	}

	rset := NewResultSet().
		SetExpiry(testExpiry).
		AddReferenceValues(rvq)

	return rset
}
