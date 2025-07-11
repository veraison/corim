// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import "github.com/veraison/corim/extensions"

const (
	ExtConciseEvidence      extensions.Point = "ConciseEvidence"
	ExtEvTriples            extensions.Point = "EvTriples"
	ExtEvidenceTriples      extensions.Point = "EvidenceTriples"
	ExtEvidenceTriplesFlags extensions.Point = "EvidenceTriplesFlags"
)

type Extensions struct {
	extensions.Extensions
}
