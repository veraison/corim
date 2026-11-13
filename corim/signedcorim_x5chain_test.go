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
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/testdata"
	"github.com/veraison/go-cose"
)

func trustAnchorsWithCert(t *testing.T, der []byte) cose.TrustAnchors {
	t.Helper()

	anchor, err := x509.ParseCertificate(der)
	require.NoError(t, err)

	return cose.TrustAnchors{
		Anchors:        []*x509.Certificate{anchor},
		RevocationMode: cose.RevocationDisabled,
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

func mustGenerateKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	return key
}

func mustCreateCert(
	t *testing.T, template, parent *x509.Certificate, pub *ecdsa.PublicKey, parentKey *ecdsa.PrivateKey,
) (cert *x509.Certificate, der []byte) {
	t.Helper()

	der, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentKey)
	require.NoError(t, err)

	cert, err = x509.ParseCertificate(der)
	require.NoError(t, err)

	return cert, der
}

func caTemplate(serial int64, commonName string) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber:          big.NewInt(serial),
		Subject:               pkix.Name{CommonName: commonName},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
}

func leafTemplate(serial *big.Int, notAfter time.Time) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: "Leaf"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
}

// mustCreateRootCA returns a fresh self-signed root CA.
func mustCreateRootCA(t *testing.T, commonName string) (*x509.Certificate, *ecdsa.PrivateKey) {
	t.Helper()

	key := mustGenerateKey(t)
	template := caTemplate(1, commonName)
	cert, _ := mustCreateCert(t, template, template, &key.PublicKey, key)

	return cert, key
}

func buildTestPKI(t *testing.T) testPKI {
	t.Helper()

	var pki testPKI

	pki.root, pki.rootKey = mustCreateRootCA(t, "Root CA")

	pki.intermediateKey = mustGenerateKey(t)
	pki.intermediate, pki.intermediateDER = mustCreateCert(
		t, caTemplate(2, "Intermediate CA"), pki.root, &pki.intermediateKey.PublicKey, pki.rootKey,
	)

	pki.leafKey = mustGenerateKey(t)
	pki.leaf, pki.leafDER = mustCreateCert(
		t, leafTemplate(big.NewInt(3), time.Now().Add(time.Hour)),
		pki.intermediate, &pki.leafKey.PublicKey, pki.intermediateKey,
	)

	return pki
}

// sign returns a decoded SignedCorim signed by the PKI leaf key, carrying the
// leaf and intermediate in its x5chain.
func (p *testPKI) sign(t *testing.T) SignedCorim {
	t.Helper()

	_, _, signed := signWithChain(t, mustECPrivateKeyJWK(t, p.leafKey), p.leafDER, p.intermediateDER)

	return signed
}

// makeCRL issues a CRL valid for the next hour, revoking the given serials.
func makeCRL(
	t *testing.T, issuer *x509.Certificate, issuerKey *ecdsa.PrivateKey, revoked ...*big.Int,
) *x509.RevocationList {
	t.Helper()

	entries := make([]x509.RevocationListEntry, 0, len(revoked))
	for _, serial := range revoked {
		entries = append(entries, x509.RevocationListEntry{
			SerialNumber:   serial,
			RevocationTime: time.Now().Add(-time.Minute),
		})
	}

	crlDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:                    big.NewInt(1),
		ThisUpdate:                time.Now().Add(-time.Minute),
		NextUpdate:                time.Now().Add(time.Hour),
		RevokedCertificateEntries: entries,
	}, issuer, issuerKey)
	require.NoError(t, err)

	crl, err := x509.ParseRevocationList(crlDER)
	require.NoError(t, err)

	return crl
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

	err = signed.VerifyWithX5Chain(trustAnchorsWithCert(t, testdata.RootCA), nil)
	assert.ErrorContains(t, err, "header not set")
}

func TestSignedCorim_VerifyWithX5Chain_rejectsUnprotectedX5Chain(t *testing.T) {
	signer, err := NewSignerFromJWK(testEndEntityKey)
	require.NoError(t, err)

	var withoutChain SignedCorim
	withoutChain.UnsignedCorim = *unsignedCorimFromCBOR(t, testGoodUnsignedCorimCBOR)
	withoutChain.Meta = *metaGood(t)

	noChainCBOR, err := withoutChain.Sign(signer)
	require.NoError(t, err)

	// The signature does not cover the unprotected header, so no re-signing is needed.
	msg := cose.NewSign1Message()
	require.NoError(t, msg.UnmarshalCBOR(noChainCBOR))
	msg.Headers.RawUnprotected = nil
	msg.Headers.Unprotected[cose.HeaderLabelX5Chain] = testdata.EndEntityDer
	out, err := msg.MarshalCBOR()
	require.NoError(t, err)

	var signed SignedCorim
	require.NoError(t, signed.FromCOSE(out))
	err = signed.VerifyWithX5Chain(trustAnchorsWithCert(t, testdata.RootCA), nil)
	assert.ErrorContains(t, err, "unprotected header")
}

func TestSignedCorim_VerifyWithX5Chain_noSign1Message(t *testing.T) {
	err := NewSignedCorim().VerifyWithX5Chain(trustAnchorsWithCert(t, testdata.RootCA), nil)
	assert.EqualError(t, err, "no Sign1 message found")
}

func TestSignedCorim_VerifyWithX5Chain_tamperedSignature(t *testing.T) {
	signedCBOR, _, signed := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())
	msg := cose.NewSign1Message()
	require.NoError(t, msg.UnmarshalCBOR(signedCBOR))
	msg.Signature[0] ^= 0xff
	tamperedCBOR, err := msg.MarshalCBOR()
	require.NoError(t, err)
	require.NoError(t, signed.FromCOSE(tamperedCBOR))

	err = signed.VerifyWithX5Chain(trustAnchorsWithCert(t, testdata.RootCA), nil)
	assert.ErrorIs(t, err, ErrX5ChainSignature)
}

func TestSignedCorim_VerifyWithX5Chain_tamperedPayload(t *testing.T) {
	signedCBOR, _, signed := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())
	msg := cose.NewSign1Message()
	require.NoError(t, msg.UnmarshalCBOR(signedCBOR))

	unsigned := signed.UnsignedCorim
	require.NotNil(t, unsigned.SetID("tampered-corim"))
	payload, err := unsigned.ToCBOR()
	require.NoError(t, err)
	require.NotEqual(t, msg.Payload, payload)
	msg.Payload = payload

	// Preserve the original signature so the otherwise valid payload fails verification.
	tamperedCBOR, err := msg.MarshalCBOR()
	require.NoError(t, err)
	require.NoError(t, signed.FromCOSE(tamperedCBOR))

	err = signed.VerifyWithX5Chain(trustAnchorsWithCert(t, testdata.RootCA), nil)
	assert.ErrorIs(t, err, ErrX5ChainSignature)
}

func TestSignedCorim_VerifyWithX5Chain_signingCertIsCA(t *testing.T) {
	_, _, signed := signWithChain(t, testEndEntityKey, testdata.RootCA, nil)

	err := signed.VerifyWithX5Chain(trustAnchorsWithCert(t, testdata.RootCA), nil)
	assert.ErrorContains(t, err, "must not be a CA")
}

func TestSignedCorim_VerifyWithX5Chain_intermediateOnlyChain(t *testing.T) {
	_, _, signed := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, testdata.IntermediateCA)

	err := signed.VerifyWithX5Chain(trustAnchorsWithCert(t, testdata.RootCA), nil)
	assert.NoError(t, err)
}

func TestSignedCorim_VerifyWithX5Chain_noTrustAnchors(t *testing.T) {
	_, _, SignedCorimOut := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())

	err := SignedCorimOut.VerifyWithX5Chain(cose.TrustAnchors{}, nil)
	assert.ErrorIs(t, err, ErrX5ChainNoTrust)
}

func TestSignedCorim_VerifyWithX5Chain_expired(t *testing.T) {
	pki := buildTestPKI(t)

	shortLeaf, shortLeafDER := mustCreateCert(
		t, leafTemplate(pki.leaf.SerialNumber, time.Now().Add(30*time.Minute)),
		pki.intermediate, &pki.leafKey.PublicKey, pki.intermediateKey,
	)

	// Only the leaf is expired at verifyAt; both CAs are still valid.
	verifyAt := shortLeaf.NotAfter.Add(time.Minute)
	require.True(t, verifyAt.Before(pki.intermediate.NotAfter))
	require.True(t, verifyAt.Before(pki.root.NotAfter))

	_, _, SignedCorimOut := signWithChain(
		t, mustECPrivateKeyJWK(t, pki.leafKey), shortLeafDER, pki.intermediateDER,
	)

	anchors := cose.TrustAnchors{
		Anchors:        []*x509.Certificate{pki.root},
		RevocationMode: cose.RevocationDisabled,
	}
	opts := &cose.X5ChainVerifyOptions{CurrentTime: shortLeaf.NotAfter.Add(-time.Minute)}
	require.NoError(t, SignedCorimOut.VerifyWithX5Chain(anchors, opts))

	opts.CurrentTime = verifyAt
	assert.ErrorContains(t, SignedCorimOut.VerifyWithX5Chain(anchors, opts), "expired")
}

func TestSignedCorim_VerifyWithX5Chain_revokedIntermediateViaRootCRL(t *testing.T) {
	pki := buildTestPKI(t)
	signed := pki.sign(t)

	err := signed.VerifyWithX5Chain(cose.TrustAnchors{
		Anchors: []*x509.Certificate{pki.root},
		CRLs: []*x509.RevocationList{
			makeCRL(t, pki.intermediate, pki.intermediateKey),
			makeCRL(t, pki.root, pki.rootKey, pki.intermediate.SerialNumber),
		},
	}, nil)
	assert.ErrorIs(t, err, ErrX5ChainRevoked)
}

func TestSignedCorim_VerifyWithX5Chain_crlPolicyMissingIssuerCRL(t *testing.T) {
	pki := buildTestPKI(t)
	signed := pki.sign(t)

	anchors := cose.TrustAnchors{
		Anchors: []*x509.Certificate{pki.root},
		CRLs:    []*x509.RevocationList{makeCRL(t, pki.root, pki.rootKey)},
	}
	opts := &cose.X5ChainVerifyOptions{CRLPolicy: cose.CRLPolicyStrict}
	assert.ErrorIs(t, signed.VerifyWithX5Chain(anchors, opts), ErrX5ChainCRLMissing)

	opts.CRLPolicy = cose.CRLPolicyPermissive
	assert.NoError(t, signed.VerifyWithX5Chain(anchors, opts))
}

func TestSignedCorim_VerifyWithX5Chain_revocationLeafOnly(t *testing.T) {
	pki := buildTestPKI(t)
	signed := pki.sign(t)

	anchors := cose.TrustAnchors{
		Anchors: []*x509.Certificate{pki.root},
		CRLs:    []*x509.RevocationList{makeCRL(t, pki.intermediate, pki.intermediateKey)},
	}
	opts := &cose.X5ChainVerifyOptions{CRLPolicy: cose.CRLPolicyStrict}
	assert.ErrorIs(t, signed.VerifyWithX5Chain(anchors, opts), ErrX5ChainCRLMissing)

	anchors.RevocationMode = cose.RevocationLeafOnly
	assert.NoError(t, signed.VerifyWithX5Chain(anchors, opts))

	anchors.CRLs = []*x509.RevocationList{
		makeCRL(t, pki.intermediate, pki.intermediateKey, pki.leaf.SerialNumber),
	}
	assert.ErrorIs(t, signed.VerifyWithX5Chain(anchors, opts), ErrX5ChainRevoked)
}

