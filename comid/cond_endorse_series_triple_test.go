// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/extensions"
)

func Test_CondEndorseSeriesTriple_NewCondEndorseSeriesTriples_OK(t *testing.T) {
	c := NewCondEndorseSeriesTriples()
	require.NotNil(t, c)
	assert.True(t, c.IsEmpty())
}

func Test_CondEndorseSeriesTriple_Add_Valid(t *testing.T) {
	c := NewCondEndorseSeriesTriples()
	// Create a valid triple
	triple := &CondEndorseSeriesTriple{
		Condition: CondEndorseSeriesCondition{
			Environment: Environment{
				Class: &Class{
					ClassID: MustNewUUIDClassID(TestUUID),
				},
			},
		},
		Series: *NewCondEndorseSeriesRecords(),
	}

	c.Add(triple)
	assert.False(t, c.IsEmpty())
	err := c.Valid()
	assert.ErrorContains(t, err, "empty conditional series")
}

func Test_CondEndorseSeriesTriple_Valid_EmptyCondition(t *testing.T) {
	ces := NewCondEndorseSeriesTriples()
	ces_triple := &CondEndorseSeriesTriple{
		Series: *NewCondEndorseSeriesRecords(),
	}
	ces.Add(ces_triple)
	err := ces.Valid()
	assert.EqualError(t, err, "error at index 0: condition validation failed: environment validation failed: environment must not be empty")
}

func Test_CondEndorseSeriesTriple_Valid_InvalidCondition(t *testing.T) {
	c := NewCondEndorseSeriesTriples()
	series := &CondEndorseSeriesTriple{
		Condition: CondEndorseSeriesCondition{
			Environment: Environment{
				Class: &Class{
					ClassID: MustNewUUIDClassID(TestUUID),
				},
			},
			Measurements: *NewMeasurements().Add(
				&Measurement{
					Key: MustNewMkey(uint64(123), UintType),
				},
			),
		},
		Series: *NewCondEndorseSeriesRecords(),
	}
	c.Add(series)
	err := c.Valid()
	assert.ErrorContains(t, err, "no measurement value set")

}

func Test_CondEndorseSeriesTriple_Valid_InvalidSeries(t *testing.T) {
	c := NewCondEndorseSeriesTriples()
	// Create series with invalid record
	series := &CondEndorseSeriesTriple{
		Condition: CondEndorseSeriesCondition{
			Environment: Environment{
				Class: &Class{
					ClassID: MustNewUUIDClassID(TestUUID),
				},
			},
		},
		Series: *NewCondEndorseSeriesRecords(),
	}
	// Add invalid record to series
	invalidRecord := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(123), UintType),
			},
		),
	}
	series.Series.Add(invalidRecord)
	c.Add(series)
	err := c.Valid()
	assert.ErrorContains(t, err, "no measurement value")
}

func Test_CondEndorseSeriesTriple_Valid_EmptySeries(t *testing.T) {
	c := NewCondEndorseSeriesTriples()
	series := &CondEndorseSeriesTriple{
		Condition: CondEndorseSeriesCondition{
			Environment: Environment{
				Class: &Class{
					ClassID: MustNewUUIDClassID(TestUUID),
				},
			},
		},
		Series: *NewCondEndorseSeriesRecords(),
	}
	c.Add(series)
	err := c.Valid()
	assert.ErrorContains(t, err, "empty conditional series")
}

func Test_CondEndorseSeriesTriple_Valid_ValidSeries(t *testing.T) {
	c := NewCondEndorseSeriesTriples()

	// Create valid series record
	validRecord := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(123), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("test-value")),
				},
			},
		),
		Addition: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(456), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("add-value")),
				},
			},
		),
	}

	series := &CondEndorseSeriesTriple{
		Condition: CondEndorseSeriesCondition{
			Environment: Environment{
				Class: &Class{
					ClassID: MustNewUUIDClassID(TestUUID),
				},
			},
			Measurements: *NewMeasurements().Add(
				&Measurement{
					Key: MustNewMkey(uint64(789), UintType),
					Val: Mval{
						RawValue: NewRawValue().SetBytes([]byte("condition-value")),
					},
				},
			),
		},
		Series: *NewCondEndorseSeriesRecords(),
	}
	series.Series.Add(validRecord)
	c.Add(series)
	err := c.Valid()
	assert.NoError(t, err)
}

