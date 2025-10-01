// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"time"

	"github.com/google/uuid"
	"github.com/veraison/corim/comid"
)

//nolint:lll
var (
	TestUUIDString  = "31fb5abf-023e-4992-aa4e-95f9c1503bfa"
	TestUUID        = comid.UUID(uuid.Must(uuid.Parse(TestUUIDString)))
	TestUUID2String = "550e8400-e29b-41d4-a716-446655440000"
	TestUUID2       = comid.UUID(uuid.Must(uuid.Parse(TestUUID2String)))
	TestUUID3String = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	TestUUID3       = comid.UUID(uuid.Must(uuid.Parse(TestUUID3String)))
	TestProfile     = "https://abc.com"
	TestTag         = "00010001-0001-0001-0001-000100010001"
	TestDeviceID    = "BAD809B1-7032-43D9-8F94-BF128E5D061D"
	TestKey         = true
	TestDate, _     = time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
)
