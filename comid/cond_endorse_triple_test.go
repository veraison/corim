// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/corim/extensions"
)

func TestCondEndorseTriple_Valid(t *testing.T) {
	testCases := []struct {
		title  string
		triple CondEndorseTriple
		err    string
	}{
		{
			title:  "bad: empty conditions",
			triple: CondEndorseTriple{},
			err:    "no condition entries",
		},
		{
			title: "bad: invalid stateful environment",
			triple: CondEndorseTriple{
				Conditions: *NewStatefulEnvironments().Add(&StatefulEnvironment{}),
			},
			err: "environment must not be empty",
		},
		{
			title: "bad: empty endorsements",
			triple: CondEndorseTriple{
				Conditions: *NewStatefulEnvironments().Add(
					&StatefulEnvironment{
						Environment: Environment{
							Instance: MustNewUUIDInstance(TestUUID),
						},
						Measurements: *NewMeasurements().Add(
							&Measurement{
								Val: Mval{
									SVN: MustNewSVN(1, "exact-value"),
								},
							},
						),
					},
				),
			},
			err: "no endorsement entries",
		},
		{
			title: "bad: invalid endorsements",
			triple: CondEndorseTriple{
				Conditions: *NewStatefulEnvironments().Add(
					&StatefulEnvironment{
						Environment: Environment{
							Instance: MustNewUUIDInstance(TestUUID),
						},
						Measurements: *NewMeasurements().Add(
							&Measurement{
								Val: Mval{
									SVN: MustNewSVN(1, "exact-value"),
								},
							},
						),
					},
				),
				Endorsements: *NewValueTriples().Add(&ValueTriple{}),
			},
			err: "environment must not be empty",
		},
		{
			title: "ok",
			triple: CondEndorseTriple{
				Conditions: *NewStatefulEnvironments().Add(
					&StatefulEnvironment{
						Environment: Environment{
							Instance: MustNewUUIDInstance(TestUUID),
						},
						Measurements: *NewMeasurements().Add(
							&Measurement{
								Val: Mval{
									SVN: MustNewSVN(1, "exact-value"),
								},
							},
						),
					},
				),
				Endorsements: *NewValueTriples().Add(
					&StatefulEnvironment{
						Environment: Environment{
							Instance: MustNewUUIDInstance(TestUUID),
						},
						Measurements: *NewMeasurements().Add(
							&Measurement{
								Val: Mval{
									SVN: MustNewSVN(1, "exact-value"),
								},
							},
						),
					},
				),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := tc.triple.Valid()
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}
		})
	}
}

func TestCondEndorseTriple_extensions(t *testing.T) {
	exts := extensions.NewMap()
	exts.Add(ExtMval, &struct{}{})

	triple := &CondEndorseTriple{}
	err := triple.RegisterExtensions(exts)
	assert.NoError(t, err)

	ret := triple.GetExtensions()
	assert.NotNil(t, ret)
	assert.EqualValues(t, ret, exts)

	exts.Add(ExtEntity, &struct{}{})
	err = triple.RegisterExtensions(exts)
	assert.ErrorContains(t, err, "unexpected extension point")
}

func TestNewCondEndorseTriples(t *testing.T) {
	triples := NewCondEndorseTriples()
	assert.NotNil(t, triples)
	assert.True(t, triples.IsEmpty())
}
