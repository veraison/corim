// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"time"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/corim"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
)

// MValExtensions contains the Intel TDX profile extensions which can appear in
// both Reference Values and Endorsed Values
type MValExtensions struct {
	TeeTcbDate     *time.Time        `cbor:"-72,keyasint,omitempty" json:"tcbdate,omitempty"`
	TeeISVSVN      *TeeSVN           `cbor:"-73,keyasint,omitempty" json:"isvsvn,omitempty"`
	TeeInstanceID  *TeeInstanceID    `cbor:"-77,keyasint,omitempty" json:"instanceid,omitempty"`
	TeePCEID       *TeePCEID         `cbor:"-80,keyasint,omitempty" json:"pceid,omitempty"`
	TeeMiscSelect  *TeeMiscSelect    `cbor:"-81,keyasint,omitempty" json:"miscselect,omitempty"`
	TeeAttributes  *TeeAttributes    `cbor:"-82,keyasint,omitempty" json:"attributes,omitempty"`
	TeeMrTee       *TeeDigest        `cbor:"-83,keyasint,omitempty" json:"mrtee,omitempty"`
	TeeMrSigner    *TeeDigest        `cbor:"-84,keyasint,omitempty" json:"mrsigner,omitempty"`
	TeeISVProdID   *TeeISVProdID     `cbor:"-85,keyasint,omitempty" json:"isvprodid,omitempty"`
	TeeTcbEvalNum  *TeeTcbEvalNumber `cbor:"-86,keyasint,omitempty" json:"tcbevalnum,omitempty"`
	TeeTcbStatus   *TeeTcbStatus     `cbor:"-88,keyasint,omitempty" json:"tcbstatus,omitempty"`
	TeeAdvisoryIDs *TeeAdvisoryIDs   `cbor:"-89,keyasint,omitempty" json:"advisoryids,omitempty"`
	TeeEpoch       *time.Time        `cbor:"-90, keyasint,omitempty" json:"epoch,omitempty"`

	TeeCryptoKeys *comid.CryptoKeys `cbor:"-91, keyasint,omitempty" json:"teecryptokeys,omitempty"`
	TeeTCBCompSvn *TeeTcbCompSvn    `cbor:"-125, keyasint,omitempty" json:"tcbcompsvn,omitempty"`
}

// Registering the profile inside init() in the same file where it is defined
// ensures that the profile will always be available, and you don't need to
// remember to register it when you want to use it. The only potential
// danger with that is if your profile ID clashes with another profile,
// which should not happen if it is a registered PEN or a URL containing a domain
// that you own.
// Note Intel profile is "2.16.840.1.113741.1.16.1",
// which is "joint-iso-itu-t.country.us.organization.intel.intel-comid.profile"

func init() {
	profileID, err := eat.NewProfile("2.16.840.1.113741.1.16.1")
	if err != nil {
		panic(err) // will not error, as the hard-coded string above is valid
	}

	extMap := extensions.NewMap().
		Add(comid.ExtReferenceValue, &MValExtensions{}).
		Add(comid.ExtEndorsedValue, &MValExtensions{})

	if err := corim.RegisterProfile(profileID, extMap); err != nil {
		// will not error, assuming our profile ID is unique, and we've
		// correctly set up the extensions Map above
		panic(err)
	}
}
