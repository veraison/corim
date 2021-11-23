// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CorimExtractCmd_unknown_argument(t *testing.T) {
	cmd := NewCorimExtractCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_CorimExtractCmd_mandatory_args_missing_corim_file(t *testing.T) {
	cmd := NewCorimExtractCmd()

	args := []string{
		"--output-dir=ignore.d/",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoRIM supplied")
}

func Test_CorimExtractCmd_non_existent_corim_file(t *testing.T) {
	cmd := NewCorimExtractCmd()

	args := []string{
		"--file=nonexistent.cbor",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()

	err := cmd.Execute()
	assert.EqualError(t, err, "error loading signed CoRIM from nonexistent.cbor: open nonexistent.cbor: file does not exist")
}

func Test_CorimExtractCmd_bad_signed_corim(t *testing.T) {
	cmd := NewCorimExtractCmd()

	args := []string{
		"--file=bad.txt",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "bad.txt", []byte("hello!"), 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error decoding signed CoRIM from bad.txt: failed CBOR decoding for COSE-Sign1 signed CoRIM: unexpected EOF")
}

func Test_CorimExtractCmd_invalid_signed_corim(t *testing.T) {
	cmd := NewCorimExtractCmd()

	args := []string{
		"--file=invalid.cbor",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "invalid.cbor", testSignedCorimInvalid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error decoding signed CoRIM from invalid.cbor: failed validation of unsigned CoRIM: empty id")
}

func Test_CorimExtractCmd_ok_save_to_default_dir(t *testing.T) {
	cmd := NewCorimExtractCmd()

	args := []string{
		"--file=ok.cbor",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testSignedCorimValid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.NoError(t, err)

	_, err = fs.Stat("000000-comid.cbor")
	assert.NoError(t, err)

}

func Test_CorimExtractCmd_ok_save_to_non_default_dir(t *testing.T) {
	cmd := NewCorimExtractCmd()

	args := []string{
		"--file=ok.cbor",
		"--output-dir=my-dir/",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testSignedCorimValid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.NoError(t, err)

	_, err = fs.Stat("my-dir/000000-comid.cbor")
	assert.NoError(t, err)
}
