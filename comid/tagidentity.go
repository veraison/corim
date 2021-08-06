// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"

	"github.com/veraison/swid"
)

type TagIdentity struct {
	TagID      swid.TagID `cbor:"0,keyasint" json:"id"`
	TagVersion uint       `cbor:"1,keyasint,omitempty" json:"version,omitempty"`
}

func (o TagIdentity) Valid() error {
	if o.TagID == (swid.TagID{}) {
		return fmt.Errorf("empty tag-id")
	}

	return nil
}
