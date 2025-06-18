// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

type ResultType uint8

const (
	ResultTypeCollectedArtifacts ResultType = iota
	ResultTypeSourceArtifacts
	ResultTypeBoth
)

// String returns the string representation of the target ResultType
func (a ResultType) String() string {
	switch a {
	case ResultTypeCollectedArtifacts:
		return "collected-artifacts"
	case ResultTypeSourceArtifacts:
		return "source-artifacts"
	case ResultTypeBoth:
		return "both"
	}
	// unreachable
	return ""
}
