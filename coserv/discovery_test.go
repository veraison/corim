package coserv

import (
	"encoding/json"
	"maps"
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
)

func validTestDiscovery() DiscoveryDocument {
	var dd DiscoveryDocument
	dd.SetVersion("1.2.3")
	dd.AddCapability("foo", []ArtifactSupport{ArtifactSupportSource, ArtifactSupportCollected})
	dd.AddEndPoint("foo", "bar")
	dd.AddEndPoint("CoSERVRequestResponse", "bar/{query}")
	return dd
}

func TestDiscoveryJsonSerDe(t *testing.T) {
	key := readTestVectorSlice(t, "sample_key.jwk")
	dd := validTestDiscovery()
	dd.AddJwk(key)
	s, err := dd.ToJSON()
	assert.NoError(t, err)
	var ddNew DiscoveryDocument
	err = ddNew.FromJSON(s)
	assert.NoError(t, err)
}

func TestDiscoveryCborSerDe(t *testing.T) {
	key := readTestVectorSlice(t, "sample_key.cbor")
	dd := validTestDiscovery()
	dd.AddCoseKey(key)
	s, err := dd.ToCBOR()
	assert.NoError(t, err)
	var ddNew DiscoveryDocument
	err = ddNew.FromCBOR(s)
	assert.NoError(t, err)
}

func TestDiscoveryJsonNoKey(t *testing.T) {
	dd := validTestDiscovery()
	s, err := dd.ToJSON()
	assert.NoError(t, err)
	var ddNew DiscoveryDocument
	err = ddNew.FromJSON(s)
	assert.NoError(t, err)
}

func TestDiscoveryCborNoKey(t *testing.T) {
	dd := validTestDiscovery()
	s, err := dd.ToCBOR()
	assert.NoError(t, err)
	var ddNew DiscoveryDocument
	err = ddNew.FromCBOR(s)
	assert.NoError(t, err)
}

func TestDiscoveryInvalidVersion(t *testing.T) {
	dd := validTestDiscovery()
	dd.SetVersion("invalid")
	assert.Error(t, dd.Validate())
}

func TestDiscoveryInvalidMime(t *testing.T) {
	dd := validTestDiscovery()
	dd.AddCapability("@", []ArtifactSupport{ArtifactSupportSource})
	assert.Error(t, dd.Validate())
}

func TestDiscoveryInvalidJwk(t *testing.T) {
	dd := validTestDiscovery()
	k := []byte("invalid")
	dd.AddJwk(k)
	assert.Error(t, dd.Validate())
}

func TestDiscoveryInvalidCoseKey(t *testing.T) {
	dd := validTestDiscovery()
	k := []byte("invalid")
	dd.AddCoseKey(k)
	assert.Error(t, dd.Validate())
}

func TestDiscoveryEmptyCapabilities(t *testing.T) {
	dd := validTestDiscovery()
	dd.CapabilitiesList = nil
	assert.EqualError(t, dd.Validate(), ErrEmptyCapabilities.Error())
}

func TestDiscoveryEmptyEndPoints(t *testing.T) {
	dd := validTestDiscovery()
	dd.ApiEndPointsMap = nil
	assert.EqualError(t, dd.Validate(), ErrEmptyApiEndPoints.Error())
}

func TestDiscoveryInvalidRequestRespEp(t *testing.T) {
	dd := validTestDiscovery()
	dd.ApiEndPointsMap["CoSERVRequestResponse"] = "foo"
	assert.ErrorAs(t, dd.Validate(), &ErrInvalidRequestResponseEndpoint)

	dd = validTestDiscovery()
	dd.AddEndPoint("CoSERVRequestResponse", "/{query")
	assert.ErrorAs(t, dd.Validate(), &ErrInvalidRequestResponseEndpoint)

	dd = validTestDiscovery()
	dd.AddEndPoint("CoSERVRequestResponse", "/{w}/{query}")
	assert.ErrorAs(t, dd.Validate(), &ErrInvalidRequestResponseEndpoint)
}

