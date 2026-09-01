// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cotl

import (
	"errors"
	"time"
)

// ErrEmptyTagsList is returned when a CoTL carries no tags: activation is
// atomic per draft-ietf-rats-corim-11 Section 6, so an empty list can never
// be meaningfully activated.
var ErrEmptyTagsList = errors.New("empty tags-list: a CoTL must activate at least one tag")

// Report is the outcome of appraising a CoTL at a point in time.
type Report struct {
	// Errors make the CoTL unusable for appraisal.
	Errors []string
	// Warnings flag conditions that do not invalidate the CoTL.
	Warnings []string
}

// Valid returns true when the report carries no errors.
func (r Report) Valid() bool {
	return len(r.Errors) == 0
}

// ValidateAt appraises the CoTL at time now, mirroring the CoTL extraction
// rules of draft-ietf-rats-corim-11 Section 8.2.3.4:
//
//   - CoTLs outside their validity window MUST be discarded
//     (expired => error; not-yet-valid => warning).
//   - The validity window is inclusive of both bounds (RFC 5280 convention).
//   - Duplicate entries in tags-list are reported as warnings.
func (o *ConciseTlTag) ValidateAt(now time.Time) Report {
	var report Report

	if o.TlValidity.NotBefore != nil {
		if now.Before(*o.TlValidity.NotBefore) {
			report.Warnings = append(report.Warnings,
				"CoTL is not yet valid (not-before is in the future)")
		}
	}

	if now.After(o.TlValidity.NotAfter) {
		report.Errors = append(report.Errors,
			"CoTL has expired (not-after is in the past)")
	}

	seen := make(map[string]bool, len(o.TagsList))
	for _, t := range o.TagsList {
		id := t.String()
		if seen[id] {
			report.Warnings = append(report.Warnings,
				"duplicate tag-id in tags-list: "+id)
			continue
		}
		seen[id] = true
	}

	return report
}
