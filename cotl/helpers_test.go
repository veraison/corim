// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cotl

import (
	cbor "github.com/fxamacker/cbor/v2"
)

func cborDecOptions() cbor.DecOptions {
	return cbor.DecOptions{
		IndefLength: cbor.IndefLengthAllowed,
		TimeTag:     cbor.DecTagOptional,
	}
}
