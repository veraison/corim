// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "fmt"

func Example_cca_realm_refval() {
	comid := Comid{}

	if err := comid.FromJSON([]byte(CCARealmRefValJSONTemplate)); err != nil {
		panic(err)
	}

	if err := comid.Valid(); err != nil {
		panic(err)
	}

	if err := extractRealmRefVals(&comid); err != nil {
		panic(err)
	}

	// output:
	// ClassID: cd1f0e5526f9460db9d8f7fde171787c
	// Digest: 4284b5694ca6c0d2cf4789a0b95ac8025c818de52304364be7cd2981b2d2edc685b322277ec25819962413d8c9b2c1f5
	// Digest: 2107bbe761fca52d95136a1354db7a4dd57b1b26be0d3da71d9eb23986b34ba615abf6514cf35e5a9ea55a032d068a78
	// Digest: 4107bbe761fca52d95136a1354db7a4dd57b1b26be0d3da71d9eb23986b34ba615abf6514cf35e5a9ea55a032d068a78
	// Digest: 2507bbe761fca52d95136a1354db7a4dd57b1b26be0d3da71d9eb23986b34ba615abf6514cf35e5a9ea55a032d068a78
	// Digest: 3107bbe761fca52d95136a1354db7a4dd57b1b26be0d3da71d9eb23986b34ba615abf6514cf35e5a9ea55a032d068a78

}

func extractRealmRefVals(c *Comid) error {
	if c.Triples.ReferenceValues == nil {
		return fmt.Errorf("no reference values triples")
	}

	for i, rv := range *c.Triples.ReferenceValues {
		if err := extractRealmRefVal(rv); err != nil {
			return fmt.Errorf("bad Realm reference value at index %d: %w", i, err)
		}
	}

	return nil
}

func extractRealmRefVal(rv ReferenceValue) error {
	class := rv.Environment.Class

	if err := extractUuID(class); err != nil {
		return fmt.Errorf("extracting uuid: %w", err)
	}

	measurements := rv.Measurements

	if err := extractMeasurements(measurements); err != nil {
		return fmt.Errorf("extracting measurements: %w", err)
	}

	return nil
}

func extractMeasurements(m Measurements) error {
	if len(m) == 0 {
		return fmt.Errorf("no measurements")
	}

	for i, m := range m {
		if err := extractMeasurement(m); err != nil {
			return fmt.Errorf("extracting measurement at index %d: %w", i, err)
		}
	}

	return nil
}

func extractMeasurement(m Measurement) error {

	if err := extractDigest(m.Val.Digests); err != nil {
		return fmt.Errorf("extracting digest: %w", err)
	}

	return nil
}

func extractUuID(c *Class) error {
	if c == nil {
		return fmt.Errorf("no class")
	}

	classID := c.ClassID

	if classID == nil {
		return fmt.Errorf("no class-id")
	}

	if classID.Type() != "uuid" {
		return fmt.Errorf("class id is not a uuid")
	}

	fmt.Printf("ClassID: %x\n", classID.Bytes())

	return nil
}
