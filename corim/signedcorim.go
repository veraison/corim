// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"crypto"
	"crypto/rand"
	"errors"
	"fmt"

	cose "github.com/veraison/go-cose"
)

var (
	ContentType = "application/rim+cbor"
)

// SignedCorim encodes a signed-corim type (i.e., a COSE Sign1 wrapped CoRIM)
// with signature and verification methods
type SignedCorim struct {
	UnsignedCorim UnsignedCorim
	message       *cose.Sign1Message
}

func (o *SignedCorim) processHdrs() error {
	var hdr = o.message.Headers

	if hdr.Protected == nil {
		return errors.New("missing mandatory protected header")
	}

	v, ok := hdr.Protected[3]
	if !ok {
		return errors.New("missing mandatory content type")
	}

	if v != ContentType {
		return fmt.Errorf("expecting content type %q, got %q instead", ContentType, v)
	}

	// TODO(tho) key id is apparently mandatory, which doesn't look right.
	// TODO(tho) Check with the CoRIM design team.
	// See https://github.com/veraison/corim/issues/14

	return nil
}

// FromCOSE decodes and effects syntactic validation on the supplied
// signed-corim, including the embedded unsigned-corim.
// On success, the unsigned-corim-map is made available via the UnsignedCorim
// field.
func (o *SignedCorim) FromCOSE(buf []byte) error {
	o.message = cose.NewSign1Message()

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
func (o *SignedCorim) Sign(signer *cose.Signer) ([]byte, error) {
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

	alg := signer.GetAlg()
	if alg == nil {
		return nil, errors.New("signer has no algorithm")
	}

	o.message.Headers.Protected[1] = alg.Value
	o.message.Headers.Protected[3] = ContentType

	err = o.message.Sign(rand.Reader, []byte(""), *signer)
	if err != nil {
		return nil, fmt.Errorf("COSE Sign1 signature failed: %w", err)
	}

	wrap, err := cose.Marshal(o.message)
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

	alg, err := cose.GetAlg(o.message.Headers)
	if err != nil {
		return fmt.Errorf("unable to get verification algorithm: %w", err)
	}

	verifier := cose.Verifier{
		Alg:       alg,
		PublicKey: pk,
	}

	err = o.message.Verify([]byte(""), verifier)
	if err != nil {
		return err
	}

	return nil
}
