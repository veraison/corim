// Copyright 2025-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
)

func TestQuery_Valid(t *testing.T) {
	artifactType := ArtifactTypeReferenceValues
	resultType := ResultTypeBoth
	instance := comid.MustNewBytesInstance(comid.MustHexDecode(t, "deadbeef"))
	tagID := *swid.NewTagID("foo")

	testCases := []struct {
		title string
		query Query
		err   string
	}{
		{
			title: "empty",
			query: Query{},
			err:   "no selector specified",
		},
		{
			title: "both",
			query: Query{
				EnvironmentSelector: &EnvironmentSelector{},
				RimSelector:         NewRimSelectorIDs(),
			},
			err: "cannot be specified at the same time",
		},
		{
			title: "env no artificat type",
			query: Query{
				EnvironmentSelector: &EnvironmentSelector{},
			},
			err: "artifact type must be specified",
		},
		{
			title: "env no result type",
			query: Query{
				EnvironmentSelector: &EnvironmentSelector{},
				ArtifactType:        &artifactType,
			},
			err: "result type must be specified",
		},
		{
			title: "env invalid environment selector",
			query: Query{
				EnvironmentSelector: &EnvironmentSelector{},
				ArtifactType:        &artifactType,
				ResultType:          &resultType,
			},
			err: "invalid environment selector",
		},
		{
			title: "env ok",
			query: Query{
				EnvironmentSelector: &EnvironmentSelector{
					Instances: &[]StatefulInstance{
						{
							Instance: instance,
						},
					},
				},
				ArtifactType: &artifactType,
				ResultType:   &resultType,
			},
		},
		{
			title: "rim invalid",
			query: Query{
				RimSelector: NewRimSelectorIDs(),
			},
			err: "invalid RIM selector: empty",
		},
		{
			title: "rim with artifact type",
			query: Query{
				ArtifactType: &artifactType,
				RimSelector: NewRimSelectorIDs().Add(&RimSelectorID{
					TagID: tagID,
					Type:  RimSelectorTypeCorim,
				}),
			},
			err: "artifact type cannot be specified",
		},
		{
			title: "rim with result type",
			query: Query{
				ResultType: &resultType,
				RimSelector: NewRimSelectorIDs().Add(&RimSelectorID{
					TagID: tagID,
					Type:  RimSelectorTypeCorim,
				}),
			},
			err: "result type cannot be specified",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := tc.query.Valid()
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}
		})
	}
}

func TestNewEnvironmentQuery(t *testing.T) {
	artifactType := ArtifactTypeReferenceValues
	resultType := ResultTypeBoth
	selector := *NewEnvironmentSelector().AddInstance(StatefulInstance{
		Instance: comid.MustNewBytesInstance(comid.MustHexDecode(t, "deadbeef")),
	})

	_, err := NewEnvironmentQuery(artifactType, EnvironmentSelector{}, resultType)
	assert.ErrorContains(t, err, "invalid environment selector")

	_, err = NewEnvironmentQuery(artifactType, selector, resultType)
	assert.NoError(t, err)
}

func TestNewRimQuery(t *testing.T) {
	tagID := *swid.NewTagID("foo")

	_, err := NewRimQuery(RimSelectorTypeCorim, swid.TagID{})
	assert.ErrorContains(t, err, "tag-id value is nil")

	_, err = NewRimQuery(RimSelectorTypeCorim, tagID)
	assert.NoError(t, err)
}
