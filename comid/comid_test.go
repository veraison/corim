// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
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


func Test_Comid_SimpleReferenceValue(t *testing.T) {
    c := NewComid()
    env := Environment{
        Instance: MustNewUUIDInstance(TestUUID),
    }
    
    // Test digest reference value
	err := c.AddDigestReferenceValue(env, "sha-256", []byte{
		0x00, 0x01, 0x02, 0x03,
		0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b,
		0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13,
		0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b,
		0x1c, 0x1d, 0x1e, 0x1f,
	})
    require.NoError(t, err)
    
    // Verify values were added
    require.NotNil(t, c.Triples.ReferenceValues)
    require.Len(t, c.Triples.ReferenceValues.Values, 1)
    
    // Verify digest value
    rv := c.Triples.ReferenceValues.Values[0]
    require.NotNil(t, rv.Measurements.Values[0].Val.Digests)
    require.Equal(t, HashAlgSHA256.ToUint64(), (*rv.Measurements.Values[0].Val.Digests)[0].HashAlgID)
}