// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/extensions"
)

func TestConciseEvidence_NewConciseEvidence(t *testing.T) {
	coev := NewConciseEvidence()
	require.NotNil(t, coev)
}

func TestConciseEvidence_AddTriples_NOK(t *testing.T) {
	expectedErr := "no evidence triples"
	var ev *EvTriples
	coev := &ConciseEvidence{}
	err := coev.AddTriples(ev)
	assert.EqualError(t, err, expectedErr)
	expectedErr = "invalid evidence triples: no Triples set inside EvTriples"
	ev = &EvTriples{}
	err = coev.AddTriples(ev)
	assert.EqualError(t, err, expectedErr)
}

func TestConciseEvidence_AddEvidenceID(t *testing.T) {
	coev := &ConciseEvidence{}
	ev := MustNewUUIDEvidenceID(TestUUID)
	require.NotNil(t, ev)
	err := coev.AddEvidenceID(ev)
	require.NoError(t, err)
}

func TestConciseEvidence_AddEvidenceID_NOK(t *testing.T) {
	coev := &ConciseEvidence{}
	expectedErr := "invalid EvidenceID: no EvidenceID"
	var e EvidenceID
	err := coev.AddEvidenceID(&e)
	assert.EqualError(t, err, expectedErr)
}

func TestConciseEvidence_AddProfile(t *testing.T) {
	coev := &ConciseEvidence{}
	err := coev.AddProfile(TestProfile)
	require.NoError(t, err)

}

func TestConciseEvidence_AddProfile_NOK(t *testing.T) {
	coev := &ConciseEvidence{}
	expectedErr := "profile string must be an absolute URL or an ASN.1 OID: no valid OID"
	var p string
	err := coev.AddProfile(p)
	assert.EqualError(t, err, expectedErr)
	expectedErr = `profile string must be an absolute URL or an ASN.1 OID: failed to extract OID from string: strconv.Atoi: parsing "not": invalid syntax`
	p = "not"
	err = coev.AddProfile(p)
	assert.EqualError(t, err, expectedErr)
}

func TestConciseEvidence_Valid_NOK(t *testing.T) {
	expectedErr := "invalid EvTriples: no Triples set inside EvTriples"
	coev := &ConciseEvidence{}
	err := coev.Valid()
	assert.EqualError(t, err, expectedErr)
}

func Test_ConciseEvidence_Extensions(t *testing.T) {
	c := NewConciseEvidence()
	assert.Nil(t, c.GetExtensions())
	assert.Equal(t, "", c.MustGetString("myparam"))

	err := c.Set("myparam", "test-param")
	assert.EqualError(t, err, "extension not found: myparam")

	type CoEvExt struct {
		MyParam string `cbor:"-1,keyasint" json:"myparam"`
	}

	extMap := extensions.NewMap().
		Add(ExtConciseEvidence, &CoEvExt{}).
		Add(ExtEvTriples, &struct{}{}).
		Add(ExtEvidenceTriples, &struct{}{}).
		Add(ExtEvidenceTriplesFlags, &struct{}{})

	err = c.RegisterExtensions(extMap)
	require.NoError(t, err)

	err = c.Set("myparam", "test-param")
	assert.NoError(t, err)
	assert.Equal(t, "test-param", c.MustGetString("-1"))
}

func Test_ConciseEvidence_RegisterExtensions_NOK(t *testing.T) {
	expectedErr := `unexpected extension point: "myPoint"`
	c := NewConciseEvidence()
	type CoEvExt struct {
		MyParam string `cbor:"-1,keyasint" json:"myparam"`
	}
	extMap := extensions.NewMap().
		Add("myPoint", &CoEvExt{})
	err := c.RegisterExtensions(extMap)
	assert.EqualError(t, err, expectedErr)
}
