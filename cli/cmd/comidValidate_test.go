// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ComidValidateCmd_unknown_argument(t *testing.T) {
	cmd := NewComidValidateCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_ComidValidateCmd_no_files(t *testing.T) {
	cmd := NewComidValidateCmd()

	// no args

	err := cmd.Execute()
	assert.EqualError(t, err, "no files supplied")
}

func Test_ComidValidateCmd_no_files_found(t *testing.T) {
	cmd := NewComidValidateCmd()

	args := []string{
		"--file=unknown",
		"--dir=unsure",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no files found")
}

func Test_ComidValidateCmd_file_with_invalid_cbor(t *testing.T) {
	var err error

	cmd := NewComidValidateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "invalid.cbor", []byte{0xff, 0xff}, 0400)
	require.NoError(t, err)

	args := []string{
		"--file=invalid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, "1/1 validation(s) failed")
}

func Test_ComidValidateCmd_file_with_invalid_comid(t *testing.T) {
	var err error

	cmd := NewComidValidateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "bad-comid.cbor", []byte{0xa0}, 0400)
	require.NoError(t, err)

	args := []string{
		"--file=bad-comid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, "1/1 validation(s) failed")
}

func Test_ComidValidateCmd_file_with_valid_comid(t *testing.T) {
	var err error

	cmd := NewComidValidateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "ok.cbor", PSARefValCBOR, 0400)
	require.NoError(t, err)

	args := []string{
		"--file=ok.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.NoError(t, err)
}

func Test_ComidValidateCmd_file_with_valid_comid_from_dir(t *testing.T) {
	var err error

	cmd := NewComidValidateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "testdir/ok.cbor", PSARefValCBOR, 0400)
	require.NoError(t, err)

	args := []string{
		"--dir=testdir",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.NoError(t, err)
}
