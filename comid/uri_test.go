// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTaggedURI(t *testing.T) {
	textValue := "http://example.com"
	urlValue := url.URL{Scheme: "http", Host: "example.com"}
	badTextValue := "http://[::"

	testCases := []struct {
		title string
		value any
		err   string
	}{
		{
			title: "ok nil",
			value: nil,
		},
		{
			title: "ok string",
			value: textValue,
		},
		{
			title: "ok *string",
			value: &textValue,
		},
		{
			title: "ok url.URL",
			value: urlValue,
		},
		{
			title: "ok *url.URL",
			value: &urlValue,
		},
		{
			title: "bad invalid string",
			value: badTextValue,
			err:   "missing ']' in host",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			_, err := NewTaggedURI(tc.value)
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}
		})
	}
}

func TestTaggedURI_ITypeChoiceValue_methods(t *testing.T) {
	u := TaggedURI("")
	err := u.Valid()
	assert.ErrorContains(t, err, "empty URI")

	u = TaggedURI("http://example.com")
	assert.Equal(t, "http://example.com", u.String())
	assert.Equal(t, "uri", u.Type())
}
