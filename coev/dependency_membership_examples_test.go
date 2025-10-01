// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"fmt"
	"log"

	"github.com/veraison/corim/comid"
)

func Example_encode_DependencyTriples() {
	coev := NewConciseEvidence()

	// Create a dependency triple showing that one domain depends on another
	dt := NewDependencyTriple()
	
	// Set the main domain
	domainClass := comid.NewClassUUID(TestUUID).
		SetVendor("ACME Ltd.").
		SetModel("Hypervisor").
		SetLayer(0)
	domainEnv := comid.Environment{Class: domainClass}
	dt.SetDomain(domainEnv)
	
	// Add dependent domains
	depClass1 := comid.NewClassUUID(TestUUID2).
		SetVendor("ACME Ltd.").
		SetModel("VM1").
		SetLayer(1)
	depEnv1 := comid.Environment{Class: depClass1}
	dt.AddDependentDomain(depEnv1)
	
	depClass2 := comid.NewClassUUID(TestUUID3).
		SetVendor("ACME Ltd.").
		SetModel("VM2").
		SetLayer(1)
	depEnv2 := comid.Environment{Class: depClass2}
	dt.AddDependentDomain(depEnv2)

	// Add to evidence triples
	evTriples := NewEvTriples().AddDependencyTriple(dt)
	
	err := coev.AddTriples(evTriples)
	if err != nil {
		log.Fatalf("could not add dependency triples: %v", err)
	}

	err = coev.AddEvidenceID(MustNewUUIDEvidenceID(TestUUID))
	if err != nil {
		log.Fatalf("could not add EvidenceID: %v", err)
	}

	cbor, err := coev.ToCBOR()
	if err != nil {
		log.Fatalf("could not encode to CBOR: %v", err)
	}

	fmt.Printf("Successfully encoded dependency triples (%d bytes)\n", len(cbor))
	// Output: Successfully encoded dependency triples (159 bytes)
}

func Example_encode_MembershipTriples() {
	coev := NewConciseEvidence()

	// Create a membership triple showing environments that belong to a domain
	mt := NewMembershipTriple()
	
	// Set the domain that contains the member environments
	domainClass := comid.NewClassUUID(TestUUID).
		SetVendor("ACME Ltd.").
		SetModel("TrustZone").
		SetLayer(0)
	domainEnv := comid.Environment{Class: domainClass}
	mt.SetDomain(domainEnv)
	
	// Add member environments
	memberClass1 := comid.NewClassUUID(TestUUID2).
		SetVendor("ACME Ltd.").
		SetModel("SecureWorld").
		SetLayer(1)
	memberEnv1 := comid.Environment{Class: memberClass1}
	mt.AddEnvironment(memberEnv1)
	
	memberClass2 := comid.NewClassUUID(TestUUID3).
		SetVendor("ACME Ltd.").
		SetModel("NormalWorld").
		SetLayer(1)
	memberEnv2 := comid.Environment{Class: memberClass2}
	mt.AddEnvironment(memberEnv2)

	// Add to evidence triples
	evTriples := NewEvTriples().AddMembershipTriple(mt)
	
	err := coev.AddTriples(evTriples)
	if err != nil {
		log.Fatalf("could not add membership triples: %v", err)
	}

	err = coev.AddEvidenceID(MustNewUUIDEvidenceID(TestUUID))
	if err != nil {
		log.Fatalf("could not add EvidenceID: %v", err)
	}

	cbor, err := coev.ToCBOR()
	if err != nil {
		log.Fatalf("could not encode to CBOR: %v", err)
	}

	fmt.Printf("Successfully encoded membership triples (%d bytes)\n", len(cbor))
	// Output: Successfully encoded membership triples (174 bytes)
}

func Example_encode_CombinedDependencyAndMembershipTriples() {
	coev := NewConciseEvidence()

	// Create both dependency and membership triples
	evTriples := NewEvTriples()
	
	// Add dependency triple
	dt := NewDependencyTriple()
	domainClass := comid.NewClassUUID(TestUUID).
		SetVendor("ACME Ltd.").
		SetModel("Platform")
	domainEnv := comid.Environment{Class: domainClass}
	dt.SetDomain(domainEnv)
	
	depClass := comid.NewClassUUID(TestUUID2).
		SetVendor("ACME Ltd.").
		SetModel("Application")
	depEnv := comid.Environment{Class: depClass}
	dt.AddDependentDomain(depEnv)
	
	evTriples.AddDependencyTriple(dt)
	
	// Add membership triple
	mt := NewMembershipTriple()
	mt.SetDomain(domainEnv)
	
	memberClass := comid.NewClassUUID(TestUUID3).
		SetVendor("ACME Ltd.").
		SetModel("Component")
	memberEnv := comid.Environment{Class: memberClass}
	mt.AddEnvironment(memberEnv)
	
	evTriples.AddMembershipTriple(mt)
	
	err := coev.AddTriples(evTriples)
	if err != nil {
		log.Fatalf("could not add triples: %v", err)
	}

	err = coev.AddEvidenceID(MustNewUUIDEvidenceID(TestUUID))
	if err != nil {
		log.Fatalf("could not add EvidenceID: %v", err)
	}

	cbor, err := coev.ToCBOR()
	if err != nil {
		log.Fatalf("could not encode to CBOR: %v", err)
	}

	fmt.Printf("Successfully encoded combined triples (%d bytes)\n", len(cbor))
	// Output: Successfully encoded combined triples (215 bytes)
}