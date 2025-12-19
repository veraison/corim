// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/eat"
)

// Example demonstrates the PSA profile functionality
// according to draft-fdb-rats-psa-endorsements-08
func Example() {
	// The PSA profile should be automatically registered via init()
	profileID, err := eat.NewProfile(PSAProfileURI)
	if err != nil {
		panic(err)
	}

	profileURI, err := profileID.Get()
	if err != nil {
		panic(err)
	}

	fmt.Printf("PSA Profile ID: %s\n", profileURI)
	fmt.Println("PSA Profile registered successfully")

	// Output:
	// PSA Profile ID: tag:arm.com,2025:psa#1.0.0
	// PSA Profile registered successfully
}

// ExamplePSASwComponentMeasurementValues demonstrates PSA software component validation
func ExamplePSASwComponentMeasurementValues() {
	// Create a valid PSA software component measurement
	values := PSASwComponentMeasurementValues{
		Digests: []PSADigest{
			{
				Algorithm: "sha-256",
				Value:     make([]byte, 32), // 32-byte digest
			},
		},
		CryptoKeys: [][]byte{make([]byte, 32)}, // 32-byte signer ID
	}

	err := values.Valid()
	if err != nil {
		panic(err)
	}

	fmt.Println("PSA software component measurement is valid")

	// Output:
	// PSA software component measurement is valid
}

// ExamplePSASoftwareComponentKeyType demonstrates using PSA software component keys
func ExamplePSASoftwareComponentKeyType() {
	// Create a PSA software component key
	key, err := comid.NewMkey(PSASoftwareComponentType, PSASoftwareComponentType)
	if err != nil {
		panic(err)
	}

	fmt.Printf("PSA software component key type: %s\n", key.Value.Type())

	// Output:
	// PSA software component key type: psa.software-component
}

// ExamplePSASwRelationship demonstrates PSA software relationships
func ExamplePSASwRelationship() {
	// Create simple measurements using uint keys
	oldMeasurement, err := comid.NewUintMeasurement(uint64(1))
	if err != nil {
		panic(err)
	}
	
	newMeasurement, err := comid.NewUintMeasurement(uint64(2))
	if err != nil {
		panic(err)
	}

	// Create update relationship
	updateRel, err := NewPSAUpdateRelationship(newMeasurement, oldMeasurement, true)
	if err != nil {
		panic(err)
	}

	err = updateRel.Valid()
	if err != nil {
		panic(err)
	}

	fmt.Printf("PSA update relationship type: %d\n", updateRel.Relation.Type)
	fmt.Printf("Security critical: %t\n", updateRel.Relation.SecurityCritical)

	// Output:
	// PSA update relationship type: 1
	// Security critical: true
}

// ExamplePSACertNumType demonstrates PSA certificate number validation
func ExamplePSACertNumType() {
	// Create a valid PSA certificate number
	certNum := PSACertNumType("1234567890123 - 56789")

	err := certNum.Valid()
	if err != nil {
		panic(err)
	}

	fmt.Printf("PSA certificate number is valid: %s\n", string(certNum))

	// Output:
	// PSA certificate number is valid: 1234567890123 - 56789
}

// ExamplePSASwRelTriples demonstrates PSA software relationship triples
func ExamplePSASwRelTriples() {
	// Create simple measurements using uint keys
	oldMeasurement, err := comid.NewUintMeasurement(uint64(1))
	if err != nil {
		panic(err)
	}
	
	newMeasurement, err := comid.NewUintMeasurement(uint64(2))
	if err != nil {
		panic(err)
	}

	// Create patch relationship
	patchRel, err := NewPSAPatchRelationship(newMeasurement, oldMeasurement, false)
	if err != nil {
		panic(err)
	}

	// Create environment with minimal content
	vendor := "PSA Example Vendor"
	env := comid.Environment{
		Class: &comid.Class{
			Vendor: &vendor,
		},
	}

	// Create triple
	triple := PSASwRelTriple{
		Environment:  env,
		Relationship: patchRel,
	}

	// Create triples collection
	triples := NewPSASwRelTriples()
	triples.Add(triple)

	err = triples.Valid()
	if err != nil {
		panic(err)
	}

	fmt.Printf("PSA software relationship triples count: %d\n", len(triples.Values))
	fmt.Printf("Relationship type: %d (patches)\n", triples.Values[0].Relationship.Relation.Type)

	// Output:
	// PSA software relationship triples count: 1
	// Relationship type: 2 (patches)
}
