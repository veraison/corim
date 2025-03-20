// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/veraison/corim/comid"
)

type EnvironmentSelector struct {
	Classes   *[]comid.Class    `cbor:"0,keyasint,omitempty"`
	Instances *[]comid.Instance `cbor:"1,keyasint,omitempty"`
	Groups    *[]comid.Group    `cbor:"2,keyasint,omitempty"`
}

func NewEnvironmentSelector() *EnvironmentSelector {
	return &EnvironmentSelector{}
}

func (o *EnvironmentSelector) AddClass(v comid.Class) *EnvironmentSelector {
	if o.Classes == nil {
		o.Classes = new([]comid.Class)
	}

	*o.Classes = append(*o.Classes, v)

	return o
}

func (o *EnvironmentSelector) AddInstance(v comid.Instance) *EnvironmentSelector {
	if o.Instances == nil {
		o.Instances = new([]comid.Instance)
	}

	*o.Instances = append(*o.Instances, v)

	return o
}

func (o *EnvironmentSelector) AddGroup(v comid.Group) *EnvironmentSelector {
	if o.Groups == nil {
		o.Groups = new([]comid.Group)
	}

	*o.Groups = append(*o.Groups, v)

	return o
}

func (o EnvironmentSelector) Valid() error {
	if o.Classes == nil && o.Groups == nil && o.Instances == nil {
		return errors.New("non-empty<> constraint violation")
	}

	if o.Classes != nil {
		if o.Instances != nil || o.Groups != nil {
			return errors.New("only one selector type is allowed")
		}
	} else if o.Instances != nil {
		if o.Groups != nil {
			return errors.New("only one selector type is allowed")
		}
	}

	return nil
}

func sqlClassSelector(v comid.Class) (string, error) {
	var classQuery []string

	if v.ClassID != nil {
		s, err := b64tnv(*v.ClassID, "class-id")
		if err != nil {
			return "", fmt.Errorf("error decoding class-id: %w", err)
		}
		classQuery = append(classQuery, fmt.Sprintf(`class-id = %s`, s))
	}

	if v.Vendor != nil {
		classQuery = append(classQuery, fmt.Sprintf(`class-vendor = %q`, *v.Vendor))
	}

	if v.Model != nil {
		classQuery = append(classQuery, fmt.Sprintf(`class-model = %q`, *v.Model))
	}

	// TODO layer and index

	return strings.Join(classQuery, " AND "), nil
}

func sqlInstanceSelector(v comid.Instance) (string, error) {
	s, err := b64tnv(v, "instance-id")
	if err != nil {
		return "", fmt.Errorf("error decoding instance-id: %w", err)
	}
	return fmt.Sprintf(`instance-id = %s`, s), nil
}

func sqlGroupSelector(v comid.Group) (string, error) {
	s, err := b64tnv(v, "group-id")
	if err != nil {
		return "", fmt.Errorf("error decoding group-id: %w", err)
	}
	return fmt.Sprintf(`group-id = %s`, s), nil
}

func (o EnvironmentSelector) ToSQL() (string, error) {
	var conditions []string

	if o.Classes != nil {
		for i, v := range *o.Classes {
			csel, err := sqlClassSelector(v)
			if err != nil {
				return "", fmt.Errorf("creating selector for class[%d]: %w", i, err)
			}

			conditions = append(conditions, "( "+csel+" )")
		}
	} else if o.Instances != nil {
		for i, v := range *o.Instances {
			isel, err := sqlInstanceSelector(v)
			if err != nil {
				return "", fmt.Errorf("creating selector for instance[%d]: %w", i, err)
			}

			conditions = append(conditions, "( "+isel+" )")
		}
	} else if o.Groups != nil {
		for i, v := range *o.Groups {
			gsel, err := sqlGroupSelector(v)
			if err != nil {
				return "", fmt.Errorf("creating selector for group[%d]: %w", i, err)
			}

			conditions = append(conditions, "( "+gsel+" )")
		}
	}

	return strings.Join(conditions, " OR "), nil
}

func b64tnv[T json.Marshaler](t T, id string) (string, error) {
	type tnv struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	}

	d, err := t.MarshalJSON()
	if err != nil {
		return "", fmt.Errorf("error encoding %s: %w", id, err)
	}

	e := &tnv{}

	err = json.Unmarshal(d, e)
	if err != nil {
		return "", fmt.Errorf("error decoding %s: %w", id, err)
	}

	return string(e.Value), nil
}
