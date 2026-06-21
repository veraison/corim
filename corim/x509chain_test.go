// Copyright 2021-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/testdata"
)

func trustedRootsWithCert(t *testing.T, der []byte) TrustedRoots {
	t.Helper()

	root, err := x509.ParseCertificate(der)
	require.NoError(t, err)

	pool := x509.NewCertPool()
	pool.AddCert(root)

	return TrustedRoots{
		Pool: pool,
	}
}

type testPKI struct {
	rootKey         *ecdsa.PrivateKey
	root            *x509.Certificate
	intermediateKey *ecdsa.PrivateKey
	intermediate    *x509.Certificate
	intermediateDER []byte
	leafKey         *ecdsa.PrivateKey
	leaf            *x509.Certificate
	leafDER         []byte
}

func buildTestPKI(t *testing.T) testPKI {
	t.Helper()

	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	rootDER := mustCreateCA(t, rootKey, "Root CA")
	root, err := x509.ParseCertificate(rootDER)
	require.NoError(t, err)

	intermediateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	intermediateTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{CommonName: "Intermediate CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	intermediateDER, err := x509.CreateCertificate(
		rand.Reader, intermediateTemplate, root, &intermediateKey.PublicKey, rootKey,
	)
	require.NoError(t, err)

	intermediate, err := x509.ParseCertificate(intermediateDER)
	require.NoError(t, err)

	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	leafTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(3),
		Subject:               pkix.Name{CommonName: "Leaf"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	leafDER, err := x509.CreateCertificate(
		rand.Reader, leafTemplate, intermediate, &leafKey.PublicKey, intermediateKey,
	)
	require.NoError(t, err)

	leaf, err := x509.ParseCertificate(leafDER)
	require.NoError(t, err)

	return testPKI{
		rootKey:         rootKey,
		root:            root,
		intermediateKey: intermediateKey,
		intermediate:    intermediate,
		intermediateDER: intermediateDER,
		leafKey:         leafKey,
		leaf:            leaf,
		leafDER:         leafDER,
	}
}

func mustECPrivateKeyJWK(t *testing.T, key *ecdsa.PrivateKey) []byte {
	t.Helper()

	pad32 := func(b []byte) []byte {
		out := make([]byte, 32)
		copy(out[32-len(b):], b)

		return out
	}

	jwk := map[string]string{
		"kty": "EC",
		"crv": "P-256",
		"x":   base64.RawURLEncoding.EncodeToString(pad32(key.X.Bytes())),
		"y":   base64.RawURLEncoding.EncodeToString(pad32(key.Y.Bytes())),
		"d":   base64.RawURLEncoding.EncodeToString(pad32(key.D.Bytes())),
	}

	out, err := json.Marshal(jwk)
	require.NoError(t, err)

	return out
}

func signWithChain(t *testing.T, keyJWK, leafDER, intermediates []byte) (cbor []byte, signedIn, signedOut SignedCorim) {
	t.Helper()

	signer, err := NewSignerFromJWK(keyJWK)
	require.NoError(t, err)

	signedIn.UnsignedCorim = *unsignedCorimFromCBOR(t, testGoodUnsignedCorimCBOR)
	signedIn.Meta = *metaGood(t)
	require.NoError(t, signedIn.AddSigningCert(leafDER))

	if len(intermediates) > 0 {
		require.NoError(t, signedIn.AddIntermediateCerts(intermediates))
	}

	var errSign error
	cbor, errSign = signedIn.Sign(signer)
	require.NoError(t, errSign)

	require.NoError(t, signedOut.FromCOSE(cbor))

	return cbor, signedIn, signedOut
}

func TestSignedCorim_VerifyX509ChainTrust_ok(t *testing.T) {
	_, _, SignedCorimOut := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())

	trusted := trustedRootsWithCert(t, testdata.RootCA)
	err := SignedCorimOut.VerifyX509ChainTrust(trusted)
	assert.NoError(t, err)

	err = SignedCorimOut.VerifyX509ChainTrust(trusted)
	assert.NoError(t, err, "second VerifyX509ChainTrust call should be idempotent")
}

