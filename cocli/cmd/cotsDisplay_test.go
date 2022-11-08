// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CotsDisplayCmd_unknown_argument(t *testing.T) {
	cmd := NewCotsDisplayCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_CotsDisplayCmd_mandatory_args_missing_cots_file(t *testing.T) {
	cmd := NewCotsDisplayCmd()

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoTS supplied")
}

func Test_CotsDisplayCmd_non_existent_cots_file(t *testing.T) {
	cmd := NewCotsDisplayCmd()

	args := []string{
		"--file=nonexistent.cbor",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()

	err := cmd.Execute()
	assert.EqualError(t, err, "error loading signed CoTS from nonexistent.cbor: open nonexistent.cbor: file does not exist")
}

func Test_CotsDisplayCmd_bad_unsigned_corim(t *testing.T) {
	cmd := NewCotsDisplayCmd()

	args := []string{
		"--file=bad.txt",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "bad.txt", []byte("hello!"), 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "error decoding signed CoTS from bad.txt: failed CBOR decoding for COSE-Sign1 signed CoRIM: cbor: invalid COSE_Sign1_Tagged object")
}

func Test_CotsDisplayCmd_invalid_unsigned_cots(t *testing.T) {
	cmd := NewCotsDisplayCmd()

	args := []string{
		"--file=../data/cots/not.cbor",
	}
	cmd.SetArgs(args)

	fs = afero.NewOsFs()

	err := cmd.Execute()
	assert.EqualError(t, err, "error decoding signed CoTS from ../data/cots/not.cbor: failed CBOR decoding for COSE-Sign1 signed CoRIM: cbor: invalid COSE_Sign1_Tagged object")
}

func Test_CotsDisplayCmd_ok_top_level_view(t *testing.T) {
	cmd := NewCotsDisplayCmd()

	args := []string{
		"--file=../data/cots/signed.cbor",
	}
	cmd.SetArgs(args)

	fs = afero.NewOsFs()

	err := cmd.Execute()
	assert.NoError(t, err)
}
