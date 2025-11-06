// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"reflect"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/veraison/corim/comid"
)

var (
	em, emError        = initCBOREncMode()
	dm, dmError        = initCBORDecMode()
	ConciseEvidenceTag = []byte{0xd9, 0x02, 0x3B}
	coevTagMap         = map[uint64]interface{}{
		37: comid.TaggedUUID{},
	}
)

func coevTags() cbor.TagSet {
	opts := cbor.TagOptions{
		EncTag: cbor.EncTagRequired,
		DecTag: cbor.DecTagRequired,
	}

	tags := cbor.NewTagSet()

	for tag, typ := range coevTagMap {
		if err := tags.Add(opts, reflect.TypeOf(typ), tag); err != nil {
			panic(err)
		}
	}

	return tags
}

func initCBOREncMode() (en cbor.EncMode, err error) {
	encOpt := cbor.EncOptions{
		Sort:        cbor.SortCoreDeterministic,
		IndefLength: cbor.IndefLengthForbidden,
		TimeTag:     cbor.EncTagRequired,
	}
	return encOpt.EncModeWithTags(coevTags())
}

func initCBORDecMode() (dm cbor.DecMode, err error) {
	decOpt := cbor.DecOptions{
		IndefLength: cbor.IndefLengthAllowed,
	}
	return decOpt.DecModeWithTags(coevTags())
}

func init() {
	if emError != nil {
		panic(emError)
	}
	if dmError != nil {
		panic(dmError)
	}
}
