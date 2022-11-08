// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

// The tests in this file were copied from softwareidentity_test.go in github.com/veraison/swid with minor edits made
// to the expected results to reflect omission of the TagVersion field.

package cots

import (
	"reflect"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"github.com/veraison/swid"
)

func makeACMEEntityWithRoles(t *testing.T, roles ...interface{}) swid.Entity {
	e := swid.Entity{
		EntityName: "ACME Ltd",
		RegID:      "acme.example",
	}

	require.Nil(t, e.SetRoles(roles...))

	return e
}

func TestTag_RoundtripPSABundle(t *testing.T) {
	tv := AbbreviatedSwidTag{
		TagID:           swid.NewTagID("example.acme.roadrunner-sw-v1-0-0"),
		SoftwareName:    "Roadrunner software bundle",
		SoftwareVersion: "1.0.0",
		Entities: swid.Entities{
			makeACMEEntityWithRoles(t,
				swid.RoleTagCreator,
				swid.RoleSoftwareCreator,
				swid.RoleAggregator,
			),
		},
		Links: &swid.Links{
			swid.Link{
				Href: "example.acme.roadrunner-hw-v1-0-0",
				Rel:  *swid.NewRel("psa-rot-compound"),
			},
			swid.Link{
				Href: "example.acme.roadrunner-sw-bl-v1-0-0",
				Rel:  *swid.NewRel(swid.RelComponent),
			},
			swid.Link{
				Href: "example.acme.roadrunner-sw-prot-v1-0-0",
				Rel:  *swid.NewRel(swid.RelComponent),
			},
			swid.Link{
				Href: "example.acme.roadrunner-sw-arot-v1-0-0",
				Rel:  *swid.NewRel(swid.RelComponent),
			},
		},
	}
	/*
		a6                                      # map(6)
		   00                                   # unsigned(0)
		   78 21                                # text(33)
		      6578616d706c652e61636d652e726f616472756e6e65722d73772d76312d302d30 # "example.acme.roadrunner-sw-v1-0-0"
		   01                                   # unsigned(1)
		   78 1a                                # text(26)
		      526f616472756e6e657220736f6674776172652062756e646c65 # "Roadrunner software bundle"
		   0d                                   # unsigned(13)
		   65                                   # text(5)
		      312e302e30                        # "1.0.0"
		   02                                   # unsigned(2)
		   a3                                   # map(3)
		      18 1f                             # unsigned(31)
		      68                                # text(8)
		         41434d45204c7464               # "ACME Ltd"
		      18 20                             # unsigned(32)
		      6c                                # text(12)
		         61636d652e6578616d706c65       # "acme.example"
		      18 21                             # unsigned(33)
		      83                                # array(3)
		         01                             # unsigned(1)
		         02                             # unsigned(2)
		         03                             # unsigned(3)
		   04                                   # unsigned(4)
		   84                                   # array(4)
		      a2                                # map(2)
		         18 26                          # unsigned(38)
		         78 21                          # text(33)
		            6578616d706c652e61636d652e726f616472756e6e65722d68772d76312d302d30 # "example.acme.roadrunner-hw-v1-0-0"
		         18 28                          # unsigned(40)
		         70                             # text(16)
		            7073612d726f742d636f6d706f756e64 # "psa-rot-compound"
		      a2                                # map(2)
		         18 26                          # unsigned(38)
		         78 24                          # text(36)
		            6578616d706c652e61636d652e726f616472756e6e65722d73772d626c2d76312d302d30 # "example.acme.roadrunner-sw-bl-v1-0-0"
		         18 28                          # unsigned(40)
		         02                             # unsigned(2)
		      a2                                # map(2)
		         18 26                          # unsigned(38)
		         78 26                          # text(38)
		            6578616d706c652e61636d652e726f616472756e6e65722d73772d70726f742d76312d302d30 # "example.acme.roadrunner-sw-prot-v1-0-0"
		         18 28                          # unsigned(40)
		         02                             # unsigned(2)
		      a2                                # map(2)
		         18 26                          # unsigned(38)
		         78 26                          # text(38)
		            6578616d706c652e61636d652e726f616472756e6e65722d73772d61726f742d76312d302d30 # "example.acme.roadrunner-sw-arot-v1-0-0"
		         18 28                          # unsigned(40)
		         02                             # unsigned(2)
	*/
	// modified to not include TagVersion, which is omitted as empty
	expectedCBOR := []byte{
		/*0xa6,*/ 0xa5, 0x00, 0x78, 0x21, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e,
		0x61, 0x63, 0x6d, 0x65, 0x2e, 0x72, 0x6f, 0x61, 0x64, 0x72, 0x75, 0x6e,
		0x6e, 0x65, 0x72, 0x2d, 0x73, 0x77, 0x2d, 0x76, 0x31, 0x2d, 0x30, 0x2d,
		0x30 /*0x0c, 0x00,*/, 0x01, 0x78, 0x1a, 0x52, 0x6f, 0x61, 0x64, 0x72, 0x75,
		0x6e, 0x6e, 0x65, 0x72, 0x20, 0x73, 0x6f, 0x66, 0x74, 0x77, 0x61, 0x72,
		0x65, 0x20, 0x62, 0x75, 0x6e, 0x64, 0x6c, 0x65, 0x0d, 0x65, 0x31, 0x2e,
		0x30, 0x2e, 0x30, 0x02, 0xa3, 0x18, 0x1f, 0x68, 0x41, 0x43, 0x4d, 0x45,
		0x20, 0x4c, 0x74, 0x64, 0x18, 0x20, 0x6c, 0x61, 0x63, 0x6d, 0x65, 0x2e,
		0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x18, 0x21, 0x83, 0x01, 0x02,
		0x03, 0x04, 0x84, 0xa2, 0x18, 0x26, 0x78, 0x21, 0x65, 0x78, 0x61, 0x6d,
		0x70, 0x6c, 0x65, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x2e, 0x72, 0x6f, 0x61,
		0x64, 0x72, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x2d, 0x68, 0x77, 0x2d, 0x76,
		0x31, 0x2d, 0x30, 0x2d, 0x30, 0x18, 0x28, 0x70, 0x70, 0x73, 0x61, 0x2d,
		0x72, 0x6f, 0x74, 0x2d, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x75, 0x6e, 0x64,
		0xa2, 0x18, 0x26, 0x78, 0x24, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
		0x2e, 0x61, 0x63, 0x6d, 0x65, 0x2e, 0x72, 0x6f, 0x61, 0x64, 0x72, 0x75,
		0x6e, 0x6e, 0x65, 0x72, 0x2d, 0x73, 0x77, 0x2d, 0x62, 0x6c, 0x2d, 0x76,
		0x31, 0x2d, 0x30, 0x2d, 0x30, 0x18, 0x28, 0x02, 0xa2, 0x18, 0x26, 0x78,
		0x26, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x61, 0x63, 0x6d,
		0x65, 0x2e, 0x72, 0x6f, 0x61, 0x64, 0x72, 0x75, 0x6e, 0x6e, 0x65, 0x72,
		0x2d, 0x73, 0x77, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x2d, 0x76, 0x31, 0x2d,
		0x30, 0x2d, 0x30, 0x18, 0x28, 0x02, 0xa2, 0x18, 0x26, 0x78, 0x26, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x2e,
		0x72, 0x6f, 0x61, 0x64, 0x72, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x2d, 0x73,
		0x77, 0x2d, 0x61, 0x72, 0x6f, 0x74, 0x2d, 0x76, 0x31, 0x2d, 0x30, 0x2d,
		0x30, 0x18, 0x28, 0x02,
	}

	roundTripper(t, tv, expectedCBOR)
}

