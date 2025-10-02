// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"github.com/veraison/eat"
)

var (
	// PSA Token Profile ID using tag URI scheme
	TokenProfileID *eat.Profile

	// PSA Platform Endorsements Profile ID using tag URI scheme 
	EndorsementsProfileID *eat.Profile
)

func init() {
	var err error

	TokenProfileID, err = eat.NewProfile("tag:trustedcomputinggroup.org,2025:psa-token")
	if err != nil {
		panic(err)
	}

	EndorsementsProfileID, err = eat.NewProfile("tag:trustedcomputinggroup.org,2025:psa-endorsements") 
	if err != nil {
		panic(err)
	}
}