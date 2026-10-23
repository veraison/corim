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

func testQuads(t *testing.T) (RefValQuad, EndValQuad, CondEndValQuad, AKQuad, CoTSStmt) {
	t.Helper()

	authority, err := comid.NewCryptoKeyTaggedBytes(testAuthority)
	require.NoError(t, err)

	authorities := comid.NewCryptoKeys().Add(authority)
	env := comid.Environment{Class: comid.NewClassBytes(testBytes)}

	newValueTriple := func() *comid.ValueTriple {
		return &comid.ValueTriple{
			Environment: env,
			Measurements: *comid.NewMeasurements().Add(
				comid.MustNewUUIDMeasurement(comid.TestUUID).SetName("foo"),
			),
		}
	}

	cotsStore := cots.NewConciseTaStore()
	eg := cots.EnvironmentGroup{}
	eg.SetEnvironment(env)
	cotsStore.AddEnvironmentGroup(eg)
	tas := cots.NewTasAndCas()
	tas.AddTaCert([]byte{0x30, 0x82, 0x01, 0x00})
	cotsStore.SetKeys(*tas)

	return RefValQuad{Authorities: authorities, RVTriple: newValueTriple()},
		EndValQuad{Authorities: authorities, EVTriple: newValueTriple()},
		CondEndValQuad{
			Authorities: authorities,
			CETriple: &comid.CondEndorseTriple{
				Conditions:   *comid.NewStatefulEnvironments().Add(newValueTriple()),
				Endorsements: *comid.NewValueTriples().Add(newValueTriple()),
			},
		},
		AKQuad{
			Authorities: authorities,
			AKTriple: &comid.KeyTriple{
				Environment: env,
				VerifKeys: comid.CryptoKeys{
					comid.MustNewPKIXBase64Key(comid.TestECPubKey),
				},
			},
		},
		CoTSStmt{Authorities: authorities, CoTS: cotsStore}
}

