// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "encoding/json"

// tnv (type'n'value) stores a JSON object with two attributes: a string "type"
// and a generic "value" (left undecoded) defined by type.  This type is used in
// a few places to implement the choice types that CBOR handles using tags.
type tnv struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}
