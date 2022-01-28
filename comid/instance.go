package comid

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/veraison/eat"
)

// Instance stores an instance identity. The supported formats are UUID and UEID.
type Instance struct {
	val interface{}
}

// NewInstance instantiates an empty instance
func NewInstance() *Instance {
	return &Instance{}
}

// SetUEID sets the identity of the target instance to the supplied UEID
func (o *Instance) SetUEID(val eat.UEID) *Instance {
	if o != nil {
		if val.Validate() != nil {
			return nil
		}
		o.val = TaggedUEID(val)
	}
	return o
}

// SetUUID sets the identity of the target instance to the supplied UUID
func (o *Instance) SetUUID(val uuid.UUID) *Instance {
	if o != nil {
		o.val = TaggedUUID(val)
	}
	return o
}

// NewInstanceUEID instantiates a new instance with the supplied UEID identity
func NewInstanceUEID(val eat.UEID) *Instance {
	return NewInstance().SetUEID(val)
}

// NewInstanceUUID instantiates a new instance with the supplied UUID identity
func NewInstanceUUID(val uuid.UUID) *Instance {
	return NewInstance().SetUUID(val)
}

// Valid checks for the validity of given instance
func (o Instance) Valid() error {
	if o.String() == "" {
		return fmt.Errorf("invalid instance id")
	}
	return nil
}

func (o Instance) GetUEID() (eat.UEID, error) {
	switch t := o.val.(type) {
	case TaggedUEID:
		return eat.UEID(t), nil
	default:
		return eat.UEID{}, fmt.Errorf("instance-id type is: %T", t)
	}
}

func (o Instance) GetUUID() (string, error) {
	switch t := o.val.(type) {
	case TaggedUUID:
		return UUID(t).String(), nil
	default:
		return "", fmt.Errorf("instance-id type is: %T", t)
	}
}

// String returns a printable string of the Instance value.  UUIDs use the
// canonical 8-4-4-4-12 format, UEIDs are hex encoded.
func (o Instance) String() string {
	switch t := o.val.(type) {
	case TaggedUUID:
		return UUID(t).String()
	case TaggedUEID:
		return hex.EncodeToString(t)
	default:
		return ""
	}
}

// MarshalCBOR serializes the target instance to CBOR
func (o Instance) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

func (o *Instance) UnmarshalCBOR(data []byte) error {
	var ueid TaggedUEID

	if dm.Unmarshal(data, &ueid) == nil {
		o.val = ueid
		return nil
	}

	var u TaggedUUID

	if dm.Unmarshal(data, &u) == nil {
		o.val = u
		return nil
	}

	return fmt.Errorf("unknown instance type (CBOR: %x)", data)
}

// UnmarshalJSON deserializes the supplied JSON type/value object into the Group
// target.  The supported formats are UUID, e.g.:
//
//   {
//     "type": "uuid",
//     "value": "69E027B2-7157-4758-BCB4-D9F167FE49EA"
//   }
//
// and UEID:
//
//   {
//     "type": "ueid",
//     "value": "Ad6tvu/erb7v3q2+796tvu8="
//   }
//
func (o *Instance) UnmarshalJSON(data []byte) error {
	var v tnv

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Type {
	case "uuid":
		var x UUID
		if err := x.UnmarshalJSON(v.Value); err != nil {
			return err
		}
		o.val = TaggedUUID(x)
	case "ueid":
		var x UEID
		if err := x.UnmarshalJSON(v.Value); err != nil {
			return err
		}
		o.val = TaggedUEID(x)
	default:
		return fmt.Errorf("unknown type %s for instance", v.Type)
	}

	return nil
}

func (o Instance) MarshalJSON() ([]byte, error) {
	var (
		v   tnv
		b   []byte
		err error
	)

	switch t := o.val.(type) {
	case TaggedUUID:
		b, err = UUID(t).MarshalJSON()
		if err != nil {
			return nil, err
		}
		v = tnv{Type: "uuid", Value: b}
	case TaggedUEID:
		b, err = UEID(t).MarshalJSON()
		if err != nil {
			return nil, err
		}
		v = tnv{Type: "ueid", Value: b}
	default:
		return nil, fmt.Errorf("unknown type %T for instance", t)
	}

	return json.Marshal(v)
}
