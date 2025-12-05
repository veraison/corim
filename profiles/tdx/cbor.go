// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"reflect"

	cbor "github.com/fxamacker/cbor/v2"
)

var (
	em, emError = initCBOREncMode()
	dm, dmError = initCBORDecMode()

	tdxTagsMap = map[uint64]interface{}{
		// TDX tags
		60010: TaggedNumericExpression{},
		60020: TaggedSetDigestExpression{},
		60021: TaggedSetStringExpression{},
	}
)

func tdxTags() cbor.TagSet {
	opts := cbor.TagOptions{
		EncTag: cbor.EncTagRequired,
		DecTag: cbor.DecTagRequired,
	}

	tags := cbor.NewTagSet()

	for tag, typ := range tdxTagsMap {
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
	return encOpt.EncModeWithTags(tdxTags())
}

func initCBORDecMode() (dm cbor.DecMode, err error) {
	decOpt := cbor.DecOptions{
		IndefLength: cbor.IndefLengthAllowed,
	}
	return decOpt.DecModeWithTags(tdxTags())
}

func init() {
	if emError != nil {
		panic(emError)
	}
	if dmError != nil {
		panic(dmError)
	}
}