func TestCheckCRLValidity_notYetValid(t *testing.T) {
	crl := &x509.RevocationList{
		ThisUpdate: time.Now().Add(time.Hour),
		NextUpdate: time.Now().Add(2 * time.Hour),
	}

	err := checkCRLValidity(crl, time.Now())
	assert.ErrorContains(t, err, "not yet valid")
}

func TestParseX509CertificateDEROrPEM_invalidDER(t *testing.T) {
	_, err := ParseX509CertificateDEROrPEM([]byte("not a certificate"))
	assert.ErrorContains(t, err, "parsing certificate")
}

func TestParseRevocationListDEROrPEM_invalidDER(t *testing.T) {
	_, err := ParseRevocationListDEROrPEM([]byte("not a crl"))
	assert.ErrorContains(t, err, "parsing CRL")
}

func TestSignedCorim_VerifyX509ChainTrust_noX5Chain(t *testing.T) {
	s := NewSignedCorim()
	err := s.VerifyX509ChainTrust(TrustedRoots{})
	assert.EqualError(t, err, "x5chain: header not set in CoRIM")
}

func TestSignedCorim_VerifyX509ChainTrust_untrustedRoot(t *testing.T) {
	_, _, SignedCorimOut := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())

	t.Run("empty explicit pool", func(t *testing.T) {
		err := SignedCorimOut.VerifyX509ChainTrust(TrustedRoots{Pool: x509.NewCertPool()})
		assert.ErrorContains(t, err, "x5chain verification failed")
		assert.ErrorContains(t, err, "chain does not anchor to supplied root certificate(s)")
	})

	t.Run("system roots", func(t *testing.T) {
		err := SignedCorimOut.VerifyX509ChainTrust(TrustedRoots{})
		assert.ErrorContains(t, err, "x5chain verification failed")
		assert.ErrorContains(t, err, "chain does not anchor to trusted roots (system CAs)")
	})
}

func TestSignedCorim_VerifyX509ChainTrust_intermediateOnlyChain(t *testing.T) {
	_, _, SignedCorimOut := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, testdata.IntermediateCA)

	err := SignedCorimOut.VerifyX509ChainTrust(trustedRootsWithCert(t, testdata.RootCA))
	assert.NoError(t, err)
}

func TestSignedCorim_VerifyX509ChainTrust_expired(t *testing.T) {
	pki := buildTestPKI(t)

	shortLeafTemplate := &x509.Certificate{
		SerialNumber:          pki.leaf.SerialNumber,
		Subject:               pki.leaf.Subject,
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(30 * time.Minute),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	shortLeafDER, err := x509.CreateCertificate(
		rand.Reader, shortLeafTemplate, pki.intermediate, &pki.leafKey.PublicKey, pki.intermediateKey,
	)
	require.NoError(t, err)

	shortLeaf, err := x509.ParseCertificate(shortLeafDER)
	require.NoError(t, err)

	require.True(t, pki.intermediate.NotAfter.After(shortLeaf.NotAfter))
	require.True(t, pki.root.NotAfter.After(shortLeaf.NotAfter))

	_, _, SignedCorimOut := signWithChain(
		t, mustECPrivateKeyJWK(t, pki.leafKey), shortLeafDER, pki.intermediateDER,
	)

	pool := x509.NewCertPool()
	pool.AddCert(pki.root)

	err = SignedCorimOut.VerifyX509ChainTrust(TrustedRoots{
		Pool:        pool,
		CurrentTime: shortLeaf.NotAfter.Add(time.Hour),
	})
	assert.ErrorContains(t, err, "expired")
}

func TestSignedCorim_VerifyX509ChainTrust_signingCertIsCA(t *testing.T) {
	root, err := x509.ParseCertificate(testdata.RootCA)
	require.NoError(t, err)

	// Exercises validateSigningCertificate before PKIX/COSE verify.
	s := SignedCorim{SigningCert: root}
	err = s.VerifyX509ChainTrust(trustedRootsWithCert(t, testdata.RootCA))
	assert.EqualError(t, err, "x5chain: signing certificate must not be a CA")
}

func TestSignedCorim_VerifyX509ChainTrust_signingCertIsCA_fromCOSE(t *testing.T) {
	_, _, SignedCorimOut := signWithChain(t, testEndEntityKey, testdata.RootCA, nil)

	err := SignedCorimOut.VerifyX509ChainTrust(trustedRootsWithCert(t, testdata.RootCA))
	assert.EqualError(t, err, "x5chain: signing certificate must not be a CA")
}

func TestSignedCorim_VerifyX509ChainTrust_revokedIntermediateViaRootCRL(t *testing.T) {
	pki := buildTestPKI(t)

	crlDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(-time.Minute),
		NextUpdate: time.Now().Add(time.Hour),
		RevokedCertificateEntries: []x509.RevocationListEntry{
			{
				SerialNumber:   pki.intermediate.SerialNumber,
				RevocationTime: time.Now().Add(-time.Minute),
			},
		},
	}, pki.root, pki.rootKey)
	require.NoError(t, err)

	crl, err := x509.ParseRevocationList(crlDER)
	require.NoError(t, err)

	_, _, SignedCorimOut := signWithChain(
		t, mustECPrivateKeyJWK(t, pki.leafKey), pki.leafDER, pki.intermediateDER,
	)

	pool := x509.NewCertPool()
	pool.AddCert(pki.root)

	err = SignedCorimOut.VerifyX509ChainTrust(TrustedRoots{
		Pool: pool,
		CRLs: []*x509.RevocationList{crl},
	})
	assert.ErrorContains(t, err, "revoked")
}

func TestSignedCorim_VerifyX509ChainTrust_revokedLeafUsesVerifiedChain(t *testing.T) {
	pki := buildTestPKI(t)

	unrelatedKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	unrelatedCA, err := x509.ParseCertificate(mustCreateCA(t, unrelatedKey, "Unrelated CA"))
	require.NoError(t, err)

	crlDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(-time.Minute),
		NextUpdate: time.Now().Add(time.Hour),
		RevokedCertificateEntries: []x509.RevocationListEntry{
			{
				SerialNumber:   pki.leaf.SerialNumber,
				RevocationTime: time.Now().Add(-time.Minute),
			},
		},
	}, pki.intermediate, pki.intermediateKey)
	require.NoError(t, err)

	crl, err := x509.ParseRevocationList(crlDER)
	require.NoError(t, err)

	intermediates := append(append([]byte{}, unrelatedCA.Raw...), pki.intermediateDER...)
	_, _, SignedCorimOut := signWithChain(
		t, mustECPrivateKeyJWK(t, pki.leafKey), pki.leafDER, intermediates,
	)

	pool := x509.NewCertPool()
	pool.AddCert(pki.root)

	err = SignedCorimOut.VerifyX509ChainTrust(TrustedRoots{
		Pool: pool,
		CRLs: []*x509.RevocationList{crl},
	})
	assert.ErrorContains(t, err, "revoked")
}

func TestSignedCorim_VerifyX509ChainTrust_tamperedPayload(t *testing.T) {
	cbor, _, SignedCorimOut := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())
	cbor[len(cbor)-1] ^= 0xff
	require.NoError(t, SignedCorimOut.FromCOSE(cbor))

	err := SignedCorimOut.VerifyX509ChainTrust(trustedRootsWithCert(t, testdata.RootCA))
	assert.ErrorContains(t, err, "x5chain: COSE signature verification failed")
	assert.ErrorContains(t, err, "verification error")
}

func TestSignedCorim_VerifyX509ChainTrust_signingKeyMismatch(t *testing.T) {
	pki := buildTestPKI(t)

	wrongKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	_, _, SignedCorimOut := signWithChain(
		t, mustECPrivateKeyJWK(t, wrongKey), pki.leafDER, pki.intermediateDER,
	)

	pool := x509.NewCertPool()
	pool.AddCert(pki.root)

	err = SignedCorimOut.VerifyX509ChainTrust(TrustedRoots{Pool: pool})
	assert.ErrorContains(t, err, "x5chain: COSE signature verification failed")
	assert.ErrorContains(t, err, "verification error")
}

