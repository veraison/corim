// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "fmt"

func Example_psa_refval() {
	comid := Comid{}

	if err := comid.FromJSON([]byte(PSARefValJSONTemplate)); err != nil {
		panic(err)
	}

	if err := comid.Valid(); err != nil {
		panic(err)
	}

	if err := extractRefVals(&comid); err != nil {
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

}

func extractRefVals(c *Comid) error {
	if c.Triples.ReferenceValues == nil {
		return fmt.Errorf("no reference values triples")
	}

	for i, rv := range c.Triples.ReferenceValues.Values {
		if err := extractPSARefVal(rv); err != nil {
			return fmt.Errorf("bad PSA reference value at index %d: %w", i, err)
		}
	}

	return nil
}

func extractPSARefVal(rv ValueTriple) error {
	class := rv.Environment.Class

	if err := extractImplementationID(class); err != nil {
		return fmt.Errorf("extracting impl-id: %w", err)
	}

	measurements := rv.Measurements
	if err := extractSwMeasurements(measurements); err != nil {
		return fmt.Errorf("extracting measurements: %w", err)
	}

	return nil
}

func extractSwMeasurements(m Measurements) error {
	if len(m.Values) == 0 {
		return fmt.Errorf("no measurements")
	}
	for i, m := range m.Values {
		if err := extractSwMeasurement(m); err != nil {
			return fmt.Errorf("extracting measurement at index %d: %w", i, err)
		}
	}
	return nil
}

func extractSwMeasurement(m Measurement) error {
	if err := extractPSARefValID(m.Key); err != nil {
		return fmt.Errorf("extracting PSA refval id: %w", err)
	}

	if err := extractDigest(m.Val.Digests); err != nil {
		return fmt.Errorf("extracting digest: %w", err)
	}

	return nil
}

func extractDigest(d *Digests) error {
	if d == nil {
		return fmt.Errorf("no digest")
	}

	if len(*d) != 1 {
		return fmt.Errorf("more than one digest")
	}

	fmt.Printf("Digest: %x\n", (*d)[0].HashValue)

	return nil
}

func extractPSARefValID(k *Mkey) error {
	if k == nil {
		return fmt.Errorf("no measurement key")
	}

	id, ok := k.Value.(*TaggedPSARefValID)

	if !ok {
		return fmt.Errorf("expected PSA refval id, found: %T", k.Value)
	}

	fmt.Printf("SignerID: %x\n", id.SignerID)

	if id.Label != nil {
		fmt.Printf("Label: %s\n", *id.Label)
	}

	if id.Version != nil {
		fmt.Printf("Version: %s\n", *id.Version)
	}

	// ignore alg-id

	return nil
}

func extractImplementationID(c *Class) error {
	if c == nil {
		return fmt.Errorf("no class")
	}

	classID := c.ClassID

	if classID == nil {
		return fmt.Errorf("no class-id")
	}

	if classID.Type() != ImplIDType {
		return fmt.Errorf("class id is not a psa.impl-id")
	}

	fmt.Printf("ImplementationID: %x\n", classID.Bytes())

	return nil
}
