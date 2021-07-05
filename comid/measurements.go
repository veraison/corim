// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "fmt"

// Measurements is an array of Measurement
type Measurements []Measurement

// NewMeasurements instantiates an empty Measurements array
func NewMeasurements() *Measurements {
	return new(Measurements)
}

// AddMeasurements adds the supplied Measurement to the target Measurement
func (o *Measurements) AddMeasurement(m *Measurement) *Measurements {
	if o != nil && m != nil {
		*o = append(*o, *m)
	}
	return o
}

func (o Measurements) Valid() error {
	for i, m := range o {
		if err := m.Valid(); err != nil {
			return fmt.Errorf("measurement at index %d: %w", i, err)
		}
	}
	return nil
}