func TestCheckCRLValidity_expired(t *testing.T) {
	crl := &x509.RevocationList{
		ThisUpdate: time.Now().Add(-2 * time.Hour),
		NextUpdate: time.Now().Add(-time.Hour),
	}

	err := checkCRLValidity(crl, time.Now())
	assert.EqualError(t, err, `x5chain: CRL from "" has expired`)
}

func TestValidateSigningCertificate_missingDigitalSignature(t *testing.T) {
	cert := &x509.Certificate{
		IsCA:     false,
		KeyUsage: x509.KeyUsageCertSign,
	}

	err := validateSigningCertificate(cert)
	assert.EqualError(t, err, "x5chain: signing certificate lacks digitalSignature key usage")
}

func TestValidateSigningCertificate_zeroKeyUsagePasses(t *testing.T) {
	cert := &x509.Certificate{
		IsCA: false,
	}

	err := validateSigningCertificate(cert)
	assert.NoError(t, err)
}

func TestCheckRevocation_revokedReportedBeforeExpiredCRL(t *testing.T) {
	pki := buildTestPKI(t)

	validCRLDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(-time.Minute),
		NextUpdate: time.Now().Add(time.Hour),
		RevokedCertificateEntries: []x509.RevocationListEntry{
			{
				SerialNumber:   pki.leaf.SerialNumber,
				RevocationTime: time.Now().Add(-time.Minute),
			},
		},
	}, pki.intermediate, pki.intermediateKey)
	require.NoError(t, err)

	expiredCRLDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(2),
		ThisUpdate: time.Now().Add(-2 * time.Hour),
		NextUpdate: time.Now().Add(-time.Hour),
	}, pki.intermediate, pki.intermediateKey)
	require.NoError(t, err)

	validCRL, err := x509.ParseRevocationList(validCRLDER)
	require.NoError(t, err)

	expiredCRL, err := x509.ParseRevocationList(expiredCRLDER)
	require.NoError(t, err)

	chain := []*x509.Certificate{pki.leaf, pki.intermediate, pki.root}

	err = checkRevocation(chain, []*x509.RevocationList{expiredCRL, validCRL}, time.Now())
	assert.ErrorContains(t, err, "revoked")
	assert.NotContains(t, err.Error(), "has expired")
}

func TestCheckRevocation_validCRLIgnoresExpiredSibling(t *testing.T) {
	pki := buildTestPKI(t)

	validCRLDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(-time.Minute),
		NextUpdate: time.Now().Add(time.Hour),
	}, pki.intermediate, pki.intermediateKey)
	require.NoError(t, err)

	expiredCRLDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(2),
		ThisUpdate: time.Now().Add(-2 * time.Hour),
		NextUpdate: time.Now().Add(-time.Hour),
	}, pki.intermediate, pki.intermediateKey)
	require.NoError(t, err)

	validCRL, err := x509.ParseRevocationList(validCRLDER)
	require.NoError(t, err)

	expiredCRL, err := x509.ParseRevocationList(expiredCRLDER)
	require.NoError(t, err)

	chain := []*x509.Certificate{pki.leaf, pki.intermediate, pki.root}

	err = checkRevocation(chain, []*x509.RevocationList{expiredCRL, validCRL}, time.Now())
	assert.NoError(t, err)
}

func TestCheckRevocation_allMatchingCRLsExpired(t *testing.T) {
	pki := buildTestPKI(t)

	expiredCRL1DER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(-2 * time.Hour),
		NextUpdate: time.Now().Add(-time.Hour),
	}, pki.intermediate, pki.intermediateKey)
	require.NoError(t, err)

	expiredCRL2DER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(2),
		ThisUpdate: time.Now().Add(-3 * time.Hour),
		NextUpdate: time.Now().Add(-2 * time.Hour),
	}, pki.intermediate, pki.intermediateKey)
	require.NoError(t, err)

	expiredCRL1, err := x509.ParseRevocationList(expiredCRL1DER)
	require.NoError(t, err)

	expiredCRL2, err := x509.ParseRevocationList(expiredCRL2DER)
	require.NoError(t, err)

	chain := []*x509.Certificate{pki.leaf, pki.intermediate, pki.root}

	err = checkRevocation(chain, []*x509.RevocationList{expiredCRL1, expiredCRL2}, time.Now())
	assert.ErrorContains(t, err, "has expired")
}

