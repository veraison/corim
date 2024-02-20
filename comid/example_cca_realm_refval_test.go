// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"strings"
)

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
	// Vendor: Workload Client Ltd
	// ClassID: cd1f0e5526f9460db9d8f7fde171787c
	// InstanceID: 4284b5694ca6c0d2cf4789a0b95ac8025c818de52304364be7cd2981b2d2edc685b322277ec25819962413d8c9b2c1f5
	// Index: rim
	// Alg: sha-384
	// Digest: 4284b5694ca6c0d2cf4789a0b95ac8025c818de52304364be7cd2981b2d2edc685b322277ec25819962413d8c9b2c1f5
	// Index: rem
	// Alg: sha-384
	// Digest: 2107bbe761fca52d95136a1354db7a4dd57b1b26be0d3da71d9eb23986b34ba615abf6514cf35e5a9ea55a032d068a78
	// Alg: sha-384
	// Digest: 2507bbe761fca52d95136a1354db7a4dd57b1b26be0d3da71d9eb23986b34ba615abf6514cf35e5a9ea55a032d068a78
	// Alg: sha-384
	// Digest: 3107bbe761fca52d95136a1354db7a4dd57b1b26be0d3da71d9eb23986b34ba615abf6514cf35e5a9ea55a032d068a78
	// Alg: sha-384
	// Digest: 3507bbe761fca52d95136a1354db7a4dd57b1b26be0d3da71d9eb23986b34ba615abf6514cf35e5a9ea55a032d068a78

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
	instance := rv.Environment.Instance

	if err := extractRealmClass(class); err != nil {
		return fmt.Errorf("extracting class: %w", err)
	}

	if err := extractRealmInstanceID(instance); err != nil {
		return fmt.Errorf("extracting realm instanceID: %w", err)
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
	if err := extractIntegrityRegisters(m.Val.IntegrityRegisters); err != nil {
		return fmt.Errorf("extracting digest: %w", err)
	}

	return nil
}

func extractRealmClass(c *Class) error {
	if c == nil {
		fmt.Println("class not present")
		return nil
	}

	if c.Vendor != nil {
		fmt.Printf("Vendor: %s\n", c.GetVendor())
	}

	classID := c.ClassID
	if classID == nil {
		fmt.Println("class-id not present")
		return nil
	}

	if classID.Type() != "uuid" {
		return fmt.Errorf("class id is not a uuid")
	}
	if err := classID.Valid(); err != nil {
		return fmt.Errorf("invalid uuid: %v", err)
	}
	fmt.Printf("ClassID: %x\n", classID.Bytes())

	return nil
}

func extractRealmInstanceID(i *Instance) error {
	if i == nil {
		return fmt.Errorf("no instance")
	}

	if i.Type() != "bytes" {
		return fmt.Errorf("instance id is not bytes")
	}

	fmt.Printf("InstanceID: %x\n", i.Bytes())

	return nil
}

func extractIntegrityRegisters(r *IntegrityRegisters) error {
	if r == nil {
		return fmt.Errorf("no integrity registers")
	}

	keys, err := extractRegisterIndex(r)
	if err != nil {
		return fmt.Errorf("unable to extract register index: %v", err)
	}

	for _, k := range keys {
		d, ok := r.m[k]
		if !ok {
			return fmt.Errorf("unable to locate register index for: %s", k)
		}
		fmt.Printf("Index: %s\n", k)
		if err := extractRealmDigests(d); err != nil {
			return fmt.Errorf("invalid Digests for key: %s, %v", k, err)
		}
	}

	return nil
}

func extractRealmDigests(digests Digests) error {

	if err := digests.Valid(); err != nil {
		return fmt.Errorf("invalid digest: %v", err)
	}
	for _, d := range digests {
		fmt.Printf("Alg: %s\n", d.AlgIDToString())
		fmt.Printf("Digest: %x\n", d.HashValue)
	}

	return nil
}

func extractRegisterIndex(r *IntegrityRegisters) ([]string, error) {
	var keys [2]string
	for k := range r.m {
		switch t := k.(type) {
		case string:
			key := strings.ToLower(t)
			switch key {
			case "rim":
				keys[0] = key
			case "rem":
				keys[1] = key
			default:
				return nil, fmt.Errorf("unexpected register index: %s", key)
			}
		default:
			return nil, fmt.Errorf("unexpected type for index: %T", t)
		}
	}
	return keys[:], nil
}
