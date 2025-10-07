// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"

	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
)

// MemberVal holds membership-related values for a specific membership record.
// It contains various types of membership information that can be associated
// with an environment.
type MemberVal struct {
	GroupID        *string     `cbor:"0,keyasint,omitempty" json:"group-id,omitempty"`
	GroupName      *string     `cbor:"1,keyasint,omitempty" json:"group-name,omitempty"`
	Role           *string     `cbor:"2,keyasint,omitempty" json:"role,omitempty"`
	Status         *string     `cbor:"3,keyasint,omitempty" json:"status,omitempty"`
	Permissions    *[]string   `cbor:"4,keyasint,omitempty" json:"permissions,omitempty"`
	OrganizationID *string     `cbor:"5,keyasint,omitempty" json:"organization-id,omitempty"`
	UEID           *eat.UEID   `cbor:"6,keyasint,omitempty" json:"ueid,omitempty"`
	UUID           *UUID       `cbor:"7,keyasint,omitempty" json:"uuid,omitempty"`
	Name           *string     `cbor:"8,keyasint,omitempty" json:"name,omitempty"`
	Extensions
}

// RegisterExtensions registers a struct as a collections of extensions
func (o *MemberVal) RegisterExtensions(exts extensions.Map) error {
	for p, v := range exts {
		switch p {
		case ExtMemberVal:
			o.Register(v)
		default:
			return fmt.Errorf("%w: %q", extensions.ErrUnexpectedPoint, p)
		}
	}

	return nil
}

// GetExtensions returns previously registered extension
func (o *MemberVal) GetExtensions() extensions.IMapValue {
	return o.IMapValue
}

// UnmarshalCBOR deserializes from CBOR
func (o *MemberVal) UnmarshalCBOR(data []byte) error {
	return encoding.PopulateStructFromCBOR(dm, data, o)
}

// MarshalCBOR serializes to CBOR
func (o MemberVal) MarshalCBOR() ([]byte, error) {
	return encoding.SerializeStructToCBOR(em, o)
}

// UnmarshalJSON deserializes from JSON
func (o *MemberVal) UnmarshalJSON(data []byte) error {
	return encoding.PopulateStructFromJSON(data, o)
}

// MarshalJSON serializes to JSON
func (o MemberVal) MarshalJSON() ([]byte, error) {
	return encoding.SerializeStructToJSON(o)
}

// SetGroupID sets the group identifier for the membership.
func (o *MemberVal) SetGroupID(groupID string) *MemberVal {
	if o != nil {
		o.GroupID = &groupID
	}
	return o
}

// SetGroupName sets the group name for the membership.
func (o *MemberVal) SetGroupName(groupName string) *MemberVal {
	if o != nil {
		o.GroupName = &groupName
	}
	return o
}

// SetRole sets the role for the membership.
func (o *MemberVal) SetRole(role string) *MemberVal {
	if o != nil {
		o.Role = &role
	}
	return o
}

// SetStatus sets the status for the membership.
func (o *MemberVal) SetStatus(status string) *MemberVal {
	if o != nil {
		o.Status = &status
	}
	return o
}

// SetPermissions sets the permissions for the membership.
func (o *MemberVal) SetPermissions(permissions []string) *MemberVal {
	if o != nil {
		o.Permissions = &permissions
	}
	return o
}

// SetOrganizationID sets the organization identifier for the membership.
func (o *MemberVal) SetOrganizationID(orgID string) *MemberVal {
	if o != nil {
		o.OrganizationID = &orgID
	}
	return o
}

// SetUEID sets the UEID for the membership.
func (o *MemberVal) SetUEID(ueid eat.UEID) *MemberVal {
	if o != nil {
		o.UEID = &ueid
	}
	return o
}

// SetUUID sets the UUID for the membership.
func (o *MemberVal) SetUUID(uuid UUID) *MemberVal {
	if o != nil {
		o.UUID = &uuid
	}
	return o
}

// SetName sets the name for the membership.
func (o *MemberVal) SetName(name string) *MemberVal {
	if o != nil {
		o.Name = &name
	}
	return o
}

// Valid returns an error if none of the membership values are set and the Extensions are empty.
func (o MemberVal) Valid() error {
	// Check if no membership values are set
	if o.GroupID == nil &&
		o.GroupName == nil &&
		o.Role == nil &&
		o.Status == nil &&
		o.Permissions == nil &&
		o.OrganizationID == nil &&
		o.UEID == nil &&
		o.UUID == nil &&
		o.Name == nil &&
		o.IsEmpty() {

		return fmt.Errorf("no membership value set")
	}

	// Validate UEID if set
	if o.UEID != nil {
		if err := UEID(*o.UEID).Valid(); err != nil {
			return fmt.Errorf("UEID validation failed: %w", err)
		}
	}

	// Validate UUID if set
	if o.UUID != nil {
		if err := o.UUID.Valid(); err != nil {
			return fmt.Errorf("UUID validation failed: %w", err)
		}
	}

	return nil
}