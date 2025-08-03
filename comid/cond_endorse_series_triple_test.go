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
		Condition: ValueTriple{
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
	assert.Contains(t, err.Error(), "no measurement entries")
}

func Test_CondEndorseSeriesTriple_Valid_EmptyCondition(t *testing.T) {
	ces := NewCondEndorseSeriesTriples()
	ces_triple := &CondEndorseSeriesTriple{
		Series: *NewCondEndorseSeriesRecords(),
	}
	ces.Add(ces_triple)
	err := ces.Valid()
	assert.EqualError(t, err, "error at index 0: stateful environment validation failed: environment validation failed: environment must not be empty")
}

func Test_CondEndorseSeriesTriple_Valid_InvalidCondition(t *testing.T) {
	c := NewCondEndorseSeriesTriples()
	series := &CondEndorseSeriesTriple{
		Condition: ValueTriple{
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
		Condition: ValueTriple{
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
	assert.ErrorContains(t, err, "no measurement entries")
}

func Test_CondEndorseSeriesTriple_Valid_EmptySeries(t *testing.T) {
	c := NewCondEndorseSeriesTriples()
	series := &CondEndorseSeriesTriple{
		Condition: ValueTriple{
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
	assert.ErrorContains(t, err, "no measurement entries")
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
		Condition: ValueTriple{
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
		Condition: ValueTriple{
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
		Condition: ValueTriple{
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
		Condition: ValueTriple{
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
		Condition: ValueTriple{
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
		Condition: ValueTriple{
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
		Condition: ValueTriple{
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
		Condition: ValueTriple{
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
