// Copyright 2021-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"errors"
	"fmt"
	"strings"

	"github.com/veraison/corim/extensions"
	cose "github.com/veraison/go-cose"
)

var (
	ContentType          = "application/rim+cbor"
	NoExternalData       = []byte("")
	HeaderLabelCorimMeta = int64(8)

	errNoSign1Message = errors.New("no Sign1 message found")

	// Sentinel errors from go-cose, re-exported so callers can match corim API
	// failures with errors.Is without naming the cose package at the call site.
	ErrX5ChainNoTrust    = cose.ErrX5ChainNoTrust
	ErrX5ChainCRLMissing = cose.ErrX5ChainCRLMissing
	ErrX5ChainRevoked    = cose.ErrX5ChainRevoked
	ErrX5ChainSignature  = cose.ErrX5ChainSignature
)

// SignedCorim encodes a signed-corim message (i.e., a COSE Sign1 wrapped CoRIM)
// with signature and verification methods
type SignedCorim struct {
	UnsignedCorim     UnsignedCorim
	Meta              Meta
	KeyID             []byte
	SigningCert       *x509.Certificate
	IntermediateCerts []*x509.Certificate
	message           *cose.Sign1Message
}

// NewSignedCorim instantiates an empty SignedCorim
func NewSignedCorim() *SignedCorim {
	return &SignedCorim{}
}

func (o *SignedCorim) RegisterExtensions(exts extensions.Map) error {
	unsignedExts := extensions.NewMap()

	for p, v := range exts {
		switch p {
		case ExtSigner:
			signerExts := extensions.NewMap().Add(ExtSigner, v)
			if err := o.Meta.RegisterExtensions(signerExts); err != nil {
				return err
			}
		default:
			unsignedExts.Add(p, v)
		}
	}

	return o.UnsignedCorim.RegisterExtensions(unsignedExts)
}

func (o *SignedCorim) processHdrs() error {
	var hdr = o.message.Headers

	if hdr.Protected == nil {
		return errors.New("missing mandatory protected header")
	}

	if v, ok := hdr.Protected[cose.HeaderLabelContentType]; ok {
		if v != ContentType {
			return fmt.Errorf("expecting content type %q, got %q instead", ContentType, v)
		}
	} else {
		return errors.New("missing mandatory content type")
	}

	if v, ok := hdr.Protected[cose.HeaderLabelKeyID]; ok {
		switch t := v.(type) {
		case []byte:
			o.KeyID = t
		default:
			return fmt.Errorf("kid: expected a []byte but got %v (%T)", t, t)
		}
	}

	if v, ok := hdr.Protected[HeaderLabelCorimMeta]; ok {
		if err := o.extractMeta(v); err != nil {
			return err
		}
	} else {
		return errors.New("missing mandatory corim.meta")
	}

	// Process optional x5chain
	if v, ok := hdr.Protected[cose.HeaderLabelX5Chain]; ok {
		if err := o.extractX5Chain(v); err != nil {
			return err
		}
	}

	return nil
}

func (o *SignedCorim) extractMeta(v interface{}) error {
	metaCBOR, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("expecting CBOR-encoded CoRIM Meta, got %T instead", v)
	}

	var meta Meta

	err := meta.FromCBOR(metaCBOR)
	if err != nil {
		return fmt.Errorf("unable to decode CoRIM Meta: %w", err)
	}

	o.Meta = meta

	return nil
}

// extractX5Chain decodes x5chain via cose.ParseX5Chain and sets SigningCert /
// IntermediateCerts only on success.
func (o *SignedCorim) extractX5Chain(x5chain interface{}) error {
	chain, err := cose.ParseX5Chain(x5chain)
	if err != nil {
		return err
	}

	o.SigningCert = chain.Leaf
	o.IntermediateCerts = chain.Intermediates

	return nil
}

// FromCOSE decodes a signed-corim message and validates the embedded
// unsigned-corim (including [UnsignedCorim.Valid]) and corim-meta. On success,
// UnsignedCorim and Meta are populated, and the message can be checked with
// [SignedCorim.Verify] or [SignedCorim.VerifyWithX5Chain].
//
// If an x5chain header is present, it must be leaf-first with contiguous
// signatures (RFC 9360 §2); decoding fails for any other order even when the
// caller only uses [SignedCorim.Verify]. Trust is not checked; use
// [SignedCorim.VerifyWithX5Chain] for that.
func (o *SignedCorim) FromCOSE(buf []byte) error {
	o.message = cose.NewSign1Message()
	o.SigningCert = nil
	o.IntermediateCerts = nil

	var err error
	// Roll back partial decode on any failure. Later steps must assign to err (not :=)
	// or this cleanup is skipped.
	defer func() {
		if err != nil {
			o.message = nil
			o.SigningCert = nil
			o.IntermediateCerts = nil
		}
	}()

	// If a tagged-corim-type-choice #6.500 of tagged-signed-corim #6.502, strip the prefix.
	// This is a remnant of an older draft of the specification before
	// https://github.com/ietf-rats-wg/draft-ietf-rats-corim/pull/337
	corimTypeChoice := []byte("\xd9\x01\xf4\xd9\x01\xf6")
	buf, _ = bytes.CutPrefix(buf, corimTypeChoice)

	if err = o.message.UnmarshalCBOR(buf); err != nil {
		return fmt.Errorf("failed CBOR decoding for COSE-Sign1 signed CoRIM: %w", err)
	}

	if err = o.processHdrs(); err != nil {
		return fmt.Errorf("processing COSE headers: %w", err)
	}

	if err = o.UnsignedCorim.FromCBOR(o.message.Payload); err != nil {
		return fmt.Errorf("failed CBOR decoding of unsigned CoRIM: %w", err)
	}

	if err = o.UnsignedCorim.Valid(); err != nil {
		return fmt.Errorf("failed validation of unsigned CoRIM: %w", err)
	}

	return nil
}

