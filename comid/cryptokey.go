// Copyright 2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/go-cose"
	"github.com/veraison/swid"
)

const (
	// PKIXBase64KeyType indicates a PEM-encoded SubjectPublicKeyInfo. See
	// https://www.rfc-editor.org/rfc/rfc7468#section-13
	PKIXBase64KeyType = "pkix-base64-key"
	// PKIXBase64CertType indicates a PEM-encoded X.509 public key
	// certificate. See https://www.rfc-editor.org/rfc/rfc7468#section-5
	PKIXBase64CertType = "pkix-base64-cert"
	// PKIXBase64CertPathType indicates a X.509 certificate chain created
	// by the concatenation of as many PEM encoded X.509 certificates as
	// needed. The certificates MUST be concatenated in order so that each
	// directly certifies the one preceding.
	PKIXBase64CertPathType = "pkix-base64-cert-path"
	// COSEKeyType represents a CBOR encoded COSE_Key or COSE_KeySet. See
	// https://www.rfc-editor.org/rfc/rfc9052#section-7
	COSEKeyType = "cose-key"
	// ThumbprintType represents a digest of a raw public key. The digest
	// value may be used to find the public key if contained in a lookup
	// table.
	ThumbprintType = "thumbprint"
	// CertThumbprintType represents a digest of a certificate. The digest
	// value may be used to find the certificate if contained in a lookup
	// table.
	CertThumbprintType = "cert-thumbprint"
	// CertPathThumbprintType represents a digest of a certification path.
	// The digest value may be used to find the certificate path if
	// contained in a lookup table.
	CertPathThumbprintType = "cert-path-thumbprint"
)

// CryptoKey is the struct implementing CoRIM crypto-key-type-choice. See
// https://www.ietf.org/archive/id/draft-ietf-rats-corim-02.html#name-crypto-keys
type CryptoKey struct {
	Value ICryptoKeyValue
}

// NewCryptoKey returns the pointer to a new CryptoKey of the specified type,
// constructed using the provided value k. The type of k depends on the
// specified crypto key type. For PKIX types, k must be a string. For COSE_Key,
// k must be a []byte. For thumbprint types, k must be a swid.HashEntry.
func NewCryptoKey(k any, typ string) (*CryptoKey, error) {
	switch typ {
	case PKIXBase64KeyType:
		v, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("value must be a string; found %T", k)
		}
		return NewPKIXBase64Key(v)
	case PKIXBase64CertType:
		v, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("value must be a string; found %T", k)
		}
		return NewPKIXBase64Cert(v)
	case PKIXBase64CertPathType:
		v, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("value must be a string; found %T", k)
		}
		return NewPKIXBase64CertPath(v)
	case COSEKeyType:
		v, ok := k.([]byte)
		if !ok {
			return nil, fmt.Errorf("value must be a []byte; found %T", k)
		}
		return NewCOSEKey(v)
	case ThumbprintType, CertThumbprintType, CertPathThumbprintType:
		v, ok := k.(swid.HashEntry)
		if !ok {
			return nil, fmt.Errorf("value must be a swid.HashEntry; found %T", k)
		}
		switch typ {
		case ThumbprintType:
			return NewThumbprint(v)
		case CertThumbprintType:
			return NewCertThumbprint(v)
		case CertPathThumbprintType:
			return NewCertPathThumbprint(v)
		default:
			// Should never here because of the the outer case clause
			panic(fmt.Sprintf("unexpected thumbprint type: %s", typ))
		}
	default:
		return nil, fmt.Errorf("unexpected CryptoKey type: %s", typ)
	}
}

// MustNewCryptoKey is the same as NewCryptoKey, but does not return an error,
// and panics if there is a problem.
func MustNewCryptoKey(k any, typ string) *CryptoKey {
	key, err := NewCryptoKey(k, typ)
	if err != nil {
		panic(err)
	}

	return key
}

// String returns the string representation of the CryptoKey.
func (o CryptoKey) String() string {
	return o.Value.String()
}

// Valid returns an error if validation of the CryptoKey fails, or nil if it
// succeeds.
func (o CryptoKey) Valid() error {
	return o.Value.Valid()
}

// PublicKey returns a crypto.PublicKey constructed from the CryptoKey's
// underlying value. This returns an error if the CryptoKey is one of the
// thumbprint types.
func (o CryptoKey) PublicKey() (crypto.PublicKey, error) {
	return o.Value.PublicKey()
}