type testExtensions struct {
	TestSVN uint `cbor:"-72,keyasint,omitempty" json:"testsvn,omitempty"`
}

func Test_CondEndorseSeriesTriple_RegisterExtensions_OK(t *testing.T) {
	extMap := extensions.NewMap().
		Add(ExtMval, &testExtensions{})
	series := &CondEndorseSeriesTriple{
		Condition: CondEndorseSeriesCondition{
			Environment: Environment{
				Class: &Class{
					ClassID: MustNewUUIDClassID(TestUUID),
				},
			},
		},
		Series: *NewCondEndorseSeriesRecords(),
	}
	err := series.RegisterExtensions(extMap)
	require.NoError(t, err)
}

func Test_CondEndorseSeriesTriple_RegisterExtensions_NOK(t *testing.T) {
	expectedErr := `condition: unexpected extension point: "ReferenceValue"`
	extMap := extensions.NewMap().
		Add(ExtReferenceValue, &testExtensions{})
	series := &CondEndorseSeriesTriple{
		Condition: CondEndorseSeriesCondition{
			Environment: Environment{
				Class: &Class{
					ClassID: MustNewUUIDClassID(TestUUID),
				},
			},
		},
		Series: *NewCondEndorseSeriesRecords(),
	}
	err := series.RegisterExtensions(extMap)
	assert.EqualError(t, err, expectedErr)
}

func Test_CondEndorseSeriesTriple_RegisterExtensions_SeriesError(t *testing.T) {
	extMap := extensions.NewMap().
		Add(ExtReferenceValue, &testExtensions{})

	// Create series with record that will cause extension registration error
	series := &CondEndorseSeriesTriple{
		Condition: CondEndorseSeriesCondition{
			Environment: Environment{
				Class: &Class{
					ClassID: MustNewUUIDClassID(TestUUID),
				},
			},
		},
		Series: *NewCondEndorseSeriesRecords(),
	}

	// Add record with measurements that will cause extension error
	record := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(123), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("test-value")),
				},
			},
		),
	}
	series.Series.Add(record)

	err := series.RegisterExtensions(extMap)
	assert.ErrorContains(t, err, "unexpected extension point")
}

func Test_CondEndorseSeriesRecord_RegisterExtensions_OK(t *testing.T) {
	record := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(123), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("test-value")),
				},
			},
		),
		Addition: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(456), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("add-value")),
				},
			},
		),
	}

	extMap := extensions.NewMap().
		Add(ExtMval, &testExtensions{})

	err := record.RegisterExtensions(extMap)
	require.NoError(t, err)
}

func Test_CondEndorseSeriesRecord_RegisterExtensions_SelectionError(t *testing.T) {
	extMap := extensions.NewMap().
		Add(ExtReferenceValue, &testExtensions{})

	record := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(123), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("test-value")),
				},
			},
		),
	}

	err := record.RegisterExtensions(extMap)
	assert.ErrorContains(t, err, "selection:")
}

func Test_CondEndorseSeriesRecord_RegisterExtensions_AdditionError(t *testing.T) {
	extMap := extensions.NewMap().
		Add(ExtReferenceValue, &testExtensions{})

	record := &CondEndorseSeriesRecord{
		Addition: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(456), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("test-value")),
				},
			},
		),
	}

	err := record.RegisterExtensions(extMap)
	assert.ErrorContains(t, err, "unexpected extension point")
}

func Test_CondEndorseSeriesRecords_RegisterExtensions_OK(t *testing.T) {
	records := NewCondEndorseSeriesRecords()

	record := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(123), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("test-value")),
				},
			},
		),
		Addition: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(456), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("add-value")),
				},
			},
		),
	}

	records.Add(record)

	extMap := extensions.NewMap().
		Add(ExtMval, &testExtensions{})

	err := records.RegisterExtensions(extMap)
	require.NoError(t, err)
}

