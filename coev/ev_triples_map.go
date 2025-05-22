// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"github.com/veraison/corim/comid"
)

type EvTriplesMap struct {
	EvidenceTriples   *comid.ValueTriples
	IdentityTriples   *comid.KeyTriples
	AttestKeysTriples *comid.KeyTriples
}