// MarshalJSON returns a []byte containing the JSON representation of the
// CryptoKey.
func (o CryptoKey) MarshalJSON() ([]byte, error) {
	value := struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}{
		Value: o.Value.String(),
	}

	switch o.Value.(type) {
	case TaggedPKIXBase64Key:
		value.Type = PKIXBase64KeyType
	case TaggedPKIXBase64Cert:
		value.Type = PKIXBase64CertType
	case TaggedPKIXBase64CertPath:
		value.Type = PKIXBase64CertPathType
	case TaggedCOSEKey:
		value.Type = COSEKeyType
	case TaggedThumbprint:
		value.Type = ThumbprintType
	case TaggedCertThumbprint:
		value.Type = CertThumbprintType
	case TaggedCertPathThumbprint:
		value.Type = CertPathThumbprintType
	default:
		return nil, fmt.Errorf("unexpected ICryptoKeyValue type: %T", o.Value)
	}

	return json.Marshal(value)
}

// UnmarshalJSON populates the CryptoKey from the JSON representation inside
// the provided []byte.
func (o *CryptoKey) UnmarshalJSON(b []byte) error {
	var value struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	if value.Type == "" {
		return errors.New("key type not set")
	}

	switch value.Type {
	case PKIXBase64KeyType:
		o.Value = TaggedPKIXBase64Key(value.Value)
	case PKIXBase64CertType:
		o.Value = TaggedPKIXBase64Cert(value.Value)
	case PKIXBase64CertPathType:
		o.Value = TaggedPKIXBase64CertPath(value.Value)
	case COSEKeyType:
		data, err := base64.StdEncoding.DecodeString(value.Value)
		if err != nil {
			return fmt.Errorf("base64 decode error: %w", err)
		}
		o.Value = TaggedCOSEKey(data)
	case ThumbprintType, CertThumbprintType, CertPathThumbprintType:
		he, err := swid.ParseHashEntry(value.Value)
		if err != nil {
			return fmt.Errorf("swid.HashEntry decode error: %w", err)
		}
		switch value.Type {
		case ThumbprintType:
			o.Value = TaggedThumbprint{digest{he}}
		case CertThumbprintType:
			o.Value = TaggedCertThumbprint{digest{he}}
		case CertPathThumbprintType:
			o.Value = TaggedCertPathThumbprint{digest{he}}
		}
	default:
		return fmt.Errorf("unexpected ICryptoKeyValue type: %q", value.Type)
	}

	return o.Valid()
}

// MarshalCBOR returns a []byte containing the CBOR representation of the
// CryptoKey.
func (o CryptoKey) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.Value)
}

// UnmarshalCBOR populates the CryptoKey from the CBOR representation inside
// the provided []byte.
func (o *CryptoKey) UnmarshalCBOR(b []byte) error {
	return dm.Unmarshal(b, &o.Value)
}

// ICryptoKeyValue is the interface implemented by the concrete CryptoKey value
// types.
type ICryptoKeyValue interface {
	// String returns the string representation of the ICryptoKeyValue.
	String() string
	// Valid returns an error if validation of the ICryptoKeyValue fails,
	// or nil if it succeeds.
	Valid() error
	// PublicKey returns a crypto.PublicKey constructed from the
	// ICryptoKeyValue's underlying value. This returns an error if the
	// ICryptoKeyValue is one of the thumbprint types.
	PublicKey() (crypto.PublicKey, error)
}

// TaggedPKIXBase64Key is a PEM-encoded SubjectPublicKeyInfo. See
// https://www.rfc-editor.org/rfc/rfc7468#section-13
type TaggedPKIXBase64Key string

func NewPKIXBase64Key(s string) (*CryptoKey, error) {
	key := TaggedPKIXBase64Key(s)
	if err := key.Valid(); err != nil {
		return nil, err
	}
	return &CryptoKey{key}, nil
}

func MustNewPKIXBase64Key(s string) *CryptoKey {
	key, err := NewPKIXBase64Key(s)
	if err != nil {
		panic(err)
	}
	return key
}

func (o TaggedPKIXBase64Key) String() string {
	return string(o)
}

func (o TaggedPKIXBase64Key) Valid() error {
	_, err := o.PublicKey()
	return err
}

func (o TaggedPKIXBase64Key) PublicKey() (crypto.PublicKey, error) {
	if string(o) == "" {
		return nil, errors.New("key value not set")
	}

	block, rest := pem.Decode([]byte(o))
	if block == nil {
		return nil, errors.New("could not decode PEM block")
	}

	if len(rest) != 0 {
		return nil, errors.New("trailing data found after PEM block")
	}

	if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf(
			"unexpected PEM block type: %q, expected \"PUBLIC KEY\"",
			block.Type,
		)
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse public key: %w", err)
	}

	return key, nil
}

// TaggedPKIXBase64Cert is a PEM-encoded X.509 public key
// certificate. See https://www.rfc-editor.org/rfc/rfc7468#section-5
type TaggedPKIXBase64Cert string

