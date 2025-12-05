// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"github.com/veraison/corim/coev"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/corim/profiles/tdx"
	"github.com/veraison/eat"
)

var ProfileID *eat.Profile

// Registering the profile inside init() in the same file where it is defined
// ensures that the profile will always be available, and you don't need to
// remember to register it when you want to use it. The only potential
// danger with that is if your profile ID clashes with another profile,
// which should not happen if it is a registered PEN or a URL containing a domain
// that you own.
// Note Intel profile is "2.16.840.1.113741.1.16.1",
// which is "joint-iso-itu-t.country.us.organization.intel.intel-comid.profile"

func init() {
	var err error
	ProfileID, err = eat.NewProfile("2.16.840.1.113741.1.16.1")
	if err != nil {
		panic(err) // will not error, as the hard-coded string above is valid
	}

	extMap := extensions.NewMap().
		Add(coev.ExtEvidenceTriples, &tdx.MValExtensions{})

	if err := coev.RegisterProfile(ProfileID, extMap); err != nil {
		// will not error, assuming our profile ID is unique, and we've
		// correctly set up the extensions Map above
		panic(err)
	}
}
