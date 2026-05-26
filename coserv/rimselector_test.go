// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/swid"
)

func TestRimSelectorID_Valid(t *testing.T) {
	testTag := *swid.NewTagID("foo")
	testCases := []struct {
		title    string
		selector RimSelectorID
		err      string
	}{
		{
			title:    "bad tag-id",
			selector: RimSelectorID{},
			err:      "tag-id value is nil",
		},
		{
			title: "bad type",
			selector: RimSelectorID{
				TagID: testTag,
				Type:  RimSelectorType(99),
			},
			err: "invalid type: 99",
		},
		{
			title: "ok",
			selector: RimSelectorID{
				TagID: testTag,
				Type:  RimSelectorTypeComid,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := tc.selector.Valid()
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}
		})
	}
}

func TestRimSelectorID_New(t *testing.T) {
	_, err := NewRimSelectorID(RimSelectorTypeComid, *swid.NewTagID("foo"))
	assert.NoError(t, err)

	_, err = NewRimSelectorID(RimSelectorType(66), *swid.NewTagID("foo"))
	assert.ErrorContains(t, err, "invalid type")
}