func NewPKIXBase64Cert(s string) (*CryptoKey, error) {
	cert := TaggedPKIXBase64Cert(s)
	if err := cert.Valid(); err != nil {
		return nil, err
	}
	return &CryptoKey{cert}, nil
}

func MustNewPKIXBase64Cert(s string) *CryptoKey {
	cert, err := NewPKIXBase64Cert(s)
	if err != nil {
		panic(err)
	}
	return cert
}

func (o TaggedPKIXBase64Cert) String() string {
	return string(o)
}

func (o TaggedPKIXBase64Cert) Valid() error {
	_, err := o.cert()
	return err
}

func (o TaggedPKIXBase64Cert) PublicKey() (crypto.PublicKey, error) {
	cert, err := o.cert()
	if err != nil {
		return nil, err
	}

	if cert.PublicKey == nil {
		return nil, errors.New("cert does not contain a crypto.PublicKey")
	}

	return cert.PublicKey, nil
}

func (o TaggedPKIXBase64Cert) cert() (*x509.Certificate, error) {
	if string(o) == "" {
		return nil, errors.New("cert value not set")
	}

	block, rest := pem.Decode([]byte(o))
	if block == nil {
		return nil, errors.New("could not decode PEM block")
	}

	if len(rest) != 0 {
		return nil, errors.New("trailing data found after PEM block")
	}

	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf(
			"unexpected PEM block type: %q, expected \"CERTIFICATE\"",
			block.Type,
		)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse x509 cert: %w", err)
	}

	return cert, nil
}

// TaggedPKIXBase64CertPath is a X.509 certificate chain created
// by the concatenation of as many PEM encoded X.509 certificates as
// needed. The certificates MUST be concatenated in order so that each
// directly certifies the one preceding.
type TaggedPKIXBase64CertPath string

func NewPKIXBase64CertPath(s string) (*CryptoKey, error) {
	cert := TaggedPKIXBase64CertPath(s)

	if err := cert.Valid(); err != nil {
		return nil, err
	}

	return &CryptoKey{cert}, nil
}

func MustNewPKIXBase64CertPath(s string) *CryptoKey {
	cert, err := NewPKIXBase64CertPath(s)

	if err != nil {
		panic(err)
	}

	return cert
}

func (o TaggedPKIXBase64CertPath) String() string {
	return string(o)
}

func (o TaggedPKIXBase64CertPath) Valid() error {
	_, err := o.certPath()
	return err
}

func (o TaggedPKIXBase64CertPath) PublicKey() (crypto.PublicKey, error) {
	certs, err := o.certPath()
	if err != nil {
		return nil, err
	}

	if len(certs) == 0 {
		return nil, errors.New("empty cert path")
	}

	if certs[0].PublicKey == nil {
		return nil, errors.New("leaf cert does not contain a crypto.PublicKey")
	}

	return certs[0].PublicKey, nil
}

func (o TaggedPKIXBase64CertPath) certPath() ([]*x509.Certificate, error) {
	if string(o) == "" {
		return nil, errors.New("cert value not set")
	}

	var certs []*x509.Certificate
	var block *pem.Block
	var rest []byte
	rest = []byte(o)
	i := 0
	for {
		if len(rest) == 0 {
			break
		}

		block, rest = pem.Decode(rest)
		if block == nil {
			return nil, fmt.Errorf("could not decode PEM block %d", i)
		}

		if block.Type != "CERTIFICATE" {
			return nil, fmt.Errorf(
				"unexpected type for PEM block %d: %q, expected \"CERTIFICATE\"",
				i, block.Type,
			)
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf(
				"could not parse x509 cert in PEM block %d: %w",
				i, err,
			)
		}

		certs = append(certs, cert)

		i++
	}

	return certs, nil
}

// TaggedCOSEKey is a CBOR encoded COSE_Key or COSE_KeySet. See
// https://www.rfc-editor.org/rfc/rfc9052#section-7
type TaggedCOSEKey []byte

func NewCOSEKey(b []byte) (*CryptoKey, error) {
	key := TaggedCOSEKey(b)

	if err := key.Valid(); err != nil {
		return nil, err
	}

	return &CryptoKey{key}, nil
}

func MustNewCOSEKey(b []byte) *CryptoKey {
	key, err := NewCOSEKey(b)

	if err != nil {
		panic(err)
	}

	return key
}

func (o TaggedCOSEKey) String() string {
	return base64.StdEncoding.EncodeToString(o)
}

