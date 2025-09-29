// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

// Package coserv provides an implementation of draft-howard-rats-coserv
package coserv

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/eat"
	"github.com/veraison/go-cose"
)

// Coserv is the internal representation of a CoSERV data item
type Coserv struct {
	Profile eat.Profile `cbor:"0,keyasint"`
	Query   Query       `cbor:"1,keyasint"`
	Results *ResultSet  `cbor:"2,keyasint,omitempty"`
}

// NewCoserv creates a new Coserv instance.
// An error is returned if the supplied profile or query are invalid
func NewCoserv(profile string, query Query) (*Coserv, error) {
	p, err := eat.NewProfile(profile)
	if err != nil {
		return nil, fmt.Errorf("invalid profile: %w", err)
	}

	if err := query.Valid(); err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}

	return &Coserv{
		Profile: *p,
		Query:   query,
	}, nil
}

// AddResults add the result set to the Coserv target after validating it
func (o *Coserv) AddResults(v ResultSet) error {
	if err := v.Valid(); err != nil {
		return fmt.Errorf("invalid result set: %w", err)
	}

	o.Results = &v

	return nil
}

// ToEDN encodes the target Coserv to CBOR Extended Diagnostic Notation (EDN)
func (o Coserv) ToEDN() (string, error) { // nolint:gocritic
	b, err := o.ToCBOR()
	if err != nil {
		return "", fmt.Errorf("failed encoding Coserv object: %w", err)
	}
	return cbor.Diagnose(b)
}

// Valid ensures that the Coserv target is correctly populated
func (o Coserv) Valid() error { // nolint:gocritic
	// TBC:
	// * artifact type mismatch should be caught on decoding
	// * ditto for profile syntax errors
	if err := o.Query.Valid(); err != nil {
		return fmt.Errorf("invalid query: %w", err)
	}
	return nil
}

// ToCBOR validates and serializes to CBOR the target Coserv
// An error is returned if either validation or encoding of the Coserv target fails
func (o Coserv) ToCBOR() ([]byte, error) { // nolint:gocritic
	if err := o.Valid(); err != nil {
		return nil, fmt.Errorf("validating Coserv: %w", err)
	}

	opts := cbor.CoreDetEncOptions()
	opts.Time = cbor.TimeRFC3339
	opts.TimeTag = 1
	em, err := opts.EncMode()
	if err != nil {
		return nil, fmt.Errorf("CBOR encoding setup failed: %w", err)
	}

	data, err := em.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("encoding Coserv to CBOR: %w", err)
	}

	return data, nil
}

// FromCBOR deserializes from CBOR into the target Coserv
// An error is returned if either decoding or validation of the CoSERV payload fails
func (o *Coserv) FromCBOR(data []byte) error {
	if err := cbor.Unmarshal(data, o); err != nil {
		return fmt.Errorf("decoding CoSERV from CBOR: %w", err)
	}

	if err := o.Valid(); err != nil {
		return fmt.Errorf("validating CoSERV: %w", err)
	}

	return nil
}

// ToBase64Url validates and serializes to base64url the target Coserv
// An error is returned if either validation or encoding of the Coserv target fails
func (o Coserv) ToBase64Url() (string, error) { // nolint:gocritic
	data, err := o.ToCBOR()
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

// FromBase64Url deserializes from base64url-encoded into the target Coserv
// An error is returned if either decoding or validation of the CoSERV payload fails
func (o *Coserv) FromBase64Url(s string) error {
	data, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("decoding CoSERV: %w", err)
	}

	return o.FromCBOR(data)
}

// Sign signs and serializes the target Coserv using the supplied go-cose Signer
func (o *Coserv) Sign(signer cose.Signer) ([]byte, error) {
	msg := cose.NewSignMessage()

	msg.Headers.Protected[cose.HeaderLabelAlgorithm] = signer.Algorithm()
	msg.Headers.Protected[cose.HeaderLabelContentType] = "application/coserv+cbor"

	payload, err := o.ToCBOR()
	if err != nil {
		return nil, err
	}

	return cose.Sign1(rand.Reader, signer, msg.Headers, payload, nil)
}

// Verify verifies the signature of a signed Coserv object using the supplied go-cose Verifier
func (o *Coserv) Verify(verifier cose.Verifier, data []byte) error {
	var msg cose.Sign1Message
	if err := msg.UnmarshalCBOR(data); err != nil {
		return fmt.Errorf("CBOR decoding signed-coserv: %w", err)
	}

	if v, ok := msg.Headers.Protected[cose.HeaderLabelContentType]; ok {
		if v != "application/coserv+cbor" {
			return fmt.Errorf("unexpected content type in signed-coserv: %v", v)
		}
	} else {
		return fmt.Errorf("missing mandatory cty parameter in signed-coserv protected headers")
	}

	if _, ok := msg.Headers.Protected[cose.HeaderLabelAlgorithm]; !ok {
		return fmt.Errorf("missing mandatory alg parameter in signed-coserv protected headers")
	}

	if err := msg.Verify(nil, verifier); err != nil {
		return fmt.Errorf("signed-coserv signature verification failed: %w", err)
	}

	if err := o.FromCBOR(msg.Payload); err != nil {
		return fmt.Errorf("CBOR decoding signed-coserv payload: %w", err)
	}

	return nil
}
