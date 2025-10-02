// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cca

import (
	"github.com/veraison/eat"
)

var (
	// CCA Token Profile ID using tag URI scheme
	TokenProfileID *eat.Profile

	// CCA Platform Endorsements Profile ID using tag URI scheme
	EndorsementsProfileID *eat.Profile 

	// CCA Realm Endorsements Profile ID using tag URI scheme
	RealmEndorsementsProfileID *eat.Profile
)

func init() {
	var err error

	TokenProfileID, err = eat.NewProfile("tag:arm.com,2025:cca-token")
	if err != nil {
		panic(err)
	}

	EndorsementsProfileID, err = eat.NewProfile("tag:arm.com,2025:cca-endorsements")
	if err != nil {
		panic(err)
	}

	RealmEndorsementsProfileID, err = eat.NewProfile("tag:arm.com,2025:cca-realm-endorsements")
	if err != nil {
		panic(err)
	}
}