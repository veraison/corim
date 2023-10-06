// Copyright 2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package extensions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Entity struct {
	EntityName string
	Roles      []int64

	Extensions
}

type TestExtensions struct {
	Address     string            `cbor:"-1,keyasint,omitempty" json:"address,omitempty"`
	Size        int               `cbor:"-2,keyasint,omitempty" json:"size,omitempty"`
	YearsOnAir  float32           `cbor:"-3,keyasint,omitempty" json:"years-on-air,omitempty"`
	StillAiring bool              `cbor:"-4,keyasint,omitempty" json:"still-airing,omitempty"`
	Ages        []int             `cbor:"-5,keyasint,omitempty" json:"ages,omitempty"`
	Jobs        map[string]string `cbor:"-6,keyasint,omitempty" json:"jobs,omitempty"`
}

func TestExtensions_Register(t *testing.T) {
	exts := Extensions{}
	assert.False(t, exts.HaveExtensions())

	exts.Register(&TestExtensions{})
	assert.True(t, exts.HaveExtensions())

	badRegister := func() {
		exts.Register(TestExtensions{})
	}

	assert.Panics(t, badRegister)
}

func TestExtensions_GetSet(t *testing.T) {
	extsVal := TestExtensions{
		Address:     "742 Evergreen Terrace",
		Size:        6,
		YearsOnAir:  33.8,
		StillAiring: true,
		Ages:        []int{2, 7, 8, 10, 37, 38},
		Jobs: map[string]string{
			"Homer": "safety inspector",
			"Marge": "housewife",
			"Bart":  "elementary school student",
			"Lisa":  "elementary school student",
		},
	}
	exts := Extensions{IExtensionsValue: &extsVal}

	v, err := exts.GetInt("size")
	assert.NoError(t, err)
	assert.Equal(t, 6, v)

	assert.Equal(t, 6, exts.MustGetInt("size"))
	assert.Equal(t, int64(6), exts.MustGetInt64("size"))
	assert.Equal(t, int32(6), exts.MustGetInt32("size"))
	assert.Equal(t, int16(6), exts.MustGetInt16("size"))
	assert.Equal(t, int8(6), exts.MustGetInt8("size"))

	assert.Equal(t, uint(6), exts.MustGetUint("size"))
	assert.Equal(t, uint64(6), exts.MustGetUint64("size"))
	assert.Equal(t, uint32(6), exts.MustGetUint32("size"))
	assert.Equal(t, uint16(6), exts.MustGetUint16("size"))
	assert.Equal(t, uint8(6), exts.MustGetUint8("size"))

	assert.InEpsilon(t, float32(33.8), exts.MustGetFloat32("years-on-air"), 0.000001)
	assert.InEpsilon(t, float64(33.8), exts.MustGetFloat64("-3"), 0.000001)

	assert.Equal(t, true, exts.MustGetBool("StillAiring"))

	_, err = exts.GetSlice("ages")
	assert.EqualError(t, err,
		`unable to cast []int{2, 7, 8, 10, 37, 38} of type []int to []interface{}`)
	assert.Nil(t, exts.MustGetSlice("ages"))

	assert.EqualValues(t, []int{2, 7, 8, 10, 37, 38}, exts.MustGetIntSlice("ages"))
	assert.EqualValues(t, []string{"2", "7", "8", "10", "37", "38"},
		exts.MustGetStringSlice("ages"))

	assert.EqualValues(t, map[string]string{
		"Homer": "safety inspector",
		"Marge": "housewife",
		"Bart":  "elementary school student",
		"Lisa":  "elementary school student",
	}, exts.MustGetStringMapString("jobs"))

	_, err = exts.GetStringMap("jobs")
	assert.EqualError(t, err,
		`unable to cast map[string]string{"Bart":"elementary school student", "Homer":"safety inspector", "Lisa":"elementary school student", "Marge":"housewife"} of type map[string]string to map[string]interface{}`)
	m := exts.MustGetStringMap("jobs")
	assert.Equal(t, map[string]any{}, m)

	s, err := exts.GetString("address")
	assert.NoError(t, err)
	assert.Equal(t, "742 Evergreen Terrace", s)

	_, err = exts.GetInt("address")
	assert.EqualError(t, err, `unable to cast "742 Evergreen Terrace" of type string to int`)

	_, err = exts.GetInt("foo")
	assert.EqualError(t, err, "extension not found: foo")

	err = exts.Set("-1", "123 Fake Street")
	assert.NoError(t, err)

	s, err = exts.GetString("address")
	assert.NoError(t, err)
	assert.Equal(t, "123 Fake Street", s)

	err = exts.Set("Size", "foo")
	assert.EqualError(t, err, `cannot set field "Size" (of type int) to foo (string)`)

	assert.Equal(t, "", exts.MustGetString("does-not-exist"))
	assert.Equal(t, 0, exts.MustGetInt("does-not-exist"))
}