func TestTag_RoundtripPSAComponent(t *testing.T) {
	tv := AbbreviatedSwidTag{
		TagID:           swid.NewTagID("example.acme.roadrunner-sw-bl-v1-0-0"),
		SoftwareName:    "Roadrunner boot loader",
		SoftwareVersion: "1.0.0",
		Entities: swid.Entities{
			makeACMEEntityWithRoles(t,
				swid.RoleTagCreator,
				swid.RoleAggregator,
			),
		},
		Payloads: &swid.Payloads{
			swid.Payload{
				ResourceCollection: swid.ResourceCollection{
					Resources: &swid.Resources{
						swid.Resource{
							Type: swid.ResourceTypePSAMeasuredSoftwareComponent,
							ResourceExtension: swid.ResourceExtension{
								PSAMeasuredSoftwareComponent: swid.PSAMeasuredSoftwareComponent{
									MeasurementValue: swid.HashEntry{
										HashAlgID: 1, // sha-256
										HashValue: []byte("aabb...eeff"),
									},
									SignerID: swid.HashEntry{
										HashAlgID: 1, // sha-256
										HashValue: []byte("5192...1234"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	/*
		a6                                      # map(6)
		   00                                   # unsigned(0)
		   78 24                                # text(36)
		      6578616d706c652e61636d652e726f616472756e6e65722d73772d626c2d76312d302d30 # "example.acme.roadrunner-sw-bl-v1-0-0"
		   01                                   # unsigned(1)
		   76                                   # text(22)
		      526f616472756e6e657220626f6f74206c6f61646572 # "Roadrunner boot loader"
		   0d                                   # unsigned(13)
		   65                                   # text(5)
		      312e302e30                        # "1.0.0"
		   02                                   # unsigned(2)
		   a3                                   # map(3)
		      18 1f                             # unsigned(31)
		      68                                # text(8)
		         41434d45204c7464               # "ACME Ltd"
		      18 20                             # unsigned(32)
		      6c                                # text(12)
		         61636d652e6578616d706c65       # "acme.example"
		      18 21                             # unsigned(33)
		      82                                # array(2)
		         01                             # unsigned(1)
		         03                             # unsigned(3)
		   06                                   # unsigned(6)
		   a1                                   # map(1)
		      13                                # unsigned(19)
		      a3                                # map(3)
		         18 1d                          # unsigned(29)
		         78 24                          # text(36)
		            61726d2e636f6d2d5053414d65617375726564536f667477617265436f6d706f6e656e74 # "arm.com-PSAMeasuredSoftwareComponent"
		         78 1b                          # text(27)
		            61726d2e636f6d2d5053414d6561737572656d656e7456616c7565 # "arm.com-PSAMeasurementValue"
		         82                             # array(2)
		            01                          # unsigned(1)
		            4b                          # bytes(11)
		               616162622e2e2e65656666   # "aabb...eeff"
		         73                             # text(19)
		            61726d2e636f6d2d5053415369676e65724964 # "arm.com-PSASignerId"
		         82                             # array(2)
		            01                          # unsigned(1)
		            4b                          # bytes(11)
		               353139322e2e2e31323334   # "5192...1234"
	*/
	expectedCBOR := []byte{
		/*0xa6,*/ 0xa5, 0x00, 0x78, 0x24, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
		0x2e, 0x61, 0x63, 0x6d, 0x65, 0x2e, 0x72, 0x6f, 0x61, 0x64, 0x72,
		0x75, 0x6e, 0x6e, 0x65, 0x72, 0x2d, 0x73, 0x77, 0x2d, 0x62, 0x6c,
		0x2d, 0x76, 0x31, 0x2d, 0x30, 0x2d, 0x30 /*0x0c, 0x00,*/, 0x01, 0x76,
		0x52, 0x6f, 0x61, 0x64, 0x72, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x20,
		0x62, 0x6f, 0x6f, 0x74, 0x20, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72,
		0x0d, 0x65, 0x31, 0x2e, 0x30, 0x2e, 0x30, 0x02, 0xa3, 0x18, 0x1f,
		0x68, 0x41, 0x43, 0x4d, 0x45, 0x20, 0x4c, 0x74, 0x64, 0x18, 0x20,
		0x6c, 0x61, 0x63, 0x6d, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70,
		0x6c, 0x65, 0x18, 0x21, 0x82, 0x01, 0x03, 0x06, 0xa1, 0x13, 0xa3,
		0x18, 0x1d, 0x78, 0x24, 0x61, 0x72, 0x6d, 0x2e, 0x63, 0x6f, 0x6d,
		0x2d, 0x50, 0x53, 0x41, 0x4d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65,
		0x64, 0x53, 0x6f, 0x66, 0x74, 0x77, 0x61, 0x72, 0x65, 0x43, 0x6f,
		0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x78, 0x1b, 0x61, 0x72,
		0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x2d, 0x50, 0x53, 0x41, 0x4d, 0x65,
		0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x56, 0x61,
		0x6c, 0x75, 0x65, 0x82, 0x01, 0x4b, 0x61, 0x61, 0x62, 0x62, 0x2e,
		0x2e, 0x2e, 0x65, 0x65, 0x66, 0x66, 0x73, 0x61, 0x72, 0x6d, 0x2e,
		0x63, 0x6f, 0x6d, 0x2d, 0x50, 0x53, 0x41, 0x53, 0x69, 0x67, 0x6e,
		0x65, 0x72, 0x49, 0x64, 0x82, 0x01, 0x4b, 0x35, 0x31, 0x39, 0x32,
		0x2e, 0x2e, 0x2e, 0x31, 0x32, 0x33, 0x34,
	}

	roundTripper(t, tv, expectedCBOR)
}

// marshal + unmarshal
func roundTripper(t *testing.T, tv interface{}, expectedCBOR []byte) interface{} {
	encMode, err := cbor.EncOptions{TimeTag: cbor.EncTagRequired}.EncMode()
	require.Nil(t, err)

	data, err := encMode.Marshal(tv)

	assert.Nil(t, err)
	t.Logf("CBOR(hex): %x\n", data)
	assert.Equal(t, expectedCBOR, data)

	decMode, err := cbor.DecOptions{TimeTag: cbor.DecTagOptional}.DecMode()
	require.Nil(t, err)

	actual := reflect.New(reflect.TypeOf(tv))
	err = decMode.Unmarshal(data, actual.Interface())

	assert.Nil(t, err)
	assert.Equal(t, tv, actual.Elem().Interface())

	// Return the an interface wrapping the roundtripped test vector.
	// In case it's needed for further processing it can be extracted
	// with a type assertion.
	return actual.Elem().Interface()
}
