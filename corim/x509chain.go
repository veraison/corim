// Copyright 2021-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// TrustedRoots holds the trust material used to validate a CoRIM's embedded
// x5chain. Prefer constructing Pool and CRLs via [TrustedRootPool].
//
//   - Pool: trusted root CAs only (not intermediates) for PKIX path validation.
//     Three cases at verify time:
//     Pool == nil (e.g. zero TrustedRoots{}): load system roots;
//     empty non-nil pool (e.g. TrustedRootPool with includeSystemRoots false
//     and no rootPaths): explicit trust only, no system roots;
//     non-empty pool: use as-is (may include system roots when built via
//     TrustedRootPool with includeSystemRoots true).
//     Prefer [TrustedRootPool] when embedding the library so trust policy is
//     explicit. Do not rely on Pool == nil for CLI-style verification; use
//     TrustedRootPool with includeSystemRoots true and an empty rootPaths slice.
//   - CRLs: optional revocation lists. Nil or empty skips revocation checking.
//     Nil entries in the slice are ignored. CRLs whose issuer is not present in
//     the PKIX-verified chain are ignored. When multiple CRLs match an issuer,
//     invalid CRLs are skipped; revoked status is checked against every valid
//     CRL; if all matching CRLs are invalid, verification fails. Revocation is
//     evaluated on the PKIX-verified chain, not the raw x5chain order.
//   - CurrentTime: timestamp for cert/CRL validity checks; zero means time.Now().
type TrustedRoots struct {
	Pool        *x509.CertPool
	CRLs        []*x509.RevocationList
	CurrentTime time.Time
}

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