// AddSigningCert adds a DER-encoded X.509 certificate to be included in the
// protected header of the COSE Sign1 message as the leaf certificate in X5Chain.
func (o *SignedCorim) AddSigningCert(der []byte) error {
	if der == nil {
		return errors.New("nil signing cert")
	}

	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return fmt.Errorf("invalid signing certificate: %w", err)
	}

	o.SigningCert = cert
	return nil
}

// AddIntermediateCerts adds DER-encoded X.509 certificates to the protected
// header x5chain. Certificates must be concatenated with no padding and be in
// chain order, starting with the signing certificate's issuer (RFC 9360 §2);
// [SignedCorim.Sign] rejects any other order.
func (o *SignedCorim) AddIntermediateCerts(der []byte) error {
	if len(der) == 0 {
		return errors.New("nil or empty intermediate certs")
	}

	certs, err := x509.ParseCertificates(der)
	if err != nil {
		return fmt.Errorf("invalid intermediate certificates: %w", err)
	}

	if len(certs) == 0 {
		return errors.New("no certificates found in intermediate cert data")
	}

	o.IntermediateCerts = certs
	return nil
}

// Sign returns the serialized signed-corim, signed by the supplied cose Signer.
// The target SignedCorim must have its UnsignedCorim field correctly populated.
func (o *SignedCorim) Sign(signer cose.Signer) ([]byte, error) {
	if signer == nil {
		return nil, errors.New("nil signer")
	}

	if err := o.UnsignedCorim.Valid(); err != nil {
		return nil, fmt.Errorf("failed validation of unsigned CoRIM: %w", err)
	}

	// Build into a local message so that a failed Sign leaves o.message unchanged.
	msg := cose.NewSign1Message()

	var err error
	msg.Payload, err = o.UnsignedCorim.ToCBOR()
	if err != nil {
		return nil, fmt.Errorf("failed CBOR encoding of unsigned CoRIM: %w", err)
	}

	metaCBOR, err := o.Meta.ToCBOR()
	if err != nil {
		return nil, fmt.Errorf("failed CBOR encoding of CoRIM Meta: %w", err)
	}

	alg := signer.Algorithm()

	if strings.Contains(alg.String(), "unknown algorithm value") {
		return nil, errors.New("signer has no algorithm")
	}

	msg.Headers.Protected.SetAlgorithm(alg)
	msg.Headers.Protected[cose.HeaderLabelContentType] = ContentType
	msg.Headers.Protected[HeaderLabelCorimMeta] = metaCBOR

	if o.KeyID != nil {
		msg.Headers.Protected[cose.HeaderLabelKeyID] = o.KeyID
	}

	if o.SigningCert != nil {
		chain := cose.X5Chain{Leaf: o.SigningCert, Intermediates: o.IntermediateCerts}
		if err = cose.SetX5Chain(msg.Headers.Protected, chain); err != nil {
			return nil, fmt.Errorf("setting x5chain: %w", err)
		}
	} else if o.IntermediateCerts != nil {
		return nil, errors.New("intermediate certificates supplied but no signing certificate")
	}

	err = msg.Sign(rand.Reader, NoExternalData, signer)
	if err != nil {
		return nil, fmt.Errorf("COSE Sign1 signature failed: %w", err)
	}

	wrap, err := msg.MarshalCBOR()
	if err != nil {
		return nil, fmt.Errorf("signed-corim marshaling failed: %w", err)
	}

	o.message = msg

	return wrap, nil
}

// Verify verifies the signature of the target SignedCorim object using the
// supplied public key
func (o *SignedCorim) Verify(pk crypto.PublicKey) error {
	if o.message == nil {
		return errNoSign1Message
	}

	protected := o.message.Headers.Protected

	alg, err := protected.Algorithm()
	if err != nil {
		return fmt.Errorf("unable to get verification algorithm: %w", err)
	}

	verifier, err := cose.NewVerifier(alg, pk)
	if err != nil {
		return fmt.Errorf("unable to instantiate verifier: %w", err)
	}

	err = o.message.Verify(NoExternalData, verifier)
	if err != nil {
		return err
	}

	return nil
}

// VerifyWithX5Chain validates the protected-header x5chain and the CoRIM COSE
// signature via [cose.Sign1Message.VerifyWithX5Chain]. The signing certificate
// must not be a CA and, when keyUsage is present, must include digitalSignature.
// An x5chain in the unprotected header is rejected.
//
// PKIX validation accepts any extended key usage, so callers must authorize the
// signer themselves, for example by the Subject, SAN, or EKU of SigningCert.
// Verification uses the x5chain in the message header and does not update
// SigningCert or IntermediateCerts; do not change SigningCert before
// authorizing.
//
// Call it after [SignedCorim.FromCOSE] or a successful [SignedCorim.Sign].
// Configure trust with [cose.TrustAnchors] and optional
// [cose.X5ChainVerifyOptions] (see go-cose for revocation defaults). For
// external-key verify without PKIX, use [SignedCorim.Verify].
//
//nolint:gocritic // match cose.Sign1Message.VerifyWithX5Chain by-value anchors
func (o *SignedCorim) VerifyWithX5Chain(
	anchors cose.TrustAnchors,
	opts *cose.X5ChainVerifyOptions,
) error {
	if o.message == nil {
		return errNoSign1Message
	}

	_, err := o.message.VerifyWithX5Chain(NoExternalData, anchors, opts)

	return err
}