func (o TaggedCOSEKey) Valid() error {
	if len(o) == 0 {
		return errors.New("empty COSE_Key bytes")
	}

	var err error

	// CBOR Major type 4 == array == COSE_KeySet. Key sets are currently
	// not supported by go-cose library.
	if ((o[0] & 0xe0) >> 5) == 4 {
		_, err = o.coseKeySet()
	} else {
		_, err = o.coseKey()
	}
	return err
}

func (o TaggedCOSEKey) PublicKey() (crypto.PublicKey, error) {
	if len(o) == 0 {
		return nil, errors.New("empty COSE_Key value")
	}

	// CBOR Major type 4 == array == COSE_KeySet. Key sets are currently
	// not supported by go-cose library.
	if ((o[0] & 0xe0) >> 5) == 4 {
		keySet, err := o.coseKeySet()
		if err != nil {
			return nil, err
		}

		if len(keySet) == 0 {
			return nil, errors.New("empty COSE_KeySet")
		} else if len(keySet) > 1 {
			return nil, errors.New("COSE_KeySet contains more than one key")
		}

		return keySet[0].PublicKey()
	}

	coseKey, err := o.coseKey()
	if err != nil {
		return nil, err
	}

	return coseKey.PublicKey()
}

func (o TaggedCOSEKey) MarshalCBOR() ([]byte, error) {
	var buf bytes.Buffer

	// encodeMarshalerType in github.com/fxamacker/cbor/v2 does not look up
	// assocated Tags, so we have to write them ourselves.
	if _, err := buf.Write([]byte{0xd9, 0x02, 0x2e}); err != nil { // tag 558
		return nil, err
	}

	if _, err := buf.Write(o); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (o *TaggedCOSEKey) UnmarshalCBOR(b []byte) error {
	// the first 3 bytes are the tag
	if !bytes.Equal([]byte{0xd9, 0x02, 0x2e}, b[:3]) {
		return errors.New("did not see CBOR tag 588 at the beginning of COSE_Key structure")
	}

	*o = b[3:]
	return nil
}

func (o TaggedCOSEKey) coseKey() (*cose.Key, error) {
	coseKey := new(cose.Key)

	if err := coseKey.UnmarshalCBOR(o); err != nil {
		return nil, err
	}

	return coseKey, nil
}

func (o TaggedCOSEKey) coseKeySet() ([]*cose.Key, error) {
	var keySet []*cose.Key

	if err := cbor.Unmarshal(o, &keySet); err != nil {
		return nil, err
	}

	return keySet, nil
}

type digest struct {
	swid.HashEntry
}

func (o digest) String() string {
	return o.HashEntry.String()
}

func (o digest) Valid() error {
	return swid.ValidHashEntry(o.HashAlgID, o.HashValue)
}

func (o digest) PublicKey() (crypto.PublicKey, error) {
	return nil, errors.New("cannot get PublicKey from a digest")
}

// ThumbprintTypeTaggedThumbprint represents a digest of a raw public key. The
// digest value may be used to find the public key if contained in a lookup
// table.
type TaggedThumbprint struct {
	digest
}

func NewThumbprint(he swid.HashEntry) (*CryptoKey, error) {
	key := &CryptoKey{TaggedThumbprint{digest{he}}}

	if err := key.Valid(); err != nil {
		return nil, err
	}

	return key, nil
}

func MustNewThumbprint(he swid.HashEntry) *CryptoKey {
	key, err := NewThumbprint(he)

	if err != nil {
		panic(err)
	}

	return key
}

// TaggedCertThumbprint represents a digest of a certificate. The digest value
// may be used to find the certificate if contained in a lookup table.
type TaggedCertThumbprint struct {
	digest
}

func NewCertThumbprint(he swid.HashEntry) (*CryptoKey, error) {
	key := &CryptoKey{TaggedCertThumbprint{digest{he}}}

	if err := key.Valid(); err != nil {
		return nil, err
	}

	return key, nil
}

func MustNewCertThumbprint(he swid.HashEntry) *CryptoKey {
	key, err := NewCertThumbprint(he)

	if err != nil {
		panic(err)
	}

	return key
}

// TaggedCertPathThumbprint represents a digest of a certification path. The
// digest value may be used to find the certificate path if contained in a
// lookup table.
type TaggedCertPathThumbprint struct {
	digest
}

func NewCertPathThumbprint(he swid.HashEntry) (*CryptoKey, error) {
	key := &CryptoKey{TaggedCertPathThumbprint{digest{he}}}

	if err := key.Valid(); err != nil {
		return nil, err
	}

	return key, nil
}

func MustNewCertPathThumbprint(he swid.HashEntry) *CryptoKey {
	key, err := NewCertPathThumbprint(he)

	if err != nil {
		panic(err)
	}

	return key
}
