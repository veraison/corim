// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CorimDisplayCmd_unknown_argument(t *testing.T) {
	cmd := NewCorimDisplayCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_CorimDisplayCmd_mandatory_args_missing_corim_file(t *testing.T) {
	cmd := NewCorimDisplayCmd()

	args := []string{
		"--show-tags",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoRIM supplied")
}

func Test_CorimDisplayCmd_non_existent_corim_file(t *testing.T) {
	cmd := NewCorimDisplayCmd()

	args := []string{
		"--file=nonexistent.cbor",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()

	err := cmd.Execute()
	assert.EqualError(t, err, "error loading signed CoRIM from nonexistent.cbor: open nonexistent.cbor: file does not exist")
}

func Test_CorimDisplayCmd_bad_signed_corim(t *testing.T) {
	cmd := NewCorimDisplayCmd()

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

func Test_CorimDisplayCmd_invalid_signed_corim(t *testing.T) {
	cmd := NewCorimDisplayCmd()

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

func Test_CorimDisplayCmd_ok_top_level_view(t *testing.T) {
	cmd := NewCorimDisplayCmd()

	args := []string{
		"--file=ok.cbor",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testSignedCorimValid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.NoError(t, err)
}

func Test_CorimDisplayCmd_ok_nested_view(t *testing.T) {
	cmd := NewCorimDisplayCmd()

	args := []string{
		"--file=ok.cbor",
		"--show-tags",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "ok.cbor", testSignedCorimValid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.NoError(t, err)
}
