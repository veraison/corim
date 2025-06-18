// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

type ArtifactType uint8

const (
	ArtifactTypeEndorsedValues ArtifactType = iota
	ArtifactTypeTrustAnchors
	ArtifactTypeReferenceValues
)

// String returns the string representation of the target ArtifactType
func (a ArtifactType) String() string {
	switch a {
	case ArtifactTypeEndorsedValues:
		return "endorsed-values"
	case ArtifactTypeReferenceValues:
		return "reference-values"
	case ArtifactTypeTrustAnchors:
		return "trust-anchors"
	}
	// unreachable
	return ""
}
