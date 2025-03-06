// Copyright 2021-2024 Contributors to the Veraison project.
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
)

// SignedCorim encodes a signed-corim message (i.e., a COSE Sign1 wrapped CoRIM)
// with signature and verification methods
type SignedCorim struct {
	UnsignedCorim     UnsignedCorim
	Meta              Meta
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

	// TODO(tho) key id is apparently mandatory, which doesn't look right.
	// TODO(tho) Check with the CoRIM design team.
	// See https://github.com/ietf-rats-wg/draft-ietf-rats-corim/issues/363

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

func (o *SignedCorim) extractX5Chain(x5chain interface{}) error {
	var buf bytes.Buffer

	switch t := x5chain.(type) {
	case []interface{}:
		for i, elem := range t {
			cert, ok := elem.([]byte)
			if !ok {
				return fmt.Errorf("accessing x5chain[%d]: got %T, want []byte", i, elem)
			}

			switch i {
			case 0:
				if err := o.AddSigningCert(cert); err != nil {
					return fmt.Errorf("decoding x5chain: %w", err)
				}
			default:
				buf.Write(cert)
			}
		}

		if buf.Len() > 0 {
			if err := o.AddIntermediateCerts(buf.Bytes()); err != nil {
				return fmt.Errorf("decoding x5chain: %w", err)
			}
		}
	case []byte:
		if err := o.AddSigningCert(t); err != nil {
			return fmt.Errorf("decoding x5chain: %w", err)
		}
	default:
		return fmt.Errorf("decoding x5chain: got %T, want []interface{} or []byte", t)
	}

	return nil
}

// FromCOSE decodes and effects syntactic validation on the supplied
// signed-corim message, including the embedded unsigned-corim and corim-meta.
// On success, the unsigned-corim-map is made available via the UnsignedCorim
// field while the corim-meta-map is decoded into the Meta field.
func (o *SignedCorim) FromCOSE(buf []byte) error {
	o.message = cose.NewSign1Message()

	// If a tagged-corim-type-choice #6.500 of tagged-signed-corim #6.502, strip the prefix.
	// This is a remnant of an older draft of the specification before
	// https://github.com/ietf-rats-wg/draft-ietf-rats-corim/pull/337
	corimTypeChoice := []byte("\xd9\x01\xf4\xd9\x01\xf6")
	buf, _ = bytes.CutPrefix(buf, corimTypeChoice)

	if err := o.message.UnmarshalCBOR(buf); err != nil {
		return fmt.Errorf("failed CBOR decoding for COSE-Sign1 signed CoRIM: %w", err)
	}

	if err := o.processHdrs(); err != nil {
		return fmt.Errorf("processing COSE headers: %w", err)
	}

	if err := o.UnsignedCorim.FromCBOR(o.message.Payload); err != nil {
		return fmt.Errorf("failed CBOR decoding of unsigned CoRIM: %w", err)
	}

	if err := o.UnsignedCorim.Valid(); err != nil {
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

// AddIntermediateCerts adds DER-encoded X.509 certificates to be included in the protected
// header of the COSE Sign1 message as part of the X5Chain.
// The certificates must be concatenated with no intermediate padding, as per X.509 convention.
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

	o.message = cose.NewSign1Message()

	var err error
	o.message.Payload, err = o.UnsignedCorim.ToCBOR()
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

	o.message.Headers.Protected.SetAlgorithm(alg)
	o.message.Headers.Protected[cose.HeaderLabelContentType] = ContentType
	o.message.Headers.Protected[HeaderLabelCorimMeta] = metaCBOR

	if o.SigningCert != nil {
		// COSE_X509 = bstr / [ 2*certs: bstr ]
		//
		// handle alt (1): bstr
		if len(o.IntermediateCerts) == 0 {
			o.message.Headers.Protected[cose.HeaderLabelX5Chain] = o.SigningCert.Raw
		} else { // handle alt (2): [ 2*certs: bstr ]
			certChain := [][]byte{o.SigningCert.Raw}
			for _, cert := range o.IntermediateCerts {
				certChain = append(certChain, cert.Raw)
			}
			o.message.Headers.Protected[cose.HeaderLabelX5Chain] = certChain
		}
	}

	err = o.message.Sign(rand.Reader, NoExternalData, signer)
	if err != nil {
		return nil, fmt.Errorf("COSE Sign1 signature failed: %w", err)
	}

	wrap, err := o.message.MarshalCBOR()
	if err != nil {
		return nil, fmt.Errorf("signed-corim marshaling failed: %w", err)
	}

	return wrap, nil
}

// Verify verifies the signature of the target SignedCorim object using the
// supplied public key
func (o *SignedCorim) Verify(pk crypto.PublicKey) error {
	if o.message == nil {
		return errors.New("no Sign1 message found")
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
