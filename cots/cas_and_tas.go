// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cots

import (
	"encoding/json"
	"fmt"
)

type TaFormat int64

const (
	TaFormatCertificate TaFormat = iota
	TaFormatTrustAnchorInfo
	TaFormatSubjectPublicKeyInfo
)

var (
	roleToString = map[TaFormat]string{
		TaFormatCertificate:          "cert",
		TaFormatTrustAnchorInfo:      "ta",
		TaFormatSubjectPublicKeyInfo: "spki",
	}

	stringToRole = map[string]TaFormat{
		"cert": TaFormatCertificate,
		"ta":  TaFormatTrustAnchorInfo,
		"spki": TaFormatSubjectPublicKeyInfo,
	}
)

//type TrustAnchor struct {
//	Format TaFormat `cbor:"0,keyasint" json:"format"`
//	Data   []byte   `cbor:"1,keyasint,omitempty" json:"data,omitempty"`
//}

type TrustAnchor struct {
	_            struct{}     `cbor:",toarray"`
	Format TaFormat `json:"format"`
	Data   []byte   `json:"data"`
}

func NewTrustAnchor() *TrustAnchor {
	return &TrustAnchor{}
}

func (o *TrustAnchor) SetFormat(format TaFormat) *TrustAnchor {
	if o != nil {
		o.Format = format
	}
	return o
}

func (o TrustAnchor) GetFormat() TaFormat {
	return o.Format
}

func (o *TrustAnchor) SetData(data []byte) *TrustAnchor {
	if o != nil {
		o.Data = data
	}
	return o
}

func (o TrustAnchor) GetData() []byte {
	return o.Data
}

// ToCBOR serializes the target Meta to CBOR
func (o TrustAnchor) ToCBOR() ([]byte, error) {
	return em.Marshal(&o)
}

// FromCBOR deserializes the supplied CBOR data into the target Meta
func (o *TrustAnchor) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, o)
}

// FromJSON deserializes the supplied JSON data into the target Meta
func (o *TrustAnchor) FromJSON(data []byte) error {
	return json.Unmarshal(data, o)
}

// ToJSON serializes the target Meta to JSON
func (o TrustAnchor) ToJSON() ([]byte, error) {
	return json.Marshal(&o)
}

type TasAndCas struct {
	Tas []TrustAnchor `cbor:"0,keyasint" json:"tas"`
	Cas [][]byte      `cbor:"1,keyasint,omitempty" json:"cas,omitempty"`
}

func NewTasAndCas() *TasAndCas {
	return &TasAndCas{}
}

func (o *TasAndCas) AddCaCert(cert []byte) *TasAndCas {
	if o != nil {
		o.Cas = append(o.Cas, cert)
	}
	return o
}

func (o TasAndCas) GetCaCerts() [][]byte {
	return o.Cas
}

func (o *TasAndCas) AddTaCert(cert []byte) *TasAndCas {
	if o != nil {
		ta := TrustAnchor{}
		ta.Format = TaFormatCertificate
		ta.Data = cert
		o.Tas = append(o.Tas, ta)
	}
	return o
}

func (o TasAndCas) GetTas() []TrustAnchor {
	return o.Tas
}

// Valid checks for validity of the fields within Meta
func (o TasAndCas) Valid() error {
	if len(o.Tas) == 0 {
		return fmt.Errorf("empty TasAndCas")
	}

	return nil
}

// ToCBOR serializes the target Meta to CBOR
func (o TasAndCas) ToCBOR() ([]byte, error) {
	return em.Marshal(&o)
}

// FromCBOR deserializes the supplied CBOR data into the target Meta
func (o *TasAndCas) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, o)
}

// FromJSON deserializes the supplied JSON data into the target Meta
func (o *TasAndCas) FromJSON(data []byte) error {
	return json.Unmarshal(data, o)
}

// ToJSON serializes the target Meta to JSON
func (o TasAndCas) ToJSON() ([]byte, error) {
	return json.Marshal(&o)
}