func TestResultSet_Valid_artifact_types(t *testing.T) {
	rvq, evq, ceq, akq, tas := testQuads(t)

	const (
		mixed  = "result set mixes results of different artifact types"
		evqCeq = "EVQ and CEQ must be present together"
		akqTas = "AKQ and TAS must be present together"
	)

	tests := []struct {
		name   string
		build  func(*ResultSet)
		expErr string
	}{
		{"RVQ only", func(r *ResultSet) { r.AddReferenceValues(rvq) }, ""},
		{"EVQ and CEQ", func(r *ResultSet) {
			r.AddEndorsedValues(evq).AddConditionalEndorsementValues(ceq)
		}, ""},
		{"AKQ and TAS", func(r *ResultSet) {
			r.AddAttestationKeys(akq).AddCoTS(tas)
		}, ""},
		{"EVQ without CEQ", func(r *ResultSet) { r.AddEndorsedValues(evq) }, evqCeq},
		{"CEQ without EVQ", func(r *ResultSet) {
			r.AddConditionalEndorsementValues(ceq)
		}, evqCeq},
		{"AKQ without TAS", func(r *ResultSet) { r.AddAttestationKeys(akq) }, akqTas},
		{"TAS without AKQ", func(r *ResultSet) { r.AddCoTS(tas) }, akqTas},
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

func TestResultSet_Valid_empty(t *testing.T) {
	rvq, _, _, _, _ := testQuads(t)

	t.Run("nothing", func(t *testing.T) {
		rset := NewResultSet().SetExpiry(testExpiry)
		assert.EqualError(t, rset.Valid(), "result set is empty")
	})

	t.Run("empty quad array is present", func(t *testing.T) {
		rset := NewResultSet().SetExpiry(testExpiry)
		rset.RVQ = &[]RefValQuad{}
		assert.NoError(t, rset.Valid())
	})

	t.Run("quads", func(t *testing.T) {
		rset := NewResultSet().SetExpiry(testExpiry).AddReferenceValues(rvq)
		assert.NoError(t, rset.Valid())
	})

	t.Run("source artifacts only", func(t *testing.T) {
		monad, err := cmw.NewMonad("application/vnd.example.refvals", []byte{0x00})
		require.NoError(t, err)

		rset := NewResultSet().SetExpiry(testExpiry).AddSourceArtifacts(*monad)
		assert.NoError(t, rset.Valid())
	})
}

func TestResultSet_Valid_rims(t *testing.T) {
	rvq, _, _, _, _ := testQuads(t)

	collection, err := cmw.NewCollection("")
	require.NoError(t, err)

	monad, err := cmw.NewMonad("application/vnd.example.refvals", []byte{0x00})
	require.NoError(t, err)

	const excl = "RIMs are exclusive with all other results"

	t.Run("RIMs only", func(t *testing.T) {
		rset := NewResultSet().SetExpiry(testExpiry).SetRIMs(*collection)
		assert.NoError(t, rset.Valid())
	})

	t.Run("RIMs not a collection", func(t *testing.T) {
		rset := NewResultSet().SetExpiry(testExpiry).SetRIMs(*monad)
		assert.EqualError(t, rset.Valid(), "RIMs CMW must be a collection")
	})

	t.Run("RIMs and quads", func(t *testing.T) {
		rset := NewResultSet().SetExpiry(testExpiry).
			SetRIMs(*collection).
			AddReferenceValues(rvq)
		assert.EqualError(t, rset.Valid(), excl)
	})

	t.Run("RIMs and source artifacts", func(t *testing.T) {
		rset := NewResultSet().SetExpiry(testExpiry).
			SetRIMs(*collection).
			AddSourceArtifacts(*monad)
		assert.EqualError(t, rset.Valid(), excl)
	})
}

func TestResultSet_Valid_quads(t *testing.T) {
	rvq, evq, ceq, akq, tas := testQuads(t)

	tests := []struct {
		name   string
		build  func(*ResultSet)
		expErr string
	}{
		{"RVQ nil triple", func(r *ResultSet) {
			r.AddReferenceValues(rvq).
				AddReferenceValues(RefValQuad{Authorities: rvq.Authorities})
		}, "invalid RVQ at index 1: missing reference value triple"},
		{"RVQ missing authorities", func(r *ResultSet) {
			r.AddReferenceValues(RefValQuad{RVTriple: rvq.RVTriple})
		}, "invalid RVQ at index 0: missing authorities"},
		{"RVQ empty authorities", func(r *ResultSet) {
			r.AddReferenceValues(RefValQuad{
				Authorities: comid.NewCryptoKeys(),
				RVTriple:    rvq.RVTriple,
			})
		}, "invalid RVQ at index 0: invalid authorities: no keys to validate"},
		{"RVQ invalid triple", func(r *ResultSet) {
			r.AddReferenceValues(RefValQuad{
				Authorities: rvq.Authorities,
				RVTriple:    &comid.ValueTriple{},
			})
		}, "invalid RVQ at index 0: invalid reference value triple: "},
		{"EVQ nil triple", func(r *ResultSet) {
			r.AddEndorsedValues(EndValQuad{Authorities: evq.Authorities}).
				AddConditionalEndorsementValues(ceq)
		}, "invalid EVQ at index 0: missing endorsed value triple"},
		{"CEQ nil triple", func(r *ResultSet) {
			r.AddEndorsedValues(evq).
				AddConditionalEndorsementValues(ceq).
				AddConditionalEndorsementValues(CondEndValQuad{Authorities: ceq.Authorities})
		}, "invalid CEQ at index 1: missing conditional endorsement triple"},
		{"AKQ nil triple", func(r *ResultSet) {
			r.AddAttestationKeys(AKQuad{Authorities: akq.Authorities}).AddCoTS(tas)
		}, "invalid AKQ at index 0: missing attestation key triple"},
		{"TAS nil CoTS", func(r *ResultSet) {
			r.AddAttestationKeys(akq).AddCoTS(CoTSStmt{Authorities: tas.Authorities})
		}, "invalid TAS at index 0: missing CoTS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rset := NewResultSet().SetExpiry(testExpiry)
			tt.build(rset)

			err := rset.Valid()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expErr)
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

	t.Run("empty", func(t *testing.T) {
		rset := NewResultSet().SetExpiry(testExpiry)
		rset.SourceArtifacts = &[]cmw.CMW{}
		assert.EqualError(t, rset.Valid(), "source artifacts must not be empty")
	})
}
