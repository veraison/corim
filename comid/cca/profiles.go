// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cca

import (
	"fmt"
	"regexp"

	"github.com/veraison/eat"
)

// tagURIPattern validates RFC 4151 tag URI format
// tag:authority,date:specific
var tagURIPattern = regexp.MustCompile(`^tag:[a-zA-Z0-9\.\-]+,\d{4}(-(0[1-9]|1[0-2])(-(0[1-9]|[12][0-9]|3[01]))?)?:.+$`)

// validateTagURI checks if the given string is a valid tag URI according to RFC 4151
func validateTagURI(uri string) error {
	if !tagURIPattern.MatchString(uri) {
		return fmt.Errorf("invalid tag URI format: %q (expected format: tag:authority,date:specific)", uri)
	}
	return nil
}

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

	// Validate and create Token Profile
	if err = validateTagURI("tag:arm.com,2025:cca-token"); err != nil {
		panic(err)
	}
	TokenProfileID, err = eat.NewProfile("tag:arm.com,2025:cca-token")
	if err != nil {
		panic(err)
	}

	// Validate and create Endorsements Profile
	if err = validateTagURI("tag:arm.com,2025:cca-endorsements"); err != nil {
		panic(err)
	}
	EndorsementsProfileID, err = eat.NewProfile("tag:arm.com,2025:cca-endorsements")
	if err != nil {
		panic(err)
	}

	// Validate and create Realm Endorsements Profile
	if err = validateTagURI("tag:arm.com,2025:cca-realm-endorsements"); err != nil {
		panic(err)
	}
	RealmEndorsementsProfileID, err = eat.NewProfile("tag:arm.com,2025:cca-realm-endorsements")
	if err != nil {
		panic(err)
	}
}
