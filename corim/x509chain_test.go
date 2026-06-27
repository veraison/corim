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

func trustAnchorsWithCert(t *testing.T, der []byte) TrustAnchors {
	t.Helper()

	anchor, err := x509.ParseCertificate(der)
	require.NoError(t, err)

	pool := x509.NewCertPool()
	pool.AddCert(anchor)

	return TrustAnchors{
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

func makeValidCRL(t *testing.T, issuer *x509.Certificate, issuerKey *ecdsa.PrivateKey) *x509.RevocationList {
	t.Helper()

	crlDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(-time.Minute),
		NextUpdate: time.Now().Add(time.Hour),
	}, issuer, issuerKey)
	require.NoError(t, err)

	crl, err := x509.ParseRevocationList(crlDER)
	require.NoError(t, err)

	return crl
}

func makeValidChainCRLs(t *testing.T, pki *testPKI) []*x509.RevocationList {
	t.Helper()

	return []*x509.RevocationList{
		makeValidCRL(t, pki.intermediate, pki.intermediateKey),
		makeValidCRL(t, pki.root, pki.rootKey),
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

func TestSignedCorim_VerifyWithX5Chain_ok(t *testing.T) {
	_, _, SignedCorimOut := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())

	err := SignedCorimOut.VerifyWithX5Chain(trustAnchorsWithCert(t, testdata.RootCA))
	assert.NoError(t, err)
}

func TestSignedCorim_VerifyWithX5Chain_noX5Chain(t *testing.T) {
	signer, err := NewSignerFromJWK(testEndEntityKey)
	require.NoError(t, err)

	var withoutChain SignedCorim
	withoutChain.UnsignedCorim = *unsignedCorimFromCBOR(t, testGoodUnsignedCorimCBOR)
	withoutChain.Meta = *metaGood(t)

	noChainCBOR, err := withoutChain.Sign(signer)
	require.NoError(t, err)

	var signed SignedCorim
	require.NoError(t, signed.FromCOSE(noChainCBOR))

	err = signed.VerifyWithX5Chain(TrustAnchors{})
	assert.EqualError(t, err, "x5chain: header not set in CoRIM")
}

func TestSignedCorim_VerifyWithX5Chain_noSign1Message(t *testing.T) {
	cert, err := x509.ParseCertificate(testdata.EndEntityDer)
	require.NoError(t, err)

	s := NewSignedCorim()
	s.SigningCert = cert

	err = s.VerifyWithX5Chain(trustAnchorsWithCert(t, testdata.RootCA))
	assert.EqualError(t, err, "no Sign1 message found")
}

func TestSignedCorim_VerifyWithX5Chain_untrustedAnchor(t *testing.T) {
	_, _, SignedCorimOut := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())

	err := SignedCorimOut.VerifyWithX5Chain(TrustAnchors{Pool: x509.NewCertPool()})
	assert.ErrorContains(t, err, "x5chain verification failed")
	var unknownAuthority x509.UnknownAuthorityError
	assert.ErrorAs(t, err, &unknownAuthority)
}

func TestSignedCorim_VerifyWithX5Chain_intermediateOnlyChain(t *testing.T) {
	_, _, SignedCorimOut := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, testdata.IntermediateCA)

	err := SignedCorimOut.VerifyWithX5Chain(trustAnchorsWithCert(t, testdata.RootCA))
	assert.NoError(t, err)
}

func TestSignedCorim_VerifyWithX5Chain_expired(t *testing.T) {
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

	err = SignedCorimOut.VerifyWithX5Chain(TrustAnchors{
		Pool:        pool,
		CurrentTime: shortLeaf.NotAfter.Add(time.Hour),
	})
	assert.ErrorContains(t, err, "expired")
}

func TestSignedCorim_VerifyWithX5Chain_signingCertIsCA_fromCOSE(t *testing.T) {
	_, _, SignedCorimOut := signWithChain(t, testEndEntityKey, testdata.RootCA, nil)

	err := SignedCorimOut.VerifyWithX5Chain(trustAnchorsWithCert(t, testdata.RootCA))
	assert.EqualError(t, err, "x5chain: signing certificate must not be a CA")
}

