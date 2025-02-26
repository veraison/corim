// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"

	"github.com/veraison/swid"
)

// Version stores a version-map with JSON and CBOR serializations.
type Version struct {
	Version string             `cbor:"0,keyasint" json:"value"`
	Scheme  swid.VersionScheme `cbor:"1,keyasint" json:"scheme"`
}

func NewVersion() *Version {
	return &Version{}
}

func (o *Version) SetVersion(v string) *Version {
	if o != nil {
		o.Version = v
	}
	return o
}

func (o *Version) SetScheme(v int64) *Version {
	if o != nil {
		if o.Scheme.SetCode(v) != nil {
			return nil
		}
	}
	return o
}

func (o Version) Valid() error {
	if o.Version == "" {
		return fmt.Errorf("empty version")
	}
	return nil
}

func (o Version) Equal(r Version) bool {
	if o.Scheme != r.Scheme || o.Version != r.Version {
		return false
	}

	return true
}

func (o Version) CompareAgainstReference(r Version) bool {
	return o.Equal(r)
}
