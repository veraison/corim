// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"fmt"
)

type EntityExtension struct {
	Param1 string `cbor:"0,keyasint" json:"param1"`
	Param2 byte   `cbor:"1,keyasint,omitempty" json:"param2,omitempty"`
}

func (o EntityExtension) Valid() error {
	if o.Param1 == "" {
		return fmt.Errorf("mandatory Parameter Param1 missing")
	}
	if o.Param2 == 0x40 {
		return fmt.Errorf("mandatory Parameter Param2 cannot be 0x40")
	}

	return nil
}

func (o *EntityExtension) FromCBOR(data []byte) error {
	if o != nil {
		err := dm.Unmarshal(data, o)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *EntityExtension) ToCBOR() ([]byte, error) {
	if o != nil {
		data, err := em.Marshal(o)
		if err != nil {
			return nil, err
		} else {
			return data, nil
		}
	} else {
		return nil, fmt.Errorf("no entity to serialize")
	}
}