func TestSignedCorim_VerifyWithX5Chain_revokedIntermediateViaRootCRL(t *testing.T) {
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

	err = SignedCorimOut.VerifyWithX5Chain(TrustAnchors{
		Pool: pool,
		CRLs: []*x509.RevocationList{
			makeValidCRL(t, pki.intermediate, pki.intermediateKey),
			crl,
		},
	})
	assert.ErrorContains(t, err, "revoked")
}

func TestSignedCorim_VerifyWithX5Chain_revokedLeafUsesVerifiedChain(t *testing.T) {
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

	err = SignedCorimOut.VerifyWithX5Chain(TrustAnchors{
		Pool: pool,
		CRLs: []*x509.RevocationList{
			crl,
			makeValidCRL(t, pki.root, pki.rootKey),
		},
	})
	assert.ErrorContains(t, err, "revoked")
}

func TestSignedCorim_VerifyWithX5Chain_tamperedPayload(t *testing.T) {
	cbor, _, SignedCorimOut := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())
	cbor[len(cbor)-1] ^= 0xff
	require.NoError(t, SignedCorimOut.FromCOSE(cbor))

	err := SignedCorimOut.VerifyWithX5Chain(trustAnchorsWithCert(t, testdata.RootCA))
	assert.ErrorContains(t, err, "x5chain: COSE signature verification failed")
	assert.ErrorContains(t, err, "verification error")
}

func TestSignedCorim_VerifyWithX5Chain_signingKeyMismatch(t *testing.T) {
	pki := buildTestPKI(t)

	wrongKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	_, _, SignedCorimOut := signWithChain(
		t, mustECPrivateKeyJWK(t, wrongKey), pki.leafDER, pki.intermediateDER,
	)

	pool := x509.NewCertPool()
	pool.AddCert(pki.root)

	err = SignedCorimOut.VerifyWithX5Chain(TrustAnchors{Pool: pool})
	assert.ErrorContains(t, err, "x5chain: COSE signature verification failed")
	assert.ErrorContains(t, err, "verification error")
}

func TestValidateLeafSigningCert_missingDigitalSignature(t *testing.T) {
	cert := &x509.Certificate{
		IsCA:     false,
		KeyUsage: x509.KeyUsageCertSign,
	}

	err := validateLeafSigningCert(cert)
	assert.EqualError(t, err, "x5chain: signing certificate lacks digitalSignature key usage")
}

func TestValidateLeafSigningCert_zeroKeyUsagePasses(t *testing.T) {
	cert := &x509.Certificate{
		IsCA: false,
	}

	err := validateLeafSigningCert(cert)
	assert.NoError(t, err)
}

func TestCheckChainRevocation_revokedReportedBeforeExpiredCRL(t *testing.T) {
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
	crls := makeValidChainCRLs(t, &pki)
	crls[0] = validCRL
	crls = append(crls, expiredCRL)

	err = checkChainRevocation(chain, crls, CrlPolicyStrict, time.Now())
	assert.ErrorContains(t, err, "revoked")
	assert.NotContains(t, err.Error(), "has expired")
}

func TestCheckChainRevocation_validCRLIgnoresExpiredSibling(t *testing.T) {
	pki := buildTestPKI(t)

	expiredCRLDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(2),
		ThisUpdate: time.Now().Add(-2 * time.Hour),
		NextUpdate: time.Now().Add(-time.Hour),
	}, pki.intermediate, pki.intermediateKey)
	require.NoError(t, err)

	expiredCRL, err := x509.ParseRevocationList(expiredCRLDER)
	require.NoError(t, err)

	chain := []*x509.Certificate{pki.leaf, pki.intermediate, pki.root}
	crls := append(makeValidChainCRLs(t, &pki), expiredCRL)

	err = checkChainRevocation(chain, crls, CrlPolicyStrict, time.Now())
	assert.NoError(t, err)
}

func TestCheckChainRevocation_allMatchingCRLsExpired(t *testing.T) {
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
	crls := makeValidChainCRLs(t, &pki)
	crls[0] = expiredCRL1
	crls = append(crls, expiredCRL2)

	err = checkChainRevocation(chain, crls, CrlPolicyStrict, time.Now())
	assert.ErrorContains(t, err, "has expired")
}

