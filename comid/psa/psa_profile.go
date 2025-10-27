// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"fmt"
	
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/corim"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
)

var ProfileID *eat.Profile

// PSAProfile defines the PSA endorsements profile as specified in
// draft-fdb-rats-psa-endorsements-08
const PSAProfileURI = "tag:arm.com,2025:psa#1.0.0"

// PSA certification number key for measurement values map extension
const PSACertNumKey = 100

// PSA software component key CBOR tag
const PSASoftwareComponentKeyTag = 800

// PSA software relations triple key
const PSASoftwareRelationsKey = 50

// PSACertNum represents a PSA Certified Security Assurance Certificate number
type PSACertNum struct {
	CertNumber string `cbor:"100,keyasint" json:"psa-cert-num"`
}

// Registering the PSA profile inside init() ensures that the profile will
// always be available and you don't need to remember to register it when
// you want to use it.
func init() {
	var err error
	ProfileID, err = eat.NewProfile(PSAProfileURI)
	if err != nil {
		panic(err) // will not error, as the profile URI is valid
	}

	// Create extensions map for PSA profile
	extMap := extensions.NewMap().
		Add(comid.ExtMval, &PSACertNum{}).
		Add(comid.ExtPSASwRelTriples, &PSASwRelTriples{})

	if err := corim.RegisterProfile(ProfileID, extMap); err != nil {
		// will not error, assuming our profile ID is unique and we've
		// correctly set up the extensions Map above
		panic(err)
	}
	
	// Register PSA measurement key types
	if err := comid.RegisterMkeyType(PSASoftwareComponentKeyTag, newMkeyPSASoftwareComponent); err != nil {
		panic(fmt.Sprintf("failed to register PSA software component key type: %v", err))
	}
}
