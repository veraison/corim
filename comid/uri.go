// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"errors"
	"fmt"
	"net/url"
)

const URIType = "uri"

type TaggedURI string

func NewTaggedURI(val any) (*TaggedURI, error) {
	if val == nil {
		ret := TaggedURI("")
		return &ret, nil
	}

	switch t := val.(type) {
	case string:
		ret := TaggedURI(t)
		if err := ret.Valid(); err != nil {
			return nil, err
		}
		return &ret, nil
	case *string:
		ret := TaggedURI(*t)
		if err := ret.Valid(); err != nil {
			return nil, err
		}
		return &ret, nil
	case url.URL:
		ret := TaggedURI(t.String())
		return &ret, nil
	case *url.URL:
		ret := TaggedURI(t.String())
		return &ret, nil
	default:
		return nil, fmt.Errorf("unexpected input type for URI: %v(%T)", t, t)
	}
}

func (o TaggedURI) Empty() bool {
	return o == ""
}

func (o TaggedURI) Type() string {
	return "uri"
}

func (o TaggedURI) String() string {
	return string(o)
}

func (o TaggedURI) Valid() error {
	if o.Empty() {
		return errors.New("empty URI")
	}

	if _, err := url.Parse(string(o)); err != nil {
		return err
	}

	return nil
}
