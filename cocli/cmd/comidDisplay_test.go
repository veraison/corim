// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ComidDisplayCmd_unknown_argument(t *testing.T) {
	cmd := NewComidDisplayCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_ComidDisplayCmd_no_files(t *testing.T) {
	cmd := NewComidDisplayCmd()

	// no args

	err := cmd.Execute()
	assert.EqualError(t, err, "no files supplied")
}

func Test_ComidDisplayCmd_no_files_found(t *testing.T) {
	cmd := NewComidDisplayCmd()

	args := []string{
		"--file=unknown",
		"--dir=unsure",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no files found")
}

func Test_ComidDisplayCmd_file_with_invalid_cbor(t *testing.T) {
	var err error

	cmd := NewComidDisplayCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "invalid.cbor", []byte{0xff, 0xff}, 0400)
	require.NoError(t, err)

	args := []string{
		"--file=invalid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, "1/1 display(s) failed")
}

func Test_ComidDisplayCmd_file_with_valid_comid(t *testing.T) {
	var err error

	cmd := NewComidDisplayCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "ok.cbor", PSARefValCBOR, 0400)
	require.NoError(t, err)

	args := []string{
		"--file=ok.cbor",
	}
	cmd.SetArgs(args)

	fmt.Printf("%x\n", PSARefValCBOR)

	err = cmd.Execute()
	assert.NoError(t, err)
}

func Test_ComidDisplayCmd_file_with_valid_comid_from_dir(t *testing.T) {
	var err error

	cmd := NewComidDisplayCmd()

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
