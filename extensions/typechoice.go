// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package extensions

import (
	"encoding/json"

	"github.com/veraison/corim/encoding"
)

var StringType = "string"

// ITypeChoiceValue is the interface that is implemented by all concrete type
// choice value types. Specific type choices define their own value interfaces
// that embed this one (and possibly include additional methods).
type ITypeChoiceValue interface {
	// String returns the string representation of the ITypeChoiceValue.
	String() string
	// Valid returns an error if validation of the ITypeChoiceValue fails,
	// or nil if it succeeds.
	Valid() error
	// Type returns the type name of this ITypeChoiceValue implementation.
	Type() string
}

func TypeChoiceValueMarshalJSON(v ITypeChoiceValue) ([]byte, error) {
	valueBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	value := encoding.TypeAndValue{
		Type:  v.Type(),
		Value: valueBytes,
	}

	return json.Marshal(value)
}
