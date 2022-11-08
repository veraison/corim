// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CorimVerifyCmd_unknown_argument(t *testing.T) {
	cmd := NewCorimVerifyCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_CorimVerifyCmd_mandatory_args_missing_corim_file(t *testing.T) {
	cmd := NewCorimVerifyCmd()

	args := []string{
		"--key=ignored.jwk",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoRIM supplied")
}

func Test_CorimVerifyCmd_mandatory_args_missing_key_file(t *testing.T) {
	cmd := NewCorimVerifyCmd()

	args := []string{
		"--file=ignored.jwk",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no key supplied")
}

func Test_CorimVerifyCmd_non_existent_signed_corim_file(t *testing.T) {
	cmd := NewCorimVerifyCmd()

	args := []string{
		"--file=nonexistent.cbor",
		"--key=ignored.jwk",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()

	err := cmd.Execute()
	assert.EqualError(t, err, "error loading signed CoRIM from nonexistent.cbor: open nonexistent.cbor: file does not exist")
}

func Test_CorimVerifyCmd_bad_signed_corim(t *testing.T) {
	cmd := NewCorimVerifyCmd()

	args := []string{
		"--file=bad.txt",
		"--key=ignored.jwk",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "bad.txt", []byte("hello!"), 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error decoding signed CoRIM from bad.txt: failed CBOR decoding for COSE-Sign1 signed CoRIM: cbor: invalid COSE_Sign1_Tagged object")
}

func Test_CorimVerifyCmd_non_existent_key_file(t *testing.T) {
	cmd := NewCorimVerifyCmd()

	args := []string{
		"--file=ok.cbor",
		"--key=nonexistent.jwk",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testSignedCorimValid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error loading verifying key from nonexistent.jwk: open nonexistent.jwk: file does not exist")
}

func Test_CorimVerifyCmd_invalid_key_file(t *testing.T) {
	cmd := NewCorimVerifyCmd()

	args := []string{
		"--file=ok.cbor",
		"--key=invalid.jwk",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testSignedCorimValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "invalid.jwk", []byte("{}"), 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error loading verifying key from invalid.jwk: failed to unmarshal JWK set: failed to parse sole key in key set: invalid key type from JSON ()")
}

func Test_CorimVerifyCmd_ok(t *testing.T) {
	cmd := NewCorimVerifyCmd()

	args := []string{
		"--file=ok.cbor",
		"--key=ok.jwk",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testSignedCorimValid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ok.jwk", testECKey, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.NoError(t, err)
}
