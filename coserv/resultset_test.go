// Copyright 2025 Contributors to the Veraison project.
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