func TestCheckChainRevocation_strictRejectsMissingNextUpdate(t *testing.T) {
	pki := buildTestPKI(t)

	crl := makeValidCRL(t, pki.intermediate, pki.intermediateKey)
	crl.NextUpdate = time.Time{}

	chain := []*x509.Certificate{pki.leaf, pki.intermediate, pki.root}

	err := checkChainRevocation(chain, []*x509.RevocationList{crl}, CrlPolicyStrict, time.Now())
	assert.ErrorContains(t, err, "no nextUpdate")
}

func TestCheckChainRevocation_permissiveAllowsMissingNextUpdate(t *testing.T) {
	pki := buildTestPKI(t)

	crl := makeValidCRL(t, pki.intermediate, pki.intermediateKey)
	crl.NextUpdate = time.Time{}

	chain := []*x509.Certificate{pki.leaf, pki.intermediate, pki.root}
	crls := []*x509.RevocationList{crl, makeValidCRL(t, pki.root, pki.rootKey)}

	err := checkChainRevocation(chain, crls, CrlPolicyPermissive, time.Now())
	assert.NoError(t, err)
}

func TestCheckChainRevocation_skipsNilCRL(t *testing.T) {
	pki := buildTestPKI(t)

	chain := []*x509.Certificate{pki.leaf, pki.intermediate, pki.root}

	err := checkChainRevocation(chain, append([]*x509.RevocationList{nil}, makeValidChainCRLs(t, &pki)...), CrlPolicyStrict, time.Now())
	assert.NoError(t, err)
}

func TestCheckChainRevocation_missingIssuerCRLFails(t *testing.T) {
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

	err = checkChainRevocation([]*x509.Certificate{leaf, chainCA}, []*x509.RevocationList{crl}, CrlPolicyStrict, time.Now())
	assert.ErrorContains(t, err, "unable to get certificate CRL")
}

func TestCheckChainRevocation_permissiveSkipsMissingIssuerCRL(t *testing.T) {
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

	err = checkChainRevocation(
		[]*x509.Certificate{leaf, chainCA},
		[]*x509.RevocationList{crl},
		CrlPolicyPermissive,
		time.Now(),
	)
	assert.NoError(t, err)
}

func TestCheckChainRevocation_okWithFullChainCRLs(t *testing.T) {
	pki := buildTestPKI(t)

	chain := []*x509.Certificate{pki.leaf, pki.intermediate, pki.root}

	err := checkChainRevocation(chain, makeValidChainCRLs(t, &pki), CrlPolicyStrict, time.Now())
	assert.NoError(t, err)
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

func TestLoadTrustAnchors_readFileError(t *testing.T) {
	_, err := LoadTrustAnchors(func(string) ([]byte, error) {
		return nil, errors.New("read failed")
	}, []string{"missing.der"}, nil)
	assert.ErrorContains(t, err, "loading trust anchor from missing.der")
	assert.ErrorContains(t, err, "read failed")
}

func TestLoadTrustAnchors_invalidTrustAnchorParse(t *testing.T) {
	_, err := LoadTrustAnchors(func(string) ([]byte, error) {
		return []byte("not-a-cert"), nil
	}, []string{"bad.der"}, nil)
	assert.ErrorContains(t, err, "parsing trust anchor from bad.der")
}

func TestLoadTrustAnchors_invalidCRLParse(t *testing.T) {
	_, err := LoadTrustAnchors(func(path string) ([]byte, error) {
		switch path {
		case "anchor.der":
			return testdata.RootCA, nil
		case "bad.crl":
			return []byte("not-a-crl"), nil
		default:
			t.Fatalf("unexpected path %q", path)
			return nil, nil
		}
	}, []string{"anchor.der"}, []string{"bad.crl"})
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

func TestLoadTrustAnchors_loadsPemTrustAnchorBundle(t *testing.T) {
	bundle := append(
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: testdata.RootCA}),
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: testdata.IntermediateCA})...,
	)

	anchors, err := LoadTrustAnchors(func(string) ([]byte, error) {
		return bundle, nil
	}, []string{"bundle.pem"}, nil)
	require.NoError(t, err)

	_, _, signed := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())
	assert.NoError(t, signed.VerifyWithX5Chain(anchors))
}