func Test_CondEndorseSeriesRecords_RegisterExtensions_Error(t *testing.T) {
	records := NewCondEndorseSeriesRecords()

	record := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(123), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("test-value")),
				},
			},
		),
		Addition: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(456), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("add-value")),
				},
			},
		),
	}

	records.Add(record)

	extMap := extensions.NewMap().
		Add(ExtReferenceValue, &testExtensions{})

	err := records.RegisterExtensions(extMap)
	assert.ErrorContains(t, err, "unexpected extension point")
}

func Test_CondEndorseSeriesTriple_GetExtensions(t *testing.T) {
	series := &CondEndorseSeriesTriple{
		Condition: CondEndorseSeriesCondition{
			Environment: Environment{
				Class: &Class{
					ClassID: MustNewUUIDClassID(TestUUID),
				},
			},
		},
		Series: *NewCondEndorseSeriesRecords(),
	}
	exts := series.GetExtensions()
	assert.Nil(t, exts)
}

func Test_CondEndorseSeriesTriple_MarshalJSON(t *testing.T) {
	c := NewCondEndorseSeriesTriples()

	// Create valid triple
	triple := &CondEndorseSeriesTriple{
		Condition: CondEndorseSeriesCondition{
			Environment: Environment{
				Class: &Class{
					ClassID: MustNewUUIDClassID(TestUUID),
				},
			},
		},
		Series: *NewCondEndorseSeriesRecords(),
	}
	c.Add(triple)

	data, err := c.MarshalJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Verify it's valid JSON
	var result []interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)
}

func Test_CondEndorseSeriesTriple_UnmarshalJSON(t *testing.T) {
	c := NewCondEndorseSeriesTriples()
	// Create invalid JSON data that should cause an unmarshaling error
	jsonData := `{"invalid": "json", "for": "collection"}`
	err := c.UnmarshalJSON([]byte(jsonData))
	assert.ErrorContains(t, err, "cannot unmarshal object into Go value of type []json.RawMessage")
	assert.True(t, c.IsEmpty(), "collection should remain empty after failed unmarshaling")
}

func Test_CondEndorseSeriesTriple_MarshalCBOR(t *testing.T) {
	c := NewCondEndorseSeriesTriples()

	// Create valid triple
	triple := &CondEndorseSeriesTriple{
		Condition: CondEndorseSeriesCondition{
			Environment: Environment{
				Class: &Class{
					ClassID: MustNewUUIDClassID(TestUUID),
				},
			},
		},
		Series: *NewCondEndorseSeriesRecords(),
	}
	c.Add(triple)

	data, err := c.MarshalCBOR()
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}

func Test_CondEndorseSeriesTriple_UnmarshalCBOR(t *testing.T) {
	c := NewCondEndorseSeriesTriples()

	// First marshal a valid triple to CBOR
	triple := &CondEndorseSeriesTriple{
		Condition: CondEndorseSeriesCondition{
			Environment: Environment{
				Class: &Class{
					ClassID: MustNewUUIDClassID(TestUUID),
				},
			},
		},
		Series: *NewCondEndorseSeriesRecords(),
	}
	c.Add(triple)

	cborData, err := c.MarshalCBOR()
	require.NoError(t, err)

	// Now unmarshal it
	c = NewCondEndorseSeriesTriples()
	err = c.UnmarshalCBOR(cborData)
	require.NoError(t, err)
	assert.False(t, c.IsEmpty())
}

func Test_CondEndorseSeriesRecord_Valid_OK(t *testing.T) {
	record := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(123), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("test-value")),
				},
			},
		),
		Addition: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(456), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("add-value")),
				},
			},
		),
	}

	assert.NoError(t, record.Valid())
}

func Test_CondEndorseSeriesRecord_Valid_EmptySelections(t *testing.T) {
	record := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements(),
		Addition: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(456), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("add-value")),
				},
			},
		),
	}

	assert.NoError(t, record.Valid())
}

func Test_CondEndorseSeriesRecord_Valid_EmptyAdditions(t *testing.T) {
	record := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(123), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("test-value")),
				},
			},
		),
		Addition: *NewMeasurements(),
	}

	assert.NoError(t, record.Valid())
}

func Test_CondEndorseSeriesRecord_Valid_InvalidSelection(t *testing.T) {
	record := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(123), UintType),
			},
		),
	}

	err := record.Valid()
	assert.ErrorContains(t, err, "no measurement value set")

}

func Test_CondEndorseSeriesRecord_Valid_InvalidAddition(t *testing.T) {
	record := &CondEndorseSeriesRecord{
		Addition: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(456), UintType),
			},
		),
	}

	err := record.Valid()
	assert.ErrorContains(t, err, "no measurement value set")

}

func Test_CondEndorseSeriesRecord_GetExtensions(t *testing.T) {
	record := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(123), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("test-value")),
				},
			},
		),
		Addition: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(456), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("add-value")),
				},
			},
		),
	}
	exts := record.GetExtensions()
	assert.Nil(t, exts)
}

func Test_CondEndorseSeriesRecords_NewCondEndorseSeriesRecords(t *testing.T) {
	records := NewCondEndorseSeriesRecords()
	require.NotNil(t, records)
	assert.True(t, records.IsEmpty())
}

func Test_CondEndorseSeriesRecords_Add_Valid(t *testing.T) {
	records := NewCondEndorseSeriesRecords()

	record := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(123), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("test-value")),
				},
			},
		),
		Addition: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(456), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("add-value")),
				},
			},
		),
	}

	result := records.Add(record)
	assert.NotNil(t, result)
	assert.False(t, records.IsEmpty())
	assert.NoError(t, records.Valid())
}

func Test_CondEndorseSeriesRecords_Valid_InvalidRecord(t *testing.T) {
	records := NewCondEndorseSeriesRecords()

	invalidRecord := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(123), UintType),
			},
		),
		Addition: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(456), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("add-value")),
				},
			},
		),
	}

	records.Add(invalidRecord)
	err := records.Valid()
	assert.ErrorContains(t, err, "error at index 0:")
}

func Test_CondEndorseSeriesRecords_GetExtensions(t *testing.T) {
	records := NewCondEndorseSeriesRecords()

	record := &CondEndorseSeriesRecord{
		Selection: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(123), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("test-value")),
				},
			},
		),
		Addition: *NewMeasurements().Add(
			&Measurement{
				Key: MustNewMkey(uint64(456), UintType),
				Val: Mval{
					RawValue: NewRawValue().SetBytes([]byte("add-value")),
				},
			},
		),
	}

	records.Add(record)

	exts := records.GetExtensions()
	assert.Nil(t, exts)
}

func Test_CondEndorseSeriesCondition_serialize_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		condition    CondEndorseSeriesCondition
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title: "minimal",
			condition: CondEndorseSeriesCondition{
				Environment: Environment{
					Class: &Class{
						ClassID: MustNewUUIDClassID(TestUUID),
					},
				},
			},
			expectedCBOR: []byte{
				0x82,       // array(2) [condition]
				0xa1,       // . [0]map(1) [environment]
				0x0,        // . . key: 0 [class]
				0xa1,       // . . value: map(1) [class]
				0x0,        // . . . key: 0 [class-id]
				0xd8, 0x25, // . . . value: tag(37) [uuid]
				0x50, //       . . . . bstr(16)
				0x31, 0xfb, 0x5a, 0xbf, 0x02, 0x3e, 0x49, 0x92,
				0xaa, 0x4e, 0x95, 0xf9, 0xc1, 0x50, 0x3b, 0xfa,
				0x80, //       . [1]array(0) [measurements]
			},
			expectedJSON: `{"environment":{"class":{"id":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}}},"measurements":null}`,
		},
		{
			title: "all fields",
			condition: CondEndorseSeriesCondition{
				Environment: Environment{
					Class: &Class{
						ClassID: MustNewUUIDClassID(TestUUID),
					},
				},
				Measurements: *NewMeasurements().Add((&Measurement{}).SetSVN(5)),
				AuthorizedBy: NewCryptoKeys().Add(MustNewPKIXBase64Key(TestECPubKey)),
			},
			expectedCBOR: []byte{
				0x83,       //       array(3) [condition]
				0xa1,       //       . [0]map(1) [environment]
				0x0,        //       . . key: 0 [class]
				0xa1,       //       . . value: map(1) [class]
				0x0,        //       . . . key: 0 [class-id]
				0xd8, 0x25, //       . . . value: tag(37) [uuid]
				0x50, //             . . . . bstr(16)
				0x31, 0xfb, 0x5a, 0xbf, 0x02, 0x3e, 0x49, 0x92,
				0xaa, 0x4e, 0x95, 0xf9, 0xc1, 0x50, 0x3b, 0xfa,
				0x81,             // . [1]array(0) [measurements]
				0xa1,             // . . [0]map(1) [measurement]
				0x01,             // . . . key: 1 [val]
				0xa1,             // . . . value: map(1) [mval]
				0x01,             // . . . . key: 1 [svn]
				0xd9, 0x02, 0x28, // . . . . value: tag(552) [svn]
				0x05,             // . . . . . 5
				0x81,             // . [2]array(1) [authorized-by]
				0xd9, 0x02, 0x2a, // . . [0]tag(554) [pkix-base64-key]
				0x78, 0xb1, //       . . . tstr(177)
				0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x42, 0x45, 0x47,
				0x49, 0x4e, 0x20, 0x50, 0x55, 0x42, 0x4c, 0x49,
				0x43, 0x20, 0x4b, 0x45, 0x59, 0x2d, 0x2d, 0x2d,
				0x2d, 0x2d, 0x0a, 0x4d, 0x46, 0x6b, 0x77, 0x45, // 32

				0x77, 0x59, 0x48, 0x4b, 0x6f, 0x5a, 0x49, 0x7a,
				0x6a, 0x30, 0x43, 0x41, 0x51, 0x59, 0x49, 0x4b,
				0x6f, 0x5a, 0x49, 0x7a, 0x6a, 0x30, 0x44, 0x41,
				0x51, 0x63, 0x44, 0x51, 0x67, 0x41, 0x45, 0x57, // 64

				0x31, 0x42, 0x76, 0x71, 0x46, 0x2b, 0x2f, 0x72,
				0x79, 0x38, 0x42, 0x57, 0x61, 0x37, 0x5a, 0x45,
				0x4d, 0x55, 0x31, 0x78, 0x59, 0x59, 0x48, 0x45,
				0x51, 0x38, 0x42, 0x0a, 0x6c, 0x4c, 0x54, 0x34, // 96

				0x4d, 0x46, 0x48, 0x4f, 0x61, 0x4f, 0x2b, 0x49,
				0x43, 0x54, 0x74, 0x49, 0x76, 0x72, 0x45, 0x65,
				0x45, 0x70, 0x72, 0x2f, 0x73, 0x66, 0x54, 0x41,
				0x50, 0x36, 0x36, 0x48, 0x32, 0x68, 0x43, 0x48, // 128

				0x64, 0x62, 0x35, 0x48, 0x45, 0x58, 0x4b, 0x74,
				0x52, 0x4b, 0x6f, 0x64, 0x36, 0x51, 0x4c, 0x63,
				0x4f, 0x4c, 0x50, 0x41, 0x31, 0x51, 0x3d, 0x3d,
				0x0a, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x45, 0x4e, // 160

				0x44, 0x20, 0x50, 0x55, 0x42, 0x4c, 0x49, 0x43,
				0x20, 0x4b, 0x45, 0x59, 0x2d, 0x2d, 0x2d, 0x2d,
				0x2d, // 177
			},
			expectedJSON: `{"environment":{"class":{"id":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}}},"measurements":[{"value":{"svn":{"type":"exact-value","value":5}}}],"authorized-by":[{"type":"pkix-base64-key","value":"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"}]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			bytes, err := em.Marshal(&tc.condition)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, bytes)

			var decoded CondEndorseSeriesCondition
			err = dm.Unmarshal(bytes, &decoded)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.condition.Environment, decoded.Environment)

			bytes, err = json.Marshal(&tc.condition)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedJSON, string(bytes))

			err = json.Unmarshal(bytes, &decoded)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.condition.Environment, decoded.Environment)
		})
	}
}
