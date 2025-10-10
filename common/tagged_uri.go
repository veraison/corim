// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package common

// TaggedURI represents a URI string. It is used by both comid and corim
// packages for entity registration IDs.
type TaggedURI string

// Empty returns true if the TaggedURI is an empty string.
func (o TaggedURI) Empty() bool {
	return o == ""
}

// String2URI converts a string pointer to a TaggedURI pointer.
// Returns nil if the input is nil or empty.
func String2URI(s *string) (*TaggedURI, error) {
	if s == nil || *s == "" {
		return nil, nil
	}

	uri := TaggedURI(*s)
	return &uri, nil
}
