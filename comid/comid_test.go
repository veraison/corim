package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	c.RegisterExtensions(&ComidExt{})

	err = c.Set("field-one", "foo")
	assert.NoError(t, err)
	assert.Equal(t, "foo", c.MustGetString("-1"))
}

func Test_Comid_ToJSONPretty(t *testing.T) {
	c := NewComid()

	_, err := c.ToJSONPretty("    ")
	assert.EqualError(t, err, "tag-identity validation failed: empty tag-id")

	c.TagIdentity = TagIdentity{TagID: *swid.NewTagID("test"), TagVersion: 1}
	c.Triples = Triples{ReferenceValues: &[]ReferenceValue{}}

	expected := `{
    "tag-identity": {
        "id": "test",
        "version": 1
    },
    "triples": {
        "reference-values": []
    }
}`
	v, err := c.ToJSONPretty("    ")
	require.NoError(t, err)
	assert.Equal(t, expected, string(v))
}
