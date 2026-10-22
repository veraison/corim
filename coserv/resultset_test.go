// Copyright 2025-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/cmw"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/cots"
)

func TestResultSet_AddAttestationKeys(t *testing.T) {
	authority, err := comid.NewCryptoKeyTaggedBytes(testAuthority)
	require.NoError(t, err)

	akq := AKQuad{
		Authorities: comid.NewCryptoKeys().Add(authority),
		AKTriple: &comid.KeyTriple{
			Environment: comid.Environment{
				Class: comid.NewClassBytes(testBytes),
			},
			VerifKeys: comid.CryptoKeys{
				comid.MustNewPKIXBase64Key(comid.TestECPubKey),
			},
		},
	}

	rset := NewResultSet().SetExpiry(testExpiry).AddAttestationKeys(akq)
	assert.NotNil(t, rset)
}

func TestResultSet_AddEndorsedValues(t *testing.T) {
	authority, err := comid.NewCryptoKeyTaggedBytes(testAuthority)
	require.NoError(t, err)

	evq := EndValQuad{
		Authorities: comid.NewCryptoKeys().Add(authority),
		EVTriple: &comid.ValueTriple{
			Environment: comid.Environment{
				Class: comid.NewClassBytes(testBytes),
			},
			Measurements: *comid.NewMeasurements().Add(
				comid.MustNewUUIDMeasurement(comid.TestUUID).SetName("foo"),
			),
		},
	}

	rset := NewResultSet().SetExpiry(testExpiry).AddEndorsedValues(evq)
	assert.NotNil(t, rset)
}

func TestResultSet_AddConditionalEndorsementValues(t *testing.T) {
	authority, err := comid.NewCryptoKeyTaggedBytes(testAuthority)
	require.NoError(t, err)

	ceq := CondEndValQuad{
		Authorities: comid.NewCryptoKeys().Add(authority),
		CETriple: &comid.CondEndorseTriple{
			Conditions: *comid.NewStatefulEnvironments().Add(&comid.ValueTriple{
				Environment: comid.Environment{
					Class: comid.NewClassBytes(testBytes),
				},
				Measurements: *comid.NewMeasurements().Add(
					comid.MustNewUUIDMeasurement(comid.TestUUID).SetName("foo"),
				),
			}),
			Endorsements: *comid.NewValueTriples().Add(&comid.ValueTriple{
				Environment: comid.Environment{
					Class: comid.NewClassBytes(testBytes),
				},
				Measurements: *comid.NewMeasurements().Add(
					comid.MustNewUUIDMeasurement(comid.TestUUID).SetName("foo"),
				),
			}),
		},
	}

	rset := NewResultSet().SetExpiry(testExpiry).AddConditionalEndorsementValues(ceq)
	assert.NotNil(t, rset)
}

func TestResultSet_AddCoTS(t *testing.T) {
	authority, err := comid.NewCryptoKeyTaggedBytes(testAuthority)
	require.NoError(t, err)

	// Create a simple CoTS structure for testing
	cotsStore := cots.NewConciseTaStore()

	// Add a basic environment group with a class
	class := comid.NewClassBytes(testBytes)
	env := comid.Environment{
		Class: class,
	}
	eg := cots.EnvironmentGroup{}
	eg.SetEnvironment(env)
	cotsStore.AddEnvironmentGroup(eg)

	// Add trust anchor keys
	testCert := []byte{0x30, 0x82, 0x01, 0x00} // Simple test cert bytes
	tas := cots.NewTasAndCas()
	tas.AddTaCert(testCert)
	cotsStore.SetKeys(*tas)

	cotsStmt := CoTSStmt{
		Authorities: comid.NewCryptoKeys().Add(authority),
		CoTS:        cotsStore,
	}

	rset := NewResultSet().SetExpiry(testExpiry).AddCoTS(cotsStmt)
	assert.NotNil(t, rset)
	assert.NotNil(t, rset.TAS)
	assert.Equal(t, 1, len(*rset.TAS))
}

func TestResultSet_AddSourceArtifacts(t *testing.T) {
	cmw0, err := cmw.NewMonad("application/vnd.example.refvals", []byte{0x00, 0x01, 0x02, 0x03})
	require.NoError(t, err)

	rset := NewResultSet().SetExpiry(testExpiry).AddSourceArtifacts(*cmw0)
	assert.NotNil(t, rset)
}

func TestResultSet_SetRIMs(t *testing.T) {
	collection, err := cmw.NewCollection("")
	require.NoError(t, err)

	rset := NewResultSet().SetExpiry(testExpiry).SetRIMs(*collection)
	assert.NotNil(t, rset)
}

func TestResultSet_Valid_artifact_types(t *testing.T) {
	var (
		rvq = RefValQuad{}
		evq = EndValQuad{}
		ceq = CondEndValQuad{}
		akq = AKQuad{}
		tas = CoTSStmt{}
	)

	const mixed = "result set mixes results of different artifact types"

	tests := []struct {
		name   string
		build  func(*ResultSet)
		expErr string
	}{
		{"RVQ only", func(r *ResultSet) { r.AddReferenceValues(rvq) }, ""},
		{"EVQ only", func(r *ResultSet) { r.AddEndorsedValues(evq) }, ""},
		{"CEQ only", func(r *ResultSet) { r.AddConditionalEndorsementValues(ceq) }, ""},
		{"EVQ and CEQ", func(r *ResultSet) {
			r.AddEndorsedValues(evq).AddConditionalEndorsementValues(ceq)
		}, ""},
		{"AKQ only", func(r *ResultSet) { r.AddAttestationKeys(akq) }, ""},
		{"TAS only", func(r *ResultSet) { r.AddCoTS(tas) }, ""},
		{"AKQ and TAS", func(r *ResultSet) {
			r.AddAttestationKeys(akq).AddCoTS(tas)
		}, ""},
		{"empty", func(r *ResultSet) {}, ""},
		{"RVQ and EVQ", func(r *ResultSet) {
			r.AddReferenceValues(rvq).AddEndorsedValues(evq)
		}, mixed},
		{"RVQ and CEQ", func(r *ResultSet) {
			r.AddReferenceValues(rvq).AddConditionalEndorsementValues(ceq)
		}, mixed},
		{"RVQ and AKQ", func(r *ResultSet) {
			r.AddReferenceValues(rvq).AddAttestationKeys(akq)
		}, mixed},
		{"RVQ and TAS", func(r *ResultSet) {
			r.AddReferenceValues(rvq).AddCoTS(tas)
		}, mixed},
		{"EVQ and AKQ", func(r *ResultSet) {
			r.AddEndorsedValues(evq).AddAttestationKeys(akq)
		}, mixed},
		{"CEQ and TAS", func(r *ResultSet) {
			r.AddConditionalEndorsementValues(ceq).AddCoTS(tas)
		}, mixed},
		{"all groups", func(r *ResultSet) {
			r.AddReferenceValues(rvq).AddEndorsedValues(evq).AddAttestationKeys(akq)
		}, mixed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rset := NewResultSet().SetExpiry(testExpiry)
			tt.build(rset)

			err := rset.Valid()
			if tt.expErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expErr)
			}
		})
	}
}

func TestResultSet_Valid_source_artifacts(t *testing.T) {
	monad, err := cmw.NewMonad("application/vnd.example.refvals", []byte{0x00, 0x01})
	require.NoError(t, err)

	collection, err := cmw.NewCollection("")
	require.NoError(t, err)

	t.Run("all monads", func(t *testing.T) {
		rset := NewResultSet().SetExpiry(testExpiry).
			AddSourceArtifacts(*monad).
			AddSourceArtifacts(*monad)
		assert.NoError(t, rset.Valid())
	})

	t.Run("collection among monads", func(t *testing.T) {
		rset := NewResultSet().SetExpiry(testExpiry).
			AddSourceArtifacts(*monad).
			AddSourceArtifacts(*collection)
		assert.EqualError(t, rset.Valid(), "source artifact at index 1 must be a monad")
	})
}
