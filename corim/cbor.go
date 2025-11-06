// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"fmt"
	"reflect"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/veraison/corim/comid"
)

var (
	em, emError = initCBOREncMode()
	dm, dmError = initCBORDecMode()
)

var (
	UnsignedCorimTag        = []byte{0xd9, 0x01, 0xf5} // 501()
	CoswidTag        uint64 = 505
	ComidTag         uint64 = 506

	corimTagsMap = map[uint64]interface{}{
		32: comid.TaggedURI(""),
	}
)

func corimTags() cbor.TagSet {
	opts := cbor.TagOptions{
		EncTag: cbor.EncTagRequired,
		DecTag: cbor.DecTagRequired,
	}

	tags := cbor.NewTagSet()

	for tag, typ := range corimTagsMap {
		if err := tags.Add(opts, reflect.TypeOf(typ), tag); err != nil {
			panic(err)
		}
	}

	return tags
}

func initCBOREncMode() (en cbor.EncMode, err error) {
	encOpt := cbor.EncOptions{
		IndefLength: cbor.IndefLengthForbidden,
		TimeTag:     cbor.EncTagRequired,
	}
	return encOpt.EncModeWithTags(corimTags())
}

func initCBORDecMode() (dm cbor.DecMode, err error) {
	decOpt := cbor.DecOptions{
		IndefLength: cbor.IndefLengthAllowed,
		TimeTag:     cbor.DecTagRequired,
	}
	return decOpt.DecModeWithTags(corimTags())
}

func registerCORIMTag(tag uint64, t interface{}) error {
	if _, exists := corimTagsMap[tag]; exists {
		return fmt.Errorf("tag %d is already registered", tag)
	}

	corimTagsMap[tag] = t

	var err error

	em, err = initCBOREncMode()
	if err != nil {
		return err
	}

	dm, err = initCBORDecMode()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	if emError != nil {
		panic(emError)
	}
	if dmError != nil {
		panic(dmError)
	}
}
