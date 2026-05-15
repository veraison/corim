// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"errors"
	"fmt"

	"github.com/veraison/swid"
)

type CoswidTagIDs []swid.TagID

func (o CoswidTagIDs) Valid() error {
	if len(o) == 0 {
		return errors.New("must not be empty")
	}

	for i, tag := range o {
		if err := tag.Valid(); err != nil {
			return fmt.Errorf("tag-id[%d]: %w", i, err)
		}
	}

	return nil
}

func (o *CoswidTagIDs) Add(v *swid.TagID) *CoswidTagIDs {
	*o = append(*o, *v)
	return o
}

// CoswidTriple relates reference measurements contained in one or more CoSWIDs
// to a Target Environment.
//
//	;# import rfc9393 as coswid
//
//	coswid-triple-record = [
//	  environment-map
//	  [ + coswid.tag-id ]
//	]
type CoswidTriple struct {
	_           struct{}     `cbor:",toarray"`
	Environment Environment  `json:"environment"`
	TagIDs      CoswidTagIDs `json:"tag-ids"`
}

func (o CoswidTriple) Valid() error {
	if err := o.Environment.Valid(); err != nil {
		return fmt.Errorf("environment: %w", err)
	}

	if err := o.TagIDs.Valid(); err != nil {
		return fmt.Errorf("tag-ids: %w", err)
	}

	return nil
}

type CoswidTriples []CoswidTriple

func NewCoswidTriples() *CoswidTriples {
	return &CoswidTriples{}
}

func (o *CoswidTriples) Add(triple *CoswidTriple) *CoswidTriples {
	*o = append(*o, *triple)
	return o
}

func (o CoswidTriples) IsEmpty() bool {
	return len(o) == 0
}

func (o CoswidTriples) Valid() error {
	for i, triple := range o {
		if err := triple.Valid(); err != nil {
			return fmt.Errorf("triple[%d]: %w", i, err)
		}
	}

	return nil
}
