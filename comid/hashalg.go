package comid

import (
    "fmt"
    "strings"
    "encoding/json"
)

type HashAlg uint64

const (
    HashAlgSHA256 HashAlg = 1
    HashAlgSHA384 HashAlg = 2
    HashAlgSHA512 HashAlg = 3
)

func (h HashAlg) Valid() bool {
    return h >= HashAlgSHA256 && h <= HashAlgSHA512
}

func HashAlgFromString(s string) HashAlg {
    switch strings.ToLower(s) {
    case "sha-256":
        return HashAlgSHA256
    case "sha-384":
        return HashAlgSHA384
    case "sha-512":
        return HashAlgSHA512
    default:
        return 0
    }
}

func (h HashAlg) String() string {
    switch h {
    case HashAlgSHA256:
        return "sha-256"
    case HashAlgSHA384:
        return "sha-384"
    case HashAlgSHA512:
        return "sha-512"
    default:
        return fmt.Sprintf("unknown(%d)", h)
    }
}

func (h HashAlg) MarshalJSON() ([]byte, error) {
    return json.Marshal(h.String())
}
func (h *HashAlg) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    *h = HashAlgFromString(s)
    if !h.Valid() {
        return fmt.Errorf("invalid hash algorithm: %s", s)
    }
    return nil
}

// ToUint64 returns 0 if invalid, otherwise the numeric value.
func (h HashAlg) ToUint64() uint64 {
    if !h.Valid() {
        return 0
    }
    return uint64(h)
}

// HashAlgFromUint64 returns 0 if v is invalid, otherwise the matching HashAlg.
func HashAlgFromUint64(v uint64) HashAlg {
    h := HashAlg(v)
    if !h.Valid() {
        return 0
    }
    return h
}