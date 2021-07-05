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
	// IAK public key: 4d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a30444151634451674145466e3074616f41775233506d724b6b594c74417344396f30354b534d366d6267664e436770754c306736567054486b5a6c3733776b354244786f56376e2b4f656565306949716b5733484d5a54334554696e694a64673d3d
	// ImplementationID: 61636d652d696d706c656d656e746174696f6e2d69642d303030303030303031
	// InstanceID: 014ca3e4f50bf248c39787020d68ffd05c88767751bf2645ca923f57a98becd296
	// IAK public key: 4d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a304441516344516741453656777165376879334f385970612b425545544c556a424e5533724558565579743958485237484a574c473758544b51643969316b565258654250444c466e66597275312f657578526e4a4d374839556f46444c64413d3d
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

	fmt.Printf("IAK public key: %x\n", k.VerifKeys[0].Key)

	return nil
}

func extractInstanceID(i *Instance) error {
	if i == nil {
		return fmt.Errorf("no instance")
	}

	instID, err := i.GetUEID()
	if err != nil {
		return fmt.Errorf("extracting implemenetation-id: %w", err)
	}

	fmt.Printf("InstanceID: %x\n", instID)

	return nil
}
