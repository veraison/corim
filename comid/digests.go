package comid

import (
    "fmt"
    "github.com/veraison/swid"
)

// Digests is an array of SWID HashEntry
type Digests []swid.HashEntry

// NewDigests instantiates an empty array of Digests
func NewDigests() *Digests {
    return new(Digests)
}

// AddDigest create a new digest from the supplied arguments and appends it to the (already instantiated) Digests target.
// The method is a no-op if it is invoked on a nil target and will refuse to add inconsistent algo/value combinations.
func (o *Digests) AddDigest(algID uint64, value []byte) *Digests {
    if o != nil {
        he := NewHashEntry(algID, value)
        if he == nil {
            return nil
        }
        *o = append(*o, *he)
    }
    return o
}

func (o Digests) Valid() error {
    if len(o) == 0 {
        return fmt.Errorf("digests must not be empty")
    }
    
    for i, m := range o {
        if err := swid.ValidHashEntry(m.HashAlgID, m.HashValue); err != nil {
            return fmt.Errorf("digest at index %d: %w", i, err)
        }
    }
    return nil
}


func NewHashEntry(algID uint64, value []byte) *swid.HashEntry {
    var he swid.HashEntry

    err := he.Set(algID, value)
    if err != nil {
        return nil
    }

    return &he
}
