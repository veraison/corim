package corim

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/swid"
)

func TestLocator_JSON_thumbprint_format_backward_compatibility(t *testing.T) {
	// Test hash entry from the GitHub issue - corrected hash value
	he := swid.HashEntry{
		HashAlgID: 1, // SHA-256
		HashValue: []byte{0xe4, 0x5b, 0x72, 0xf5, 0xc0, 0xc0, 0xb5, 0x72, 0xdb, 0x4d, 0x8d, 0x3a, 0xb7, 0xe9, 0x7f, 0x36, 0x8f, 0xf7, 0x4e, 0x62, 0x34, 0x7a, 0x82, 0x4d, 0xec, 0xb6, 0x7a, 0x84, 0xe5, 0x22, 0x4d, 0x75},
	}

	t.Run("marshal uses colon format for backward compatibility", func(t *testing.T) {
		locator := Locator{
			Href:       "https://parent.example/rims/ccb3aa85-61b4-40f1-848e-02ad6e8a254b",
			Thumbprint: &he,
		}

		jsonBytes, err := json.Marshal(locator)
		require.NoError(t, err)

		jsonStr := string(jsonBytes)

		// Should contain colon format, not semicolon
		assert.Contains(t, jsonStr, "sha-256:")
		assert.NotContains(t, jsonStr, "sha-256;")

		// Specific expected output from the issue
		assert.Contains(t, jsonStr, "sha-256:5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU=")
	})

	t.Run("unmarshal supports old format (colon)", func(t *testing.T) {
		// JSON with old format (colon)
		jsonInput := `{
			"href": "https://parent.example/rims/ccb3aa85-61b4-40f1-848e-02ad6e8a254b",
			"thumbprint": "sha-256:5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
		}`

		var locator Locator
		err := json.Unmarshal([]byte(jsonInput), &locator)
		require.NoError(t, err)

		assert.Equal(t, "https://parent.example/rims/ccb3aa85-61b4-40f1-848e-02ad6e8a254b", string(locator.Href))
		assert.NotNil(t, locator.Thumbprint)
		assert.Equal(t, uint64(1), locator.Thumbprint.HashAlgID) // SHA-256
		assert.Equal(t, he.HashValue, locator.Thumbprint.HashValue)
	})

	t.Run("unmarshal supports new format (semicolon)", func(t *testing.T) {
		// JSON with new format (semicolon)
		jsonInput := `{
			"href": "https://parent.example/rims/ccb3aa85-61b4-40f1-848e-02ad6e8a254b",
			"thumbprint": "sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
		}`

		var locator Locator
		err := json.Unmarshal([]byte(jsonInput), &locator)
		require.NoError(t, err)

		assert.Equal(t, "https://parent.example/rims/ccb3aa85-61b4-40f1-848e-02ad6e8a254b", string(locator.Href))
		assert.NotNil(t, locator.Thumbprint)
		assert.Equal(t, uint64(1), locator.Thumbprint.HashAlgID) // SHA-256
		assert.Equal(t, he.HashValue, locator.Thumbprint.HashValue)
	})

	t.Run("round-trip test maintains colon format", func(t *testing.T) {
		// Start with old format JSON
		jsonInput := `{
			"href": "https://parent.example/rims/ccb3aa85-61b4-40f1-848e-02ad6e8a254b",
			"thumbprint": "sha-256:5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
		}`

		var locator Locator
		err := json.Unmarshal([]byte(jsonInput), &locator)
		require.NoError(t, err)

		// Marshal back to JSON
		jsonBytes, err := json.Marshal(locator)
		require.NoError(t, err)

		jsonStr := string(jsonBytes)

		// Should still use colon format
		assert.Contains(t, jsonStr, "sha-256:")
		assert.NotContains(t, jsonStr, "sha-256;")
	})

	t.Run("UnsignedCorim integration test", func(t *testing.T) {
		ucorim := NewUnsignedCorim().
			SetID("5c57e8f4-46cd-421b-91c9-08cf93e13cfc").
			AddDependentRim("https://parent.example/rims/ccb3aa85-61b4-40f1-848e-02ad6e8a254b", &he)

		jsonBytes, err := ucorim.ToJSON()
		require.NoError(t, err)

		jsonStr := string(jsonBytes)

		// Should contain colon format for backward compatibility
		assert.Contains(t, jsonStr, "sha-256:")
		assert.NotContains(t, jsonStr, "sha-256;")

		// Should match the expected format from the GitHub issue
		expectedThumbprint := "sha-256:5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
		assert.Contains(t, jsonStr, expectedThumbprint)
	})

	t.Run("no thumbprint case", func(t *testing.T) {
		locator := Locator{
			Href: "https://example.com/rim",
		}

		jsonBytes, err := json.Marshal(locator)
		require.NoError(t, err)

		var result Locator
		err = json.Unmarshal(jsonBytes, &result)
		require.NoError(t, err)

		assert.Equal(t, "https://example.com/rim", string(result.Href))
		assert.Nil(t, result.Thumbprint)
	})
}
