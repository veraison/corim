// Copyright 2021-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

// TrustAnchors holds trust material for x5chain validation.
// Prefer building via [LoadTrustAnchors].
//
// Pool semantics:
//   - nil — load OS trust store at verify time (used when no trust anchors are
//     supplied via [LoadTrustAnchors])
//   - non-nil pool — verify only against those anchors (explicit override; no
//     system roots)
//
// CRL semantics:
//   - empty CRLs — skip revocation checks
//   - non-empty CRLs — post-PKIX revocation checks; [CrlPolicy] selects strict vs
//     permissive behavior when an in-chain issuer has no matching CRL.
//     [CrlPolicyStrict] also requires each matching CRL to carry a nextUpdate that
//     has not passed.
//
// See [SignedCorim.VerifyWithX5Chain].
type TrustAnchors struct {
	Pool *x509.CertPool
	CRLs []*x509.RevocationList
	// CrlPolicy selects revocation behavior when CRLs is non-empty. Zero value is
	// [CrlPolicyStrict] (OpenSSL CRL_CHECK_ALL).
	CrlPolicy   CrlPolicy
	CurrentTime time.Time
}

// CrlPolicy selects how missing issuer CRLs are handled when CRLs is non-empty.
type CrlPolicy int

const (
	// CrlPolicyStrict requires every in-chain issuer to have a valid matching CRL
	// (OpenSSL CRL_CHECK_ALL) with a non-zero nextUpdate that has not passed.
	// This is the default (zero value).
	CrlPolicyStrict CrlPolicy = iota
	// CrlPolicyPermissive skips revocation for in-chain issuers with no matching CRL.
	// When matching CRLs exist but are all invalid, verification still fails.
	CrlPolicyPermissive
)

func newSystemCertPool() (*x509.CertPool, error) {
	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("loading system cert pool: %w", err)
	}
	if pool == nil {
		pool = x509.NewCertPool()
	}

	return pool, nil
}

func intermediatesFromChain(chain []*x509.Certificate) *x509.CertPool {
	pool := x509.NewCertPool()

	for i := 1; i < len(chain); i++ {
		pool.AddCert(chain[i])
	}

	return pool
}

func filterCRLsForIssuer(issuer *x509.Certificate, crls []*x509.RevocationList) []*x509.RevocationList {
	matched := make([]*x509.RevocationList, 0, len(crls))

	for _, crl := range crls {
		if crl == nil {
			continue
		}

		if crl.CheckSignatureFrom(issuer) == nil {
			matched = append(matched, crl)
		}
	}

	return matched
}

func checkCRLValidity(crl *x509.RevocationList, now time.Time, policy CrlPolicy) error {
	issuer := crl.Issuer.String()

	if !crl.ThisUpdate.IsZero() && now.Before(crl.ThisUpdate) {
		return fmt.Errorf("x5chain: CRL from %q is not yet valid", issuer)
	}

	if policy == CrlPolicyStrict && crl.NextUpdate.IsZero() {
		return fmt.Errorf("x5chain: CRL from %q has no nextUpdate", issuer)
	}

	if !crl.NextUpdate.IsZero() && now.After(crl.NextUpdate) {
		return fmt.Errorf("x5chain: CRL from %q has expired", issuer)
	}

	return nil
}

func isSerialRevoked(serial *big.Int, crl *x509.RevocationList) bool {
	for _, entry := range crl.RevokedCertificateEntries {
		if entry.SerialNumber.Cmp(serial) == 0 {
			return true
		}
	}

	return false
}

func checkChainRevocation(
	chain []*x509.Certificate,
	crls []*x509.RevocationList,
	policy CrlPolicy,
	now time.Time,
) error {
	if len(crls) == 0 {
		return nil
	}

	for i, cert := range chain {
		if i+1 >= len(chain) {
			break
		}

		issuer := chain[i+1]
		issuerCRLs := filterCRLsForIssuer(issuer, crls)
		if len(issuerCRLs) == 0 {
			if policy == CrlPolicyPermissive {
				continue
			}

			return fmt.Errorf("x5chain verification failed: unable to get certificate CRL")
		}

		var (
			validityErr   error
			validCRLFound bool
		)

		for _, crl := range issuerCRLs {
			if err := checkCRLValidity(crl, now, policy); err != nil {
				validityErr = err
				continue
			}

			validCRLFound = true

			if isSerialRevoked(cert.SerialNumber, crl) {
				return fmt.Errorf("x5chain: certificate %q is revoked", cert.Subject)
			}
		}

		if !validCRLFound {
			return validityErr
		}
	}

	return nil
}

// validateLeafSigningCert checks leaf signing-cert policy before PKIX.
// keyUsage is optional; when the extension is present, digitalSignature is required.
func validateLeafSigningCert(cert *x509.Certificate) error {
	if cert.IsCA {
		return fmt.Errorf("x5chain: signing certificate must not be a CA")
	}

	if cert.KeyUsage != 0 && cert.KeyUsage&x509.KeyUsageDigitalSignature == 0 {
		return fmt.Errorf("x5chain: signing certificate lacks digitalSignature key usage")
	}

	return nil
}