func crlsForIssuer(issuer *x509.Certificate, crls []*x509.RevocationList) []*x509.RevocationList {
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

func crlIssuerLabel(crl *x509.RevocationList) string {
	return crl.Issuer.String()
}

func checkCRLValidity(crl *x509.RevocationList, now time.Time) error {
	issuer := crlIssuerLabel(crl)

	if !crl.ThisUpdate.IsZero() && now.Before(crl.ThisUpdate) {
		return fmt.Errorf("x5chain: CRL from %q is not yet valid", issuer)
	}

	if !crl.NextUpdate.IsZero() && now.After(crl.NextUpdate) {
		return fmt.Errorf("x5chain: CRL from %q has expired", issuer)
	}

	return nil
}

func serialRevoked(serial *big.Int, crls []*x509.RevocationList) bool {
	for _, crl := range crls {
		for _, entry := range crl.RevokedCertificateEntries {
			if entry.SerialNumber.Cmp(serial) == 0 {
				return true
			}
		}
	}

	return false
}

func checkRevocation(chain []*x509.Certificate, crls []*x509.RevocationList, now time.Time) error {
	if len(crls) == 0 {
		return nil
	}

	for i, cert := range chain {
		if i+1 >= len(chain) {
			break
		}

		issuer := chain[i+1]
		issuerCRLs := crlsForIssuer(issuer, crls)
		if len(issuerCRLs) == 0 {
			continue
		}

		var (
			validityErr   error
			validCRLFound bool
		)

		for _, crl := range issuerCRLs {
			if err := checkCRLValidity(crl, now); err != nil {
				validityErr = err
				continue
			}

			validCRLFound = true

			if serialRevoked(cert.SerialNumber, []*x509.RevocationList{crl}) {
				return fmt.Errorf("x5chain: certificate %q is revoked", cert.Subject)
			}
		}

		if !validCRLFound && validityErr != nil {
			return validityErr
		}
	}

	return nil
}

// validateSigningCertificate checks leaf signing-cert policy before PKIX.
// When KeyUsage is absent (zero), verification proceeds per RFC 5280 (no
// keyUsage restriction implied); when present, digitalSignature must be set.
func validateSigningCertificate(cert *x509.Certificate) error {
	if cert.IsCA {
		return fmt.Errorf("x5chain: signing certificate must not be a CA")
	}

	if cert.KeyUsage != 0 && cert.KeyUsage&x509.KeyUsageDigitalSignature == 0 {
		return fmt.Errorf("x5chain: signing certificate lacks digitalSignature key usage")
	}

	return nil
}

// selectVerifiedChain picks one PKIX result when Verify returns multiple chains.
// Prefer the chain that shares the most certificates (by DER) with the presented
// x5chain, so verification follows the COSE header rather than an arbitrary path.
func selectVerifiedChain(presented []*x509.Certificate, verifiedChains [][]*x509.Certificate) []*x509.Certificate {
	if len(verifiedChains) == 0 {
		return nil
	}

	bestIdx := 0
	bestScore := verifiedChainOverlap(presented, verifiedChains[0])
	for i := 1; i < len(verifiedChains); i++ {
		if score := verifiedChainOverlap(presented, verifiedChains[i]); score > bestScore {
			bestScore = score
			bestIdx = i
		}
	}

	return verifiedChains[bestIdx]
}

func verifiedChainOverlap(presented, verified []*x509.Certificate) int {
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

func verifyPKIXChain(
	chain []*x509.Certificate,
	trusted TrustedRoots,
	now time.Time,
) ([]*x509.Certificate, error) {
	pool := trusted.Pool
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
		var unknownAuthority x509.UnknownAuthorityError
		if errors.As(err, &unknownAuthority) {
			anchorHint := "supplied root certificate(s)"
			if trusted.Pool == nil {
				anchorHint = "trusted roots (system CAs)"
			}

			return nil, fmt.Errorf(
				"x5chain verification failed: chain does not anchor to %s: %w",
				anchorHint,
				err,
			)
		}

		return nil, fmt.Errorf("x5chain verification failed: %w", err)
	}
	if len(verifiedChains) == 0 {
		return nil, fmt.Errorf("x5chain verification failed: no verified chain")
	}

	return selectVerifiedChain(chain, verifiedChains), nil
}

// VerifyX509ChainTrust validates the embedded x5chain and CoRIM COSE signature.
//
// Call [SignedCorim.FromCOSE] first so SigningCert, IntermediateCerts, and the
// COSE message needed for signature verification are populated from the
// signed-corim bytes. For key-based verification without PKIX trust, use
// [SignedCorim.Verify] instead.
//
// Supply trust material via [TrustedRoots]; see [TrustedRootPool] for file-based
// loading. A single-certificate x5chain (leaf only) is valid; intermediates
// are optional and taken from IntermediateCerts when present.
//
// Error prefix convention:
//   - "x5chain verification failed:" - PKIX path validation only
//   - "x5chain:" - signing-cert policy, CRL, or COSE signature checks
//
// Verification order:
//  1. validateSigningCertificate on the leaf signing cert
//  2. PKIX path validation with ExtKeyUsageAny against trusted.Pool
//  3. CRL checks on the PKIX-verified chain (when trusted.CRLs is non-empty)
//  4. COSE signature verification using the PKIX-verified leaf public key
func (o *SignedCorim) VerifyX509ChainTrust(trusted TrustedRoots) error {
	if o.SigningCert == nil {
		return errors.New("x5chain: header not set in CoRIM")
	}

	chain := make([]*x509.Certificate, 0, 1+len(o.IntermediateCerts))
	chain = append(chain, o.SigningCert)
	chain = append(chain, o.IntermediateCerts...)

	now := trusted.CurrentTime
	if now.IsZero() {
		now = time.Now()
	}

	if err := validateSigningCertificate(o.SigningCert); err != nil {
		return err
	}

	verifiedChain, err := verifyPKIXChain(chain, trusted, now)
	if err != nil {
		return err
	}

	if err := checkRevocation(verifiedChain, trusted.CRLs, now); err != nil {
		return err
	}

	if err := o.Verify(verifiedChain[0].PublicKey); err != nil {
		return fmt.Errorf("x5chain: COSE signature verification failed: %w", err)
	}

	return nil
}

// ParseX509CertificateDEROrPEM parses certificate bytes. When input is PEM,
// only the first block is used; any following blocks are ignored. PEM blocks
// must have type CERTIFICATE; otherwise input is treated as DER.
func ParseX509CertificateDEROrPEM(data []byte) (*x509.Certificate, error) {
	if block, _ := pem.Decode(data); block != nil {
		if block.Type != "CERTIFICATE" {
			return nil, fmt.Errorf("invalid PEM block type %q", block.Type)
		}

		data = block.Bytes
	}

	cert, err := x509.ParseCertificate(data)
	if err != nil {
		return nil, fmt.Errorf("parsing certificate: %w", err)
	}

	return cert, nil
}

// ParseRevocationListDEROrPEM parses CRL bytes. When input is PEM, only the
// first block is used; any following blocks are ignored. PEM blocks must have
// type X509 CRL; otherwise input is treated as DER.
func ParseRevocationListDEROrPEM(data []byte) (*x509.RevocationList, error) {
	if block, _ := pem.Decode(data); block != nil {
		if block.Type != "X509 CRL" {
			return nil, fmt.Errorf("invalid PEM block type %q", block.Type)
		}

		data = block.Bytes
	}

	crl, err := x509.ParseRevocationList(data)
	if err != nil {
		return nil, fmt.Errorf("parsing CRL: %w", err)
	}

	return crl, nil
}

// TrustedRootPool loads trust material for x5chain validation.
//
// readFile is injectable (typically os.ReadFile) for testing. rootPaths and
// crlPaths are read in order; the first read or parse error fails fast and
// returns a zero TrustedRoots. Pool is always non-nil on success. Duplicate
// root certificates (same DER) in rootPaths are added once. When
// includeSystemRoots is true, system root CAs are merged with certificates from
// rootPaths; when false, Pool contains only rootPaths. For CLI-style verification
// with no --root files, pass includeSystemRoots true and an empty rootPaths
// slice (do not use a zero TrustedRoots with Pool == nil). CurrentTime is left
// zero (verify time defaults to time.Now()).
func TrustedRootPool(
	readFile func(string) ([]byte, error),
	rootPaths, crlPaths []string,
	includeSystemRoots bool,
) (TrustedRoots, error) {
	var (
		pool *x509.CertPool
		err  error
	)

	if includeSystemRoots {
		pool, err = newSystemCertPool()
		if err != nil {
			return TrustedRoots{}, err
		}
	} else {
		pool = x509.NewCertPool()
	}

	trusted := TrustedRoots{
		Pool: pool,
		CRLs: make([]*x509.RevocationList, 0, len(crlPaths)),
	}

	addedRoots := make(map[string]struct{})

	for _, path := range rootPaths {
		data, err := readFile(path)
		if err != nil {
			return TrustedRoots{}, fmt.Errorf("loading root certificate from %s: %w", path, err)
		}

		cert, err := ParseX509CertificateDEROrPEM(data)
		if err != nil {
			return TrustedRoots{}, fmt.Errorf("parsing root certificate from %s: %w", path, err)
		}

		if _, seen := addedRoots[string(cert.Raw)]; seen {
			continue
		}

		addedRoots[string(cert.Raw)] = struct{}{}
		trusted.Pool.AddCert(cert)
	}

	for _, path := range crlPaths {
		data, err := readFile(path)
		if err != nil {
			return TrustedRoots{}, fmt.Errorf("loading CRL from %s: %w", path, err)
		}

		crl, err := ParseRevocationListDEROrPEM(data)
		if err != nil {
			return TrustedRoots{}, fmt.Errorf("parsing CRL from %s: %w", path, err)
		}

		trusted.CRLs = append(trusted.CRLs, crl)
	}

	return trusted, nil
}