func TestCheckRevocation_skipsNilCRL(t *testing.T) {
	pki := buildTestPKI(t)

	crlDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(-time.Minute),
		NextUpdate: time.Now().Add(time.Hour),
	}, pki.intermediate, pki.intermediateKey)
	require.NoError(t, err)

	crl, err := x509.ParseRevocationList(crlDER)
	require.NoError(t, err)

	chain := []*x509.Certificate{pki.leaf, pki.intermediate, pki.root}

	err = checkRevocation(chain, []*x509.RevocationList{nil, crl}, time.Now())
	assert.NoError(t, err)
}

func TestCheckRevocation_crlIssuerNotInChainIgnored(t *testing.T) {
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	chainCA, err := x509.ParseCertificate(mustCreateCA(t, caKey, "Chain CA"))
	require.NoError(t, err)

	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	leafTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "Test Leaf"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, chainCA, &leafKey.PublicKey, caKey)
	require.NoError(t, err)

	leaf, err := x509.ParseCertificate(leafDER)
	require.NoError(t, err)

	unrelatedKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	unrelatedCA, err := x509.ParseCertificate(mustCreateCA(t, unrelatedKey, "Unrelated CA"))
	require.NoError(t, err)

	crlDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(-time.Minute),
		NextUpdate: time.Now().Add(time.Hour),
	}, unrelatedCA, unrelatedKey)
	require.NoError(t, err)

	crl, err := x509.ParseRevocationList(crlDER)
	require.NoError(t, err)

	err = checkRevocation([]*x509.Certificate{leaf, chainCA}, []*x509.RevocationList{crl}, time.Now())
	assert.NoError(t, err, "CRL from an issuer outside the chain must be ignored")
}

func TestParseRevocationListDEROrPEM_pem(t *testing.T) {
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Test CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	require.NoError(t, err)

	ca, err := x509.ParseCertificate(caDER)
	require.NoError(t, err)

	crlDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(-time.Minute),
		NextUpdate: time.Now().Add(time.Hour),
	}, ca, caKey)
	require.NoError(t, err)

	pemCRL := pem.EncodeToMemory(&pem.Block{
		Type:  "X509 CRL",
		Bytes: crlDER,
	})

	crl, err := ParseRevocationListDEROrPEM(pemCRL)
	require.NoError(t, err)
	assert.NotNil(t, crl)
}

func TestParseRevocationListDEROrPEM_der(t *testing.T) {
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Test CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	require.NoError(t, err)

	ca, err := x509.ParseCertificate(caDER)
	require.NoError(t, err)

	crlDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(-time.Minute),
		NextUpdate: time.Now().Add(time.Hour),
	}, ca, caKey)
	require.NoError(t, err)

	crl, err := ParseRevocationListDEROrPEM(crlDER)
	require.NoError(t, err)
	assert.NotNil(t, crl)
}

func TestParseX509CertificateDEROrPEM_der(t *testing.T) {
	cert, err := ParseX509CertificateDEROrPEM(testdata.RootCA)
	require.NoError(t, err)
	assert.NotNil(t, cert)
}

func TestParseX509CertificateDEROrPEM_invalidPEMType(t *testing.T) {
	pemCert := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: testdata.RootCA,
	})

	_, err := ParseX509CertificateDEROrPEM(pemCert)
	assert.EqualError(t, err, `invalid PEM block type "RSA PRIVATE KEY"`)
}

func TestParseRevocationListDEROrPEM_invalidPEMType(t *testing.T) {
	pemCRL := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: []byte{0x01, 0x02},
	})

	_, err := ParseRevocationListDEROrPEM(pemCRL)
	assert.EqualError(t, err, `invalid PEM block type "CERTIFICATE"`)
}

func mustCreateCA(t *testing.T, key *ecdsa.PrivateKey, commonName string) []byte {
	t.Helper()

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: commonName},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)

	return der
}

func TestTrustedRootPool_readFileError(t *testing.T) {
	_, err := TrustedRootPool(func(string) ([]byte, error) {
		return nil, errors.New("read failed")
	}, []string{"missing.der"}, nil, true)
	assert.ErrorContains(t, err, "loading root certificate from missing.der")
	assert.ErrorContains(t, err, "read failed")
}

