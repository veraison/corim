// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "fmt"

// VerifKey stores the verification key material associated to a signing key.
// Key is - typically, but not necessarily - a public key.  Chain is an optional
// X.509 certificate chain corresponding to the public key in Key, encoded as an
// array of one or more base64-encoded DER PKIX certificates.  The certificate
// containing the public key in Key MUST be the first certificate.  This MAY be
// followed by additional certificates, with each subsequent certificate being
// the one used to certify the previous one.
type VerifKey struct {
	Key   string    `cbor:"0,keyasint" json:"key"`
	Chain *[]string `cbor:"1,keyasint,omitempty" json:"chain,omitempty"`
}

// NewVerifKey instantiates an empty VerifKey
func NewVerifKey() *VerifKey {
	return &VerifKey{}
}

// SetKey sets the Key in the target object to the supplied value
func (o *VerifKey) SetKey(key string) *VerifKey {
	if o != nil {
		o.Key = key
	}
	return o
}

// AddCert adds the supplied base64-encoded DER PKIX certificate in the target
// object
func (o *VerifKey) AddCert(cert string) *VerifKey {
	if o != nil {
		if o.Chain == nil {
			o.Chain = new([]string)
		}
		*o.Chain = append(*o.Chain, cert)
	}
	return o
}

func (o VerifKey) Valid() error {
	if o.Key == "" {
		return fmt.Errorf("verification key not set")
	}
	return nil
}