func TestLoadTrustAnchors_loadsPemCrlBundle(t *testing.T) {
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

	bundle := append(
		pem.EncodeToMemory(&pem.Block{Type: "X509 CRL", Bytes: firstCRLDER}),
		pem.EncodeToMemory(&pem.Block{Type: "X509 CRL", Bytes: secondCRLDER})...,
	)

	anchors, err := LoadTrustAnchors(func(path string) ([]byte, error) {
		switch path {
		case "anchor.der":
			return testdata.RootCA, nil
		case "crls.pem":
			return bundle, nil
		default:
			t.Fatalf("unexpected path %q", path)
			return nil, nil
		}
	}, []string{"anchor.der"}, []string{"crls.pem"})
	require.NoError(t, err)
	require.Len(t, anchors.CRLs, 2)
}

func TestLoadTrustAnchors_dedupesDuplicateTrustAnchors(t *testing.T) {
	anchors, err := LoadTrustAnchors(func(string) ([]byte, error) {
		return testdata.RootCA, nil
	}, []string{"anchor-a.der", "anchor-b.der"}, nil)
	require.NoError(t, err)

	_, _, signed := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())
	err = signed.VerifyWithX5Chain(anchors)
	assert.NoError(t, err)
}

func TestLoadTrustAnchors_wrongTrustAnchorFails(t *testing.T) {
	_, _, signed := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())

	wrongAnchorKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	wrongAnchors, err := LoadTrustAnchors(func(string) ([]byte, error) {
		return mustCreateCA(t, wrongAnchorKey, "Wrong Trust Anchor"), nil
	}, []string{"wrong-anchor.der"}, nil)
	require.NoError(t, err)

	err = signed.VerifyWithX5Chain(wrongAnchors)
	assert.ErrorContains(t, err, "x5chain verification failed")
	var unknownAuthority x509.UnknownAuthorityError
	assert.ErrorAs(t, err, &unknownAuthority)
}

func TestLoadTrustAnchors_emptyPathsUsesSystemStore(t *testing.T) {
	anchors, err := LoadTrustAnchors(func(string) ([]byte, error) {
		t.Fatal("readFile should not be called when trustAnchorPaths is empty")
		return nil, nil
	}, nil, nil)
	require.NoError(t, err)
	require.Nil(t, anchors.Pool)
}

func TestLoadTrustAnchors_loadsCRLs(t *testing.T) {
	pki := buildTestPKI(t)

	crlDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(-time.Minute),
		NextUpdate: time.Now().Add(time.Hour),
	}, pki.intermediate, pki.intermediateKey)
	require.NoError(t, err)

	anchors, err := LoadTrustAnchors(func(path string) ([]byte, error) {
		switch path {
		case "anchor.der":
			return testdata.RootCA, nil
		case "issuer.crl":
			return crlDER, nil
		default:
			t.Fatalf("unexpected path %q", path)
			return nil, nil
		}
	}, []string{"anchor.der"}, []string{"issuer.crl"})
	require.NoError(t, err)

	require.Len(t, anchors.CRLs, 1)
	assert.Equal(t, crlDER, anchors.CRLs[0].Raw)
}

func TestLoadTrustAnchors_invalidPemCrlType(t *testing.T) {
	pemCRL := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: []byte{0x01, 0x02}})

	_, err := LoadTrustAnchors(func(path string) ([]byte, error) {
		switch path {
		case "anchor.der":
			return testdata.RootCA, nil
		case "bad.crl":
			return pemCRL, nil
		default:
			t.Fatalf("unexpected path %q", path)
			return nil, nil
		}
	}, []string{"anchor.der"}, []string{"bad.crl"})
	assert.ErrorContains(t, err, "parsing CRL from bad.crl")
	assert.ErrorContains(t, err, `invalid PEM block type "CERTIFICATE"`)
}

func TestSignedCorim_VerifyWithX5Chain_crlNotYetValidFails(t *testing.T) {
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

	err = SignedCorimOut.VerifyWithX5Chain(TrustAnchors{
		Pool: pool,
		CRLs: []*x509.RevocationList{crl},
	})
	assert.ErrorContains(t, err, "not yet valid")
}
