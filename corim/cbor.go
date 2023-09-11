// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"reflect"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/veraison/corim/v2/comid"
)

var (
	em, emError = initCBOREncMode()
	dm, dmError = initCBORDecMode()
)

var (
	CoswidTag = []byte{0xd9, 0x01, 0xf9} // 505()
	ComidTag  = []byte{0xd9, 0x01, 0xfa} // 506()
)

func corimTags() cbor.TagSet {
	corimTagsMap := map[uint64]interface{}{
		32: comid.TaggedURI(""),
	}

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
		IndefLength: cbor.IndefLengthForbidden,
		TimeTag:     cbor.DecTagRequired,
	}
	return decOpt.DecModeWithTags(corimTags())
}

func init() {
	if emError != nil {
		panic(emError)
	}
	if dmError != nil {
		panic(dmError)
	}
}
