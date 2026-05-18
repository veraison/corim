// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/corim/comid"
)

func TestNewProfileFromString(t *testing.T) {
	profile, err := NewProfileFromString("1.2.3.4")
	assert.NoError(t, err)
	assert.Equal(t, comid.OIDType, profile.Type())

	profile, err = NewProfileFromString("http://example.com")
	assert.NoError(t, err)
	assert.Equal(t, comid.URIType, profile.Type())

	_, err = NewProfileFromString("http://[::")
	assert.ErrorContains(t, err, "missing ']' in host")
}

func TestNewXXXProfile(t *testing.T) {
	profile, err := NewURIProfile("http://[::")
	assert.Nil(t, profile)
	assert.ErrorContains(t, err, "missing ']' in host")

	assert.Panics(t, func() { MustNewURIProfile("http://[::") })

	profile, err = NewURIProfile("http://example.com")
	assert.NotNil(t, profile)
	assert.NoError(t, err)

	profile, err = NewOIDProfile("foo")
	assert.Nil(t, profile)
	assert.ErrorContains(t, err, "invalid OID")

	assert.Panics(t, func() { MustNewOIDProfile("foo") })

	profile, err = NewURIProfile("1.2.3.4")
	assert.NotNil(t, profile)
	assert.NoError(t, err)
}

func TestNewProfile_unknown_type(t *testing.T) {
	_, err := NewProfile("foo", "test")
	assert.ErrorContains(t, err, "unknown profile")
}

func TestProfile_UnmarshalJSON_bad(t *testing.T) {
	testCases := []struct {
		title string
		text  string
		err   string
	}{
		{
			title: "invalid JSON",
			text:  "foo",
			err:   "invalid character 'o'",
		},
		{
			title: "invalid type-and-value",
			text:  `"foo"`,
			err:   "cannot unmarshal string into Go value of type struct",
		},
		{
			title: "unknown type",
			text:  `{"type": "foo", "value": "bar"}`,
			err:   "unknown profile type: foo",
		},
		{
			title: "bad value",
			text:  `{"type": "oid", "value": "foo"}`,
			err:   "invalid OID",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var decoded Profile
			err := decoded.UnmarshalJSON([]byte(tc.text))
			assert.ErrorContains(t, err, tc.err)
		})
	}
}

type intProfile int

func (o intProfile) String() string {
	return fmt.Sprint(int(o))
}

func (o intProfile) Valid() error {
	return nil
}

func (o intProfile) Type() string {
	return "int"
}

func badNewIntProfile(val any) (*Profile, error) {
	switch t := val.(type) {
	case int:
		ret := intProfile(t)
		return &Profile{&ret}, nil
	default:
		return nil, fmt.Errorf("invalid value for int profile: %v (%T)", t, t)
	}
}

func newIntProfile(val any) (*Profile, error) {
	if val == nil {
		ret := intProfile(0)
		return &Profile{&ret}, nil
	}

	return badNewIntProfile(val)
}

func TestRegisterProfileType(t *testing.T) {
	tag := uint64(0xdeadbeef)

	err := RegisterProfileType(tag, badNewIntProfile)
	assert.ErrorContains(t, err, "invalid value for int profile")

	err = RegisterProfileType(32, newIntProfile)
	assert.ErrorContains(t, err, "tag 32 is already registered")

	err = RegisterProfileType(tag, newIntProfile)
	assert.NoError(t, err)

	err = RegisterProfileType(tag, newIntProfile)
	assert.ErrorContains(t, err, `name "int" already exists`)
}