// selectVerifiedChain prefers the PKIX result with the most x5chain (DER) overlap.
func selectVerifiedChain(presented []*x509.Certificate, verifiedChains [][]*x509.Certificate) []*x509.Certificate {
	if len(verifiedChains) == 0 {
		return nil
	}

	bestIdx := 0
	bestScore := countDEROverlap(presented, verifiedChains[0])
	for i := 1; i < len(verifiedChains); i++ {
		if score := countDEROverlap(presented, verifiedChains[i]); score > bestScore {
			bestScore = score
			bestIdx = i
		}
	}

	return verifiedChains[bestIdx]
}

func countDEROverlap(presented, verified []*x509.Certificate) int {
	presentedDER := make(map[string]struct{}, len(presented))
	for _, cert := range presented {
		presentedDER[string(cert.Raw)] = struct{}{}
	}

	score := 0
	for _, cert := range verified {
		if _, ok := presentedDER[string(cert.Raw)]; ok {
			score++
		}
	}

	return score
}

// verifyPKIXChain validates chain against anchors. KeyUsages uses ExtKeyUsageAny
// (permissive EKU policy; full cert profile validation is out of scope).
func verifyPKIXChain(
	chain []*x509.Certificate,
	anchors TrustAnchors,
	now time.Time,
) ([]*x509.Certificate, error) {
	pool := anchors.Pool
	if pool == nil {
		var err error
		pool, err = newSystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("x5chain verification failed: %w", err)
		}
	}

	verifiedChains, err := chain[0].Verify(x509.VerifyOptions{
		Roots:         pool,
		Intermediates: intermediatesFromChain(chain),
		CurrentTime:   now,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	})
	if err != nil {
		return nil, fmt.Errorf("x5chain verification failed: %w", err)
	}
	if len(verifiedChains) == 0 {
		return nil, fmt.Errorf("x5chain verification failed: no verified chain")
	}

	return selectVerifiedChain(chain, verifiedChains), nil
}

func addTrustAnchorsFromDEROrPEM(pool *x509.CertPool, addedAnchors map[string]struct{}, data []byte) error {
	if pool.AppendCertsFromPEM(data) {
		return nil
	}

	cert, err := x509.ParseCertificate(data)
	if err != nil {
		return fmt.Errorf("parsing certificate: %w", err)
	}

	if _, seen := addedAnchors[string(cert.Raw)]; seen {
		return nil
	}

	addedAnchors[string(cert.Raw)] = struct{}{}
	pool.AddCert(cert)

	return nil
}

func crlsFromDEROrPEM(data []byte) ([]*x509.RevocationList, error) {
	block, rest := pem.Decode(data)
	if block == nil {
		crl, err := x509.ParseRevocationList(data)
		if err != nil {
			return nil, fmt.Errorf("parsing CRL: %w", err)
		}

		return []*x509.RevocationList{crl}, nil
	}

	crls := make([]*x509.RevocationList, 0, 1)

	for {
		if block.Type != "X509 CRL" {
			return nil, fmt.Errorf("invalid PEM block type %q", block.Type)
		}

		crl, err := x509.ParseRevocationList(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing CRL: %w", err)
		}

		crls = append(crls, crl)

		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
	}

	return crls, nil
}

// LoadTrustAnchors loads trust anchors and CRLs from files into a [TrustAnchors] value.
// PEM trust-anchor files may bundle multiple certificates
// ([x509.CertPool.AppendCertsFromPEM]); PEM CRL files may contain multiple blocks.
// Duplicate DER anchors in trustAnchorPaths are added once.
//
// When trustAnchorPaths is empty, Pool is nil and verification uses the OS trust
// store. When trustAnchorPaths is non-empty, only those anchors are trusted.
func LoadTrustAnchors(
	readFile func(string) ([]byte, error),
	trustAnchorPaths, crlPaths []string,
) (TrustAnchors, error) {
	anchors := TrustAnchors{
		CRLs: make([]*x509.RevocationList, 0, len(crlPaths)),
	}

	if len(trustAnchorPaths) > 0 {
		pool := x509.NewCertPool()
		anchors.Pool = pool
		addedAnchors := make(map[string]struct{})

		for _, path := range trustAnchorPaths {
			data, err := readFile(path)
			if err != nil {
				return TrustAnchors{}, fmt.Errorf("loading trust anchor from %s: %w", path, err)
			}

			if err := addTrustAnchorsFromDEROrPEM(pool, addedAnchors, data); err != nil {
				return TrustAnchors{}, fmt.Errorf("parsing trust anchor from %s: %w", path, err)
			}
		}
	}

	for _, path := range crlPaths {
		data, err := readFile(path)
		if err != nil {
			return TrustAnchors{}, fmt.Errorf("loading CRL from %s: %w", path, err)
		}

		crls, err := crlsFromDEROrPEM(data)
		if err != nil {
			return TrustAnchors{}, fmt.Errorf("parsing CRL from %s: %w", path, err)
		}

		anchors.CRLs = append(anchors.CRLs, crls...)
	}

	return anchors, nil
}