func TestSignedCorim_VerifyWithX5Chain_revocationDisabledIgnoresCRLs(t *testing.T) {
	pki := buildTestPKI(t)
	signed := pki.sign(t)

	err := signed.VerifyWithX5Chain(cose.TrustAnchors{
		Anchors: []*x509.Certificate{pki.root},
		CRLs: []*x509.RevocationList{
			makeCRL(t, pki.intermediate, pki.intermediateKey, pki.leaf.SerialNumber),
			makeCRL(t, pki.root, pki.rootKey),
		},
		RevocationMode: cose.RevocationDisabled,
	}, nil)
	assert.NoError(t, err)
}

func TestSignedCorim_VerifyWithX5Chain_signingKeyMismatch(t *testing.T) {
	pki := buildTestPKI(t)

	_, _, SignedCorimOut := signWithChain(
		t, mustECPrivateKeyJWK(t, mustGenerateKey(t)), pki.leafDER, pki.intermediateDER,
	)

	err := SignedCorimOut.VerifyWithX5Chain(cose.TrustAnchors{
		Anchors:        []*x509.Certificate{pki.root},
		RevocationMode: cose.RevocationDisabled,
	}, nil)
	assert.ErrorIs(t, err, ErrX5ChainSignature)
}

func writeTestFile(t *testing.T, dir, name string, data []byte) string {
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	return path
}

// TestX5Chain_endToEnd drives only the exported API, the way a caller such as
// cocli does: sign a CoRIM with an x5chain, decode it, load trust anchors and
// CRLs from files, and verify.
func TestX5Chain_endToEnd(t *testing.T) {
	pki := buildTestPKI(t)

	var in SignedCorim
	require.NoError(t, in.UnsignedCorim.FromCBOR(testGoodUnsignedCorimCBOR))
	in.Meta = *NewMeta().SetSigner("ACME Ltd.", nil)
	require.NoError(t, in.AddSigningCert(pki.leafDER))
	require.NoError(t, in.AddIntermediateCerts(pki.intermediateDER))

	signer, err := cose.NewSigner(cose.AlgorithmES256, pki.leafKey)
	require.NoError(t, err)

	signedCBOR, err := in.Sign(signer)
	require.NoError(t, err)

	var out SignedCorim
	require.NoError(t, out.FromCOSE(signedCBOR))
	require.Equal(t, pki.leafDER, out.SigningCert.Raw)
	require.Len(t, out.IntermediateCerts, 1)

	dir := t.TempDir()
	rootPEM := writeTestFile(t, dir, "root.pem",
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: pki.root.Raw}))
	rootCRL := writeTestFile(t, dir, "root.crl", makeCRL(t, pki.root, pki.rootKey).Raw)

	t.Run("trusted chain verifies", func(t *testing.T) {
		intermediateCRL := writeTestFile(t, dir, "intermediate.crl",
			makeCRL(t, pki.intermediate, pki.intermediateKey).Raw)

		anchors, err := cose.LoadTrustAnchors(os.ReadFile, []string{rootPEM}, []string{intermediateCRL, rootCRL})
		require.NoError(t, err)

		assert.NoError(t, out.VerifyWithX5Chain(anchors, nil))
	})

	t.Run("revoked leaf is rejected", func(t *testing.T) {
		revokingCRL := writeTestFile(t, dir, "revoking.crl",
			makeCRL(t, pki.intermediate, pki.intermediateKey, pki.leaf.SerialNumber).Raw)

		anchors, err := cose.LoadTrustAnchors(os.ReadFile, []string{rootPEM}, []string{revokingCRL, rootCRL})
		require.NoError(t, err)

		assert.ErrorIs(t, out.VerifyWithX5Chain(anchors, nil), ErrX5ChainRevoked)
	})

	t.Run("untrusted anchor is rejected", func(t *testing.T) {
		other, _ := mustCreateRootCA(t, "Other Root CA")
		otherPEM := writeTestFile(t, dir, "other.pem",
			pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: other.Raw}))

		anchors, err := cose.LoadTrustAnchors(os.ReadFile, []string{otherPEM}, nil)
		require.NoError(t, err)
		anchors.RevocationMode = cose.RevocationDisabled

		assert.ErrorIs(t, out.VerifyWithX5Chain(anchors, nil), ErrX5ChainNoTrust)
	})
}

// signWithRawX5Chain signs a CoRIM with testEndEntityKey and puts x5chain into
// the protected header as given, bypassing the order check in Sign.
func signWithRawX5Chain(t *testing.T, x5chain interface{}) []byte {
	t.Helper()

	signer, err := NewSignerFromJWK(testEndEntityKey)
	require.NoError(t, err)

	var unsigned SignedCorim
	unsigned.UnsignedCorim = *unsignedCorimFromCBOR(t, testGoodUnsignedCorimCBOR)
	unsigned.Meta = *metaGood(t)

	noChainCBOR, err := unsigned.Sign(signer)
	require.NoError(t, err)

	msg := cose.NewSign1Message()
	require.NoError(t, msg.UnmarshalCBOR(noChainCBOR))
	msg.Headers.RawProtected = nil
	msg.Headers.Protected[cose.HeaderLabelX5Chain] = x5chain
	msg.Signature = nil
	require.NoError(t, msg.Sign(rand.Reader, NoExternalData, signer))

	out, err := msg.MarshalCBOR()
	require.NoError(t, err)

	return out
}

// unorderedX5Chain puts the root before the intermediate, so the leaf is not
// issued by the next certificate.
func unorderedX5Chain() [][]byte {
	return [][]byte{testdata.EndEntityDer, testdata.RootCA, testdata.IntermediateCA}
}

func TestSignedCorim_Sign_unorderedIntermediatesKeepsMessage(t *testing.T) {
	_, _, signed := signWithChain(t, testEndEntityKey, testdata.EndEntityDer, certChain())
	decoded := signed.message
	require.NotNil(t, decoded)

	// Root before intermediate: SetX5Chain rejects the order.
	require.NoError(t, signed.AddIntermediateCerts(append(append([]byte{}, testdata.RootCA...), testdata.IntermediateCA...)))

	signer, err := NewSignerFromJWK(testEndEntityKey)
	require.NoError(t, err)

	_, err = signed.Sign(signer)
	assert.ErrorContains(t, err, "setting x5chain")
	assert.ErrorContains(t, err, "not issued by")
	assert.Same(t, decoded, signed.message)

	pk, err := NewPublicKeyFromJWK(testEndEntityKey)
	require.NoError(t, err)
	assert.NoError(t, signed.Verify(pk))
}

func TestSignedCorim_FromCOSE_rejectsUnorderedChain(t *testing.T) {
	var signed SignedCorim

	err := signed.FromCOSE(signWithRawX5Chain(t, unorderedX5Chain()))
	assert.ErrorContains(t, err, "not issued by")
	assert.Nil(t, signed.message)
	assert.Nil(t, signed.SigningCert)
	assert.Empty(t, signed.IntermediateCerts)
}

func TestSignedCorim_extractX5Chain_assignsOnSuccess(t *testing.T) {
	var signed SignedCorim

	err := signed.extractX5Chain([][]byte{testdata.EndEntityDer, testdata.IntermediateCA, testdata.RootCA})
	require.NoError(t, err)
	assert.Equal(t, testdata.EndEntityDer, signed.SigningCert.Raw)
	require.Len(t, signed.IntermediateCerts, 2)
	assert.Equal(t, testdata.IntermediateCA, signed.IntermediateCerts[0].Raw)
	assert.Equal(t, testdata.RootCA, signed.IntermediateCerts[1].Raw)
}
