// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/swid"
)

func Test_Comid_Extensions(t *testing.T) {
	c := NewComid()
	assert.Nil(t, c.GetExtensions())
	assert.Equal(t, "", c.MustGetString("field-one"))

	err := c.Set("field-one", "foo")
	assert.EqualError(t, err, "extension not found: field-one")

	type ComidExt struct {
		FieldOne string `cbor:"-1,keyasint" json:"field-one"`
	}

	extMap := extensions.NewMap().
		Add(ExtComid, &ComidExt{}).
		Add(ExtEntity, &struct{}{}).
		Add(ExtReferenceValue, &struct{}{}).
		Add(ExtEndorsedValueFlags, &struct{}{})
	err = c.RegisterExtensions(extMap)
	require.NoError(t, err)

	err = c.Set("field-one", "foo")
	assert.NoError(t, err)
	assert.Equal(t, "foo", c.MustGetString("-1"))
}

func Test_Comid_ToJSONPretty(t *testing.T) {
	c := NewComid()

	_, err := c.ToJSONPretty("    ")
	assert.EqualError(t, err, "tag-identity validation failed: empty tag-id")

	c.TagIdentity = TagIdentity{TagID: *swid.NewTagID("test"), TagVersion: 1}
	c.Triples = Triples{
		ReferenceValues: NewValueTriples().Add(&ValueTriple{
			Environment: Environment{
				Instance: MustNewUUIDInstance(TestUUID),
			},
			Measurements: *NewMeasurements().Add(&Measurement{
				Val: Mval{
					RawValue: NewRawValue().SetBytes(MustHexDecode(t, "deadbeef")),
				},
			}),
		}),
	}

	expected := `{
    "tag-identity": {
        "id": "test",
        "version": 1
    },
    "triples": {
        "reference-values": [
            {
                "environment": {
                    "instance": {
                        "type": "uuid",
                        "value": "31fb5abf-023e-4992-aa4e-95f9c1503bfa"
                    }
                },
                "measurements": [
                    {
                        "value": {
                            "raw-value": {
                                "type": "bytes",
                                "value": "3q2+7w=="
                            }
                        }
                    }
                ]
            }
        ]
    }
}`
	v, err := c.ToJSONPretty("    ")
	require.NoError(t, err)
	assert.Equal(t, expected, string(v))
}

func Test_String2URI_nok(t *testing.T) {
	s := "@@@"
	_, err := String2URI(&s)
	assert.EqualError(t, err, `expecting an absolute URI: "@@@" is not an absolute URI`)
}

func TestComid_iterators(t *testing.T) {
	c := NewTestComid(t)

	assert.Equal(t, *slices.Collect(c.IterAttestVerifKeys())[0], (*c.Triples.AttestVerifKeys)[0])
	assert.Equal(t, *slices.Collect(c.IterDevIdentityKeys())[0], (*c.Triples.DevIdentityKeys)[0])
	assert.Equal(t, *slices.Collect(c.IterRefVals())[0], c.Triples.ReferenceValues.Values[0])
	assert.Equal(t, *slices.Collect(c.IterEndVals())[0], c.Triples.EndorsedValues.Values[0])
}
