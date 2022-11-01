// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "fmt"

type CCAPlatformConfigID string

type TaggedCCAPlatformConfigID CCAPlatformConfigID

func (o CCAPlatformConfigID) Empty() bool {
	return o == ""
}

func (o *CCAPlatformConfigID) Set(v string) error {
	if v == "" {
		return fmt.Errorf("empty input string")
	}
	*o = CCAPlatformConfigID(v)
	return nil
}

func (o CCAPlatformConfigID) Get() (CCAPlatformConfigID, error) {
	if o == "" {
		return "", fmt.Errorf("empty CCA platform config ID")
	}
	return o, nil
}
