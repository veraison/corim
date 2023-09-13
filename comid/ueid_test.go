package comid

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewTaggedUEID(t *testing.T) {
	ueid := UEID(TestUEID)
	tagged := TaggedUEID(TestUEID)
	bytes := MustHexDecode(t, TestUEIDString)

	for _, v := range []any{
		TestUEID,
		&TestUEID,
		ueid,
		&ueid,
		tagged,
		&tagged,
		bytes,
		base64.StdEncoding.EncodeToString(bytes),
	} {
		ret, err := NewTaggedUEID(v)
		require.NoError(t, err)
		assert.Equal(t, []byte(TestUEID), ret.Bytes())
	}
}
