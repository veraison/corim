package comid

import (
    "encoding/json"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func Test_HashAlgFromString(t *testing.T) {
    tests := []struct {
        name string
        input string
        want HashAlg
    }{
        {"sha-256", "sha-256", HashAlgSHA256},
        {"SHA-256", "SHA-256", HashAlgSHA256},
        {"sha-384", "sha-384", HashAlgSHA384},
        {"sha-512", "sha-512", HashAlgSHA512},
        {"invalid", "invalid", 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := HashAlgFromString(tt.input)
            assert.Equal(t, tt.want, got)
        })
    }
}

func Test_HashAlg_String(t *testing.T) {
    tests := []struct {
        name string
        hash HashAlg
        want string
    }{
        {"sha-256", HashAlgSHA256, "sha-256"},
        {"sha-384", HashAlgSHA384, "sha-384"},
        {"sha-512", HashAlgSHA512, "sha-512"},
        {"invalid", 99, "unknown(99)"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.hash.String()
            assert.Equal(t, tt.want, got)
        })
    }
}

func Test_HashAlg_JSON(t *testing.T) {
    tests := []struct {
        name string
        hash HashAlg
        want string
    }{
        {"sha-256", HashAlgSHA256, `"sha-256"`},
        {"sha-384", HashAlgSHA384, `"sha-384"`},
        {"sha-512", HashAlgSHA512, `"sha-512"`},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            data, err := json.Marshal(tt.hash)
            require.NoError(t, err)
            assert.Equal(t, tt.want, string(data))

            var got HashAlg
            err = json.Unmarshal(data, &got)
            require.NoError(t, err)
            assert.Equal(t, tt.hash, got)
        })
    }
}

// ...existing code...

func Test_HashAlg_Uint64(t *testing.T) {
    tests := []struct {
        name string
        hash HashAlg
        want uint64
    }{
        {"sha-256", HashAlgSHA256, 1},
        {"sha-384", HashAlgSHA384, 2},
        {"sha-512", HashAlgSHA512, 3},
        {"invalid", HashAlg(99), 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 1) Check forward conversion (HashAlg → uint64).
            got := tt.hash.ToUint64()
            assert.Equal(t, tt.want, got)

            // 2) Check backward conversion (uint64 → HashAlg).
            // For valid alg we expect round-trip equality, for invalid we do not.
            recon := HashAlgFromUint64(tt.want)
            if tt.name == "invalid" {
                // The want is 0, so we expect recon == 0, not 99.
                assert.Equal(t, HashAlg(0), recon)
            } else {
                // Valid case - recon should match original hash.
                assert.Equal(t, tt.hash, recon)
            }
        })
    }
}