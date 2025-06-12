// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"github.com/google/uuid"
	"github.com/veraison/corim/comid"
)

//nolint:lll
var (
	TestUUIDString = "31fb5abf-023e-4992-aa4e-95f9c1503bfa"
	TestUUID       = comid.UUID(uuid.Must(uuid.Parse(TestUUIDString)))
)
