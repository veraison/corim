// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CorimSignCmd_unknown_argument(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_CorimSignCmd_mandatory_args_missing_corim_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--key=ignored.jwk",
		"--meta=ignored.json",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoRIM supplied")
}

func Test_CorimSignCmd_mandatory_args_missing_meta_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--file=ignored.cbor",
		"--key=ignored.jwk",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoRIM Meta supplied")
}

func Test_CorimSignCmd_mandatory_args_missing_key_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--file=ignored.cbor",
		"--meta=ignored.json",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no key supplied")
}

func Test_CorimSignCmd_non_existent_unsigned_corim_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--file=nonexistent.cbor",
		"--key=ignored.jwk",
		"--meta=ignored.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()

	err := cmd.Execute()
	assert.EqualError(t, err, "error loading unsigned CoRIM from nonexistent.cbor: open nonexistent.cbor: file does not exist")
}

func Test_CorimSignCmd_bad_unsigned_corim(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--file=bad.txt",
		"--key=ignored.jwk",
		"--meta=ignored.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "bad.txt", []byte("hello!"), 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error decoding unsigned CoRIM from bad.txt: unexpected EOF")
}

func Test_CorimSignCmd_invalid_unsigned_corim(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--file=invalid.cbor",
		"--key=ignored.jwk",
		"--meta=ignored.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "invalid.cbor", testCorimInvalid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error validating CoRIM: tags validation failed: no tags")
}

func Test_CorimSignCmd_non_existent_meta_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--file=ok.cbor",
		"--key=ignored.jwk",
		"--meta=nonexistent.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testCorimValid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error loading CoRIM Meta from nonexistent.json: open nonexistent.json: file does not exist")
}

func Test_CorimSignCmd_bad_meta_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--file=ok.cbor",
		"--key=ignored.jwk",
		"--meta=bad.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testCorimValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "bad.json", []byte("{"), 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error decoding CoRIM Meta from bad.json: unexpected end of JSON input")
}

func Test_CorimSignCmd_invalid_meta_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--file=ok.cbor",
		"--key=ignored.jwk",
		"--meta=invalid.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testCorimValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "invalid.json", testMetaInvalid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error validating CoRIM Meta: invalid meta: empty name")
}

func Test_CorimSignCmd_non_existent_key_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--file=ok.cbor",
		"--key=nonexistent.jwk",
		"--meta=ok.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testCorimValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ok.json", testMetaValid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error loading signing key from nonexistent.jwk: open nonexistent.jwk: file does not exist")
}

func Test_CorimSignCmd_invalid_key_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--file=ok.cbor",
		"--key=invalid.jwk",
		"--meta=ok.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testCorimValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ok.json", testMetaValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "invalid.jwk", []byte("{}"), 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error loading signing key from invalid.jwk: failed to unmarshal JWK set: failed to parse sole key in key set: invalid key type from JSON ()")
}

func Test_CorimSignCmd_ok_with_default_output_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--file=ok.cbor",
		"--key=ok.jwk",
		"--meta=ok.json",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testCorimValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ok.json", testMetaValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ok.jwk", testECKey, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.NoError(t, err)

	_, err = fs.Stat("signed-ok.cbor")
	assert.NoError(t, err)
}

func Test_CorimSignCmd_ok_with_custom_output_file(t *testing.T) {
	cmd := NewCorimSignCmd()

	args := []string{
		"--file=ok.cbor",
		"--key=ok.jwk",
		"--meta=ok.json",
		"--output=my-signed-corim.cbor",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testCorimValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ok.json", testMetaValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ok.jwk", testECKey, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.NoError(t, err)

	_, err = fs.Stat("my-signed-corim.cbor")
	assert.NoError(t, err)
}