func TestTrustedRootPool_invalidRootParse(t *testing.T) {
	_, err := TrustedRootPool(func(string) ([]byte, error) {
		return []byte("not-a-cert"), nil
	}, []string{"bad.der"}, nil, false)
	assert.ErrorContains(t, err, "parsing root certificate from bad.der")
}

func TestTrustedRootPool_invalidCRLParse(t *testing.T) {
	_, err := TrustedRootPool(func(path string) ([]byte, error) {
		switch path {
		case "root.der":
			return testdata.RootCA, nil
		case "bad.crl":
			return []byte("not-a-crl"), nil
		default:
			t.Fatalf("unexpected path %q", path)
			return nil, nil
		}
	}, []string{"root.der"}, []string{"bad.crl"}, false)
	assert.ErrorContains(t, err, "parsing CRL from bad.crl")
}

func TestSelectVerifiedChain_prefersPresentedOverlap(t *testing.T) {
	leaf := &x509.Certificate{Raw: []byte("leaf")}
	inter := &x509.Certificate{Raw: []byte("inter")}
	rootA := &x509.Certificate{Raw: []byte("root-a")}
	rootB := &x509.Certificate{Raw: []byte("root-b")}

	presented := []*x509.Certificate{leaf, inter, rootA}
	chains := [][]*x509.Certificate{
		{leaf, inter, rootB},
		{leaf, inter, rootA},
	}

	selected := selectVerifiedChain(presented, chains)
	require.NotNil(t, selected)
	assert.Equal(t, rootA.Raw, selected[len(selected)-1].Raw)
}

func TestTrustedRootPool_dedupes_duplicate_roots(t *testing.T) {
	trusted, err := TrustedRootPool(func(string) ([]byte, error) {
		return testdata.RootCA, nil
	}, []string{"root-a.der", "root-b.der"}, nil, false)
	require.NoError(t, err)

	_, _, signed := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())
	err = signed.VerifyX509ChainTrust(trusted)
	assert.NoError(t, err)
}

func TestTrustedRootPool_excludesSystemRootsWhenDisabled(t *testing.T) {
	trusted, err := TrustedRootPool(func(path string) ([]byte, error) {
		if path != "root.der" {
			t.Fatalf("unexpected path %q", path)
		}

		return testdata.RootCA, nil
	}, []string{"root.der"}, nil, false)
	require.NoError(t, err)

	_, _, signed := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())
	assert.NoError(t, signed.VerifyX509ChainTrust(trusted))

	wrongRootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	wrongTrusted, err := TrustedRootPool(func(string) ([]byte, error) {
		return mustCreateCA(t, wrongRootKey, "Wrong Root CA"), nil
	}, []string{"wrong-root.der"}, nil, false)
	require.NoError(t, err)

	err = signed.VerifyX509ChainTrust(wrongTrusted)
	assert.ErrorContains(t, err, "x5chain verification failed")
}

func TestTrustedRootPool_includeSystemRoots(t *testing.T) {
	trusted, err := TrustedRootPool(func(string) ([]byte, error) {
		return testdata.RootCA, nil
	}, []string{"root.der"}, nil, true)
	require.NoError(t, err)
	require.NotNil(t, trusted.Pool)

	_, _, signed := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())
	assert.NoError(t, signed.VerifyX509ChainTrust(trusted))
}

func TestTrustedRootPool_emptyRootPathsWithSystemRoots(t *testing.T) {
	trusted, err := TrustedRootPool(func(string) ([]byte, error) {
		t.Fatal("readFile should not be called when rootPaths is empty")
		return nil, nil
	}, nil, nil, true)
	require.NoError(t, err)
	require.NotNil(t, trusted.Pool)
}

