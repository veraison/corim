// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"bytes"
	"crypto"
	"crypto/rand"
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
	UnsignedCorim UnsignedCorim
	Meta          Meta
	message       *cose.Sign1Message
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

	v, ok := hdr.Protected[cose.HeaderLabelAlgorithm]
	if !ok {
		return errors.New("missing mandatory algorithm")
	}

	// TODO: make this consistent, either int64 or cose.Algorithm
	// cose.Algorithm is an alias to int64 defined in veraison/go-cose
	switch v.(type) {
	case int64:
	case cose.Algorithm:
	default:
		return fmt.Errorf("expecting integer CoRIM Algorithm, got %T instead", v)
	}

	v, ok = hdr.Protected[cose.HeaderLabelContentType]
	if !ok {
		return errors.New("missing mandatory content type")
	}

	if v != ContentType {
		return fmt.Errorf("expecting content type %q, got %q instead", ContentType, v)
	}

	v, ok = hdr.Protected[cose.HeaderLabelKeyID]
	if !ok {
		return errors.New("missing mandatory key id")
	}

	_, ok = v.([]byte)
	if !ok {
		return fmt.Errorf("expecting byte string CoRIM Key ID, got %T instead", v)
	}

	v, ok = hdr.Protected[HeaderLabelCorimMeta]
	if !ok {
		return errors.New("missing mandatory corim.meta")
	}

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

// Sign returns the serialized signed-corim, signed by the supplied cose Signer.
// The target SignedCorim must have its UnsignedCorim field correctly
// populated.
func (o *SignedCorim) Sign(signer cose.Signer, kid []byte) ([]byte, error) {
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
	o.message.Headers.Protected[cose.HeaderLabelKeyID] = kid
	o.message.Headers.Protected[HeaderLabelCorimMeta] = metaCBOR

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
