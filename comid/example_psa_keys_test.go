// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "fmt"

func Example_psa_keys() {
	comid := Comid{}

	if err := comid.FromJSON([]byte(PSAKeysJSONTemplate)); err != nil {
		panic(err)
	}

	if err := comid.Valid(); err != nil {
		panic(err)
	}

	if err := extractKeys(&comid); err != nil {
		panic(err)
	}

	// output:
	// ImplementationID: 61636d652d696d706c656d656e746174696f6e2d69642d303030303030303031
	// InstanceID: 01ceebae7b8927a3227e5303cf5e0f1f7b34bb542ad7250ac03fbcde36ec2f1508
	// IAK public key: 2d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0a4d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a304441516344516741455731427671462b2f727938425761375a454d553178595948455138420a6c4c54344d46484f614f2b4943547449767245654570722f7366544150363648326843486462354845584b74524b6f6436514c634f4c504131513d3d0a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d
	// ImplementationID: 61636d652d696d706c656d656e746174696f6e2d69642d303030303030303031
	// InstanceID: 014ca3e4f50bf248c39787020d68ffd05c88767751bf2645ca923f57a98becd296
	// IAK public key: 2d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0a4d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a304441516344516741455731427671462b2f727938425761375a454d553178595948455138420a6c4c54344d46484f614f2b4943547449767245654570722f7366544150363648326843486462354845584b74524b6f6436514c634f4c504131513d3d0a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d
}

func extractKeys(c *Comid) error {
	if c.Triples.AttestVerifKeys == nil {
		return fmt.Errorf("no reference values triples")
	}

	for i, k := range *c.Triples.AttestVerifKeys {
		if err := extractPSAKey(k); err != nil {
			return fmt.Errorf("bad PSA verification key value at index %d: %w", i, err)
		}
	}

	return nil
}

func extractPSAKey(k AttestVerifKey) error {
	class := k.Environment.Class

	if err := extractImplementationID(class); err != nil {
		return fmt.Errorf("extracting impl-id: %w", err)
	}

	instance := k.Environment.Instance

	if err := extractInstanceID(instance); err != nil {
		return fmt.Errorf("extracting inst-id: %w", err)
	}

	if len(k.VerifKeys) != 1 {
		return fmt.Errorf("more than one key")
	}

	fmt.Printf("IAK public key: %x\n", k.VerifKeys[0])

	return nil
}

func extractInstanceID(i *Instance) error {
	if i == nil {
		return fmt.Errorf("no instance")
	}

	fmt.Printf("InstanceID: %x\n", i.Bytes())

	return nil
}