func TestDiscoveryNoRequestRespEp(t *testing.T) {
	dd := validTestDiscovery()
	delete(dd.ApiEndPointsMap, "CoSERVRequestResponse")
	assert.EqualError(t, dd.Validate(), ErrNoRequestResponseEndpoint.Error())
}

func TestDiscoveryValidJsonDe(t *testing.T) {
	var dd DiscoveryDocument
	tests := []string{
		"discovery-single-capability.json",
		"discovery-unsigned.json",
	}
	for _, test := range tests {
		j := readTestVectorSlice(t, test)
		assert.NoError(t, dd.FromJSON(j))
	}
}

func TestDiscoveryValidCborDe(t *testing.T) {
	var dd DiscoveryDocument
	tests := []string{
		"discovery-single-capability.cbor",
		"discovery-unsigned.cbor",
	}
	for _, test := range tests {
		c := readTestVectorSlice(t, test)
		assert.NoError(t, dd.FromCBOR(c))
	}
}

func TestDiscoveryCapIter(t *testing.T) {
	dd := validTestDiscovery()
	dd.AddCapability("bar", []ArtifactSupport{ArtifactSupportSource})
	caps := make([]Capability, 0)
	for mt, rt := range dd.Capabilities() {
		caps = append(caps, Capability{mt, rt})
	}
	assert.Equal(t, caps, dd.CapabilitiesList)
}

func TestDiscoveryApiEpIter(t *testing.T) {
	dd := validTestDiscovery()
	dd.AddEndPoint("bar", "baz")
	eps := maps.Collect(dd.ApiEndPoints())
	assert.Equal(t, eps, dd.ApiEndPointsMap)
}

func TestDiscoveryArsupEncode(t *testing.T) {
	valid := func(sup ArtifactSupport, j, c []byte) {
		var (
			s   []byte
			err error
		)
		s, err = json.Marshal(sup)
		assert.NoError(t, err)
		assert.Equal(t, s, j, "json encoding error")

		s, err = cbor.Marshal(sup)
		assert.NoError(t, err)
		assert.Equal(t, s, c, "cbor encoding error")
	}

	valid(ArtifactSupportSource,
		[]byte("\"source\""),
		[]byte{0x66, 0x73, 0x6F, 0x75, 0x72, 0x63, 0x65},
	)
	valid(ArtifactSupportCollected,
		[]byte("\"collected\""),
		[]byte{0x69, 0x63, 0x6F, 0x6C, 0x6C, 0x65, 0x63, 0x74, 0x65, 0x64},
	)
	valid(ArtifactSupportRims,
		[]byte("\"rims\""),
		[]byte{0x64, 0x72, 0x69, 0x6D, 0x73},
	)
}

func TestDiscoveryArsupDecode(t *testing.T) {
	var (
		j   []byte
		c   []byte
		err error
	)
	valid := func(j, c []byte, exp ArtifactSupport) {
		var sup ArtifactSupport
		err = sup.UnmarshalJSON(j)
		assert.NoError(t, err)
		assert.Equal(t, sup, exp)

		err = sup.UnmarshalCBOR(c)
		assert.NoError(t, err)
		assert.Equal(t, sup, exp)
	}

	invalid := func(j, c []byte) {
		var sup ArtifactSupport
		assert.Error(t, sup.UnmarshalJSON(j))
		assert.Error(t, sup.UnmarshalCBOR(c))
	}

	j = []byte("\"source\"")
	c = []byte{0x66, 0x73, 0x6F, 0x75, 0x72, 0x63, 0x65}
	valid(j, c, ArtifactSupportSource)

	j = []byte("\"collected\"")
	c = []byte{0x69, 0x63, 0x6F, 0x6C, 0x6C, 0x65, 0x63, 0x74, 0x65, 0x64}
	valid(j, c, ArtifactSupportCollected)

	j = []byte("\"rims\"")
	c = []byte{0x64, 0x72, 0x69, 0x6D, 0x73}
	valid(j, c, ArtifactSupportRims)

	j = []byte("\"invalid\"")
	c = []byte{0x69, 0x6E, 0x76, 0x61, 0x6C, 0x69, 0x64}
	invalid(j, c)
}