func TestTrustedRootPool_loadsCRLs(t *testing.T) {
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	caDER := mustCreateCA(t, caKey, "Test CA")
	ca, err := x509.ParseCertificate(caDER)
	require.NoError(t, err)

	crlDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(-time.Minute),
		NextUpdate: time.Now().Add(time.Hour),
	}, ca, caKey)
	require.NoError(t, err)

	trusted, err := TrustedRootPool(func(path string) ([]byte, error) {
		switch path {
		case "root.der":
			return testdata.RootCA, nil
		case "issuer.crl":
			return crlDER, nil
		default:
			t.Fatalf("unexpected path %q", path)
			return nil, nil
		}
	}, []string{"root.der"}, []string{"issuer.crl"}, false)
	require.NoError(t, err)

	require.Len(t, trusted.CRLs, 1)
	assert.Equal(t, crlDER, trusted.CRLs[0].Raw)
}

func TestSignedCorim_VerifyX509ChainTrust_wrongExplicitRootFails(t *testing.T) {
	_, _, SignedCorimOut := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())

	wrongRootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	wrongRootDER := mustCreateCA(t, wrongRootKey, "Wrong Root CA")

	trusted, err := TrustedRootPool(func(string) ([]byte, error) {
		return wrongRootDER, nil
	}, []string{"wrong-root.der"}, nil, false)
	require.NoError(t, err)

	err = SignedCorimOut.VerifyX509ChainTrust(trusted)
	assert.ErrorContains(t, err, "x5chain verification failed")
	assert.ErrorContains(t, err, "chain does not anchor to supplied root certificate(s)")
}

func TestSignedCorim_VerifyX509ChainTrust_crlNotYetValidFails(t *testing.T) {
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	caDER := mustCreateCA(t, caKey, "Test CA")
	ca, err := x509.ParseCertificate(caDER)
	require.NoError(t, err)

	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	leafTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{CommonName: "Test Leaf"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, ca, &leafKey.PublicKey, caKey)
	require.NoError(t, err)

	crlDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(time.Hour),
		NextUpdate: time.Now().Add(2 * time.Hour),
	}, ca, caKey)
	require.NoError(t, err)

	crl, err := x509.ParseRevocationList(crlDER)
	require.NoError(t, err)

	_, _, SignedCorimOut := signWithChain(t, mustECPrivateKeyJWK(t, leafKey), leafDER, caDER)

	pool := x509.NewCertPool()
	pool.AddCert(ca)

	err = SignedCorimOut.VerifyX509ChainTrust(TrustedRoots{
		Pool: pool,
		CRLs: []*x509.RevocationList{crl},
	})
	assert.ErrorContains(t, err, "not yet valid")
}

func TestParseX509CertificateDEROrPEM_pem(t *testing.T) {
	pemCert := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: testdata.RootCA,
	})

	cert, err := ParseX509CertificateDEROrPEM(pemCert)
	require.NoError(t, err)
	assert.NotNil(t, cert)
}

func TestParseX509CertificateDEROrPEM_pemUsesFirstBlockOnly(t *testing.T) {
	first := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: testdata.RootCA,
	})
	second := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: testdata.IntermediateCA,
	})

	cert, err := ParseX509CertificateDEROrPEM(append(first, second...))
	require.NoError(t, err)

	expected, err := x509.ParseCertificate(testdata.RootCA)
	require.NoError(t, err)
	assert.Equal(t, expected.Raw, cert.Raw)
}

func TestParseRevocationListDEROrPEM_pemUsesFirstBlockOnly(t *testing.T) {
	pki := buildTestPKI(t)

	firstCRLDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(-time.Minute),
		NextUpdate: time.Now().Add(time.Hour),
	}, pki.intermediate, pki.intermediateKey)
	require.NoError(t, err)

	secondCRLDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(2),
		ThisUpdate: time.Now().Add(-time.Minute),
		NextUpdate: time.Now().Add(time.Hour),
	}, pki.root, pki.rootKey)
	require.NoError(t, err)

	first := pem.EncodeToMemory(&pem.Block{Type: "X509 CRL", Bytes: firstCRLDER})
	second := pem.EncodeToMemory(&pem.Block{Type: "X509 CRL", Bytes: secondCRLDER})

	crl, err := ParseRevocationListDEROrPEM(append(first, second...))
	require.NoError(t, err)
	assert.Equal(t, firstCRLDER, crl.Raw)
}
