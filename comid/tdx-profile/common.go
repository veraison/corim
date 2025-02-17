// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

func isType[T any](v any) bool {
	_, ok := v.(T)
	return ok
}
