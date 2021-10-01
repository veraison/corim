// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"

	"github.com/veraison/eat"
)

// UEID is an Unique Entity Identifier
type UEID eat.UEID

// TaggedUEID is an alias to allow automatic tagging of an UEID type
type TaggedUEID UEID

func (o UEID) Empty() bool {
	return len(o) == 0
}

// Valid checks that the target UEID is in one of the defined formats: IMEI, EUI or RAND
func (o UEID) Valid() error {
	if err := eat.UEID(o).Validate(); err != nil {
		return fmt.Errorf("UEID validation failed: %w", err)
	}
	return nil
}

// UnmarshalJSON deserializes the supplied string into the UEID target
func (o *UEID) UnmarshalJSON(data []byte) error {
	var b []byte

	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}

	u := UEID(b)

	if err := u.Valid(); err != nil {
		return err
	}

	*o = u

	return nil
}

func (o UEID) MarshalJSON() ([]byte, error) {
	return json.Marshal([]byte(o))
}
