// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package corim

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"errors"
	"fmt"
	"reflect"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
	cose "github.com/veraison/go-cose"
)

type Signer struct {
	Name string           `cbor:"0,keyasint" json:"name"`
	URI  *comid.TaggedURI `cbor:"1,keyasint,omitempty" json:"uri,omitempty"`

	Extensions
}

func NewSigner() *Signer {
	return &Signer{}
}

// RegisterExtensions registers a struct as a collections of extensions
func (o *Signer) RegisterExtensions(exts extensions.Map) error {
	for p, v := range exts {
		switch p {
		case ExtSigner:
			o.Extensions.Register(v)
		default:
			return fmt.Errorf("%w: %q", extensions.ErrUnexpectedPoint, p)
		}
	}

	return nil
}

// GetExtensions returns previously registered extension
func (o *Signer) GetExtensions() extensions.IMapValue {
	return o.Extensions.IMapValue
}

// SetName sets the target Signer's name to the supplied value
func (o *Signer) SetName(name string) *Signer {
	if o != nil {
		if name == "" {
			return nil
		}
		o.Name = name
	}
	return o
}

// SetURI sets the target Signer's URI to the supplied value
func (o *Signer) SetURI(uri string) *Signer {
	if o != nil {
		if uri == "" {
			return nil
		}

		taggedURI, err := comid.String2URI(&uri)
		if err != nil {
			return nil
		}

		o.URI = taggedURI
	}
	return o
}

// Valid checks the validity of individual fields within Signer
func (o Signer) Valid() error {
	if o.Name == "" {
		return errors.New("empty name")
	}

	if o.URI != nil {
		if err := comid.IsAbsoluteURI(string(*o.URI)); err != nil {
			return fmt.Errorf("invalid URI: %w", err)
		}
	}

	return o.Extensions.validSigner(&o)
}

// UnmarshalCBOR deserializes from CBOR
func (o *Signer) UnmarshalCBOR(data []byte) error {
	return encoding.PopulateStructFromCBOR(dm, data, o)
}

// MarshalCBOR serializes to CBOR
func (o Signer) MarshalCBOR() ([]byte, error) {
	return encoding.SerializeStructToCBOR(em, o)
}

// UnmarshalJSON deserializes from JSON
func (o *Signer) UnmarshalJSON(data []byte) error {
	return encoding.PopulateStructFromJSON(data, o)
}

// MarshalJSON serializes to JSON
func (o Signer) MarshalJSON() ([]byte, error) {
	return encoding.SerializeStructToJSON(o)
}

const noAlg = cose.Algorithm(-65537)

func getAlgAndKeyFromJWK(j []byte) (cose.Algorithm, crypto.Signer, error) {
	var (
		err error
		k   jwk.Key
		crv elliptic.Curve
		alg cose.Algorithm
	)

	k, err = jwk.ParseKey(j)
	if err != nil {
		return noAlg, nil, err
	}

	var key crypto.Signer

	err = k.Raw(&key)
	if err != nil {
		return noAlg, nil, err
	}

	switch v := key.(type) {
	case *ecdsa.PrivateKey:
		alg = ellipticCurveToAlg(v.Curve)
		if alg == noAlg {
			return noAlg, nil, fmt.Errorf("unknown elliptic curve %v", crv)
		}
	case ed25519.PrivateKey:
		alg = cose.AlgorithmEd25519
	case *rsa.PrivateKey:
		alg = rsaJWKToAlg(k)
		if alg == noAlg {
			return noAlg, nil, fmt.Errorf("unknown RSA algorithm %q", k.Algorithm().String())
		}
	default:
		return noAlg, nil, fmt.Errorf("unknown private key type %v", reflect.TypeOf(key))
	}

	return alg, key, nil
}

func getKidFromJWK(j []byte) ([]byte, error) {
	k, err := jwk.ParseKey(j)
	if err != nil {
		return nil, err
	}

	if k.KeyID() != "" {
		return []byte(k.KeyID()), nil
	}

	// Generate a key ID from the JWK Thumbprint if none exist
	// See https://datatracker.ietf.org/doc/html/rfc7638
	kid, err := k.Thumbprint(crypto.SHA256)
	if err != nil {
		return nil, err
	}
	return kid, nil
}

func ellipticCurveToAlg(c elliptic.Curve) cose.Algorithm {
	switch c {
	case elliptic.P256():
		return cose.AlgorithmES256
	case elliptic.P384():
		return cose.AlgorithmES384
	case elliptic.P521():
		return cose.AlgorithmES512
	default:
		return noAlg
	}
}

func rsaJWKToAlg(k jwk.Key) cose.Algorithm {
	switch k.Algorithm().String() {
	case "PS256":
		return cose.AlgorithmPS256
	case "PS384":
		return cose.AlgorithmPS384
	case "PS512":
		return cose.AlgorithmPS512
	default:
		return noAlg
	}
}

func NewSignerFromJWK(j []byte) (cose.Signer, error) {
	alg, key, err := getAlgAndKeyFromJWK(j)
	if err != nil {
		return nil, err
	}

	return cose.NewSigner(alg, key)
}

func NewPublicKeyFromJWK(j []byte) (crypto.PublicKey, error) {
	_, key, err := getAlgAndKeyFromJWK(j)
	if err != nil {
		return nil, err
	}

	return key.Public(), nil
}
