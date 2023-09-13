package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testGroup uint64

func newTestGroup(val any) (*Group, error) {
	v := testGroup(7)
	return &Group{&v}, nil
}

func (o testGroup) Type() string {
	return "test-value"
}

func (o testGroup) String() string {
	return "test"
}

func (o testGroup) Valid() error {
	return nil
}

type testGroupBadType struct {
	testGroup
}

func newTestGroupBadType(val any) (*Group, error) {
	v := testGroupBadType{testGroup(7)}
	return &Group{&v}, nil
}

func (o testGroupBadType) Type() string {
	return "uuid"
}

func Test_RegisterGroupType(t *testing.T) {
	err := RegisterGroupType(32, newTestGroup)
	assert.EqualError(t, err, "tag 32 is already registered")

	err = RegisterGroupType(99993, newTestGroupBadType)
	assert.EqualError(t, err, `Group type with name "uuid" already exists`)

	err = RegisterGroupType(99993, newTestGroup)
	require.NoError(t, err)

}

func TestGroup_UmarshalJSON(t *testing.T) {
	var group Group

	err := group.UnmarshalJSON([]byte(`{`))
	assert.EqualError(t, err, "group decoding failure: unexpected end of JSON input")

	err = group.UnmarshalJSON([]byte(`{"type":"uuid","value":"aaaa"}`))
	assert.EqualError(t, err, "cannot unmarshal group: bad UUID: invalid UUID length: 4")
}
