// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "fmt"

func Example_cca_refval() {
	comid := Comid{}

	if err := comid.FromJSON([]byte(CCARefValJSONTemplate)); err != nil {
		panic(err)
	}

	if err := comid.Valid(); err != nil {
		panic(err)
	}

	if err := extractCcaRefVals(&comid); err != nil {
		panic(err)
	}

	// output:
	// ImplementationID: 61636d652d696d706c656d656e746174696f6e2d69642d303030303030303031
	// SignerID: acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b
	// Label: BL
	// Version: 2.1.0
	// Digest: 87428fc522803d31065e7bce3cf03fe475096631e5e07bbd7a0fde60c4cf25c7
	// SignerID: acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b
	// Label: PRoT
	// Version: 1.3.5
	// Digest: 0263829989b6fd954f72baaf2fc64bc2e2f01d692d4de72986ea808f6e99813f
	// SignerID: acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b
	// Label: ARoT
	// Version: 0.1.4
	// Digest: a3a5e715f0cc574a73c3f9bebb6bc24f32ffd5b67b387244c2c909da779a1478
	// Label: cca-platform-config
	// Raw value: {72617776616c75650a72617776616c75650a}
}

func extractCcaRefVals(c *Comid) error {
	if c.Triples.ReferenceValues == nil {
		return fmt.Errorf("no reference values triples")
	}

	for i, rv := range *c.Triples.ReferenceValues {
		if err := extractCCARefVal(rv); err != nil {
			return fmt.Errorf("bad PSA reference value at index %d: %w", i, err)
		}
	}

	return nil
}

func extractCCARefVal(rv ReferenceValue) error {
	class := rv.Environment.Class

	if err := extractImplementationID(class); err != nil {
		return fmt.Errorf("extracting impl-id: %w", err)
	}

	for i, m := range rv.Measurements {
		if m.Key == nil {
			return fmt.Errorf("missing mKey at index %d", i)
		}
		if !m.Key.IsSet() {
			return fmt.Errorf("mKey not set at index %d", i)
		}
		if m.Key.IsPSARefValID() {
			if err := extractSwMeasurement(m); err != nil {
				return fmt.Errorf("extracting measurement at index %d: %w", i, err)
			}
		}
		if m.Key.IsCCARefValID() {
			extractCCARefValID(m.Key)
			if err := extractRawValue(m.Val.RawValue); err != nil {
				return fmt.Errorf("extracting raw vlue: %w", err)
			}

			return nil
		}
	}

	return nil
}

func extractRawValue(r *RawValue) error {
	if r == nil {
		return fmt.Errorf("no raw value")
	}

	fmt.Printf("Raw value: %x\n", *r)

	return nil
}

func extractCCARefValID(k *Mkey) error {
	if k == nil {
		return fmt.Errorf("no measurement key")
	}

	id, err := k.GetCCARefValID()
	if err != nil {
		return fmt.Errorf("getting CCA refval id: %w", err)
	}

	if id.Label != nil {
		fmt.Printf("Label: %s\n", *id.Label)
	}
	return nil
}
