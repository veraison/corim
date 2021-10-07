// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
)

func Test_ComidCreateCmd_unknown_argument(t *testing.T) {
	cmd := NewComidCreateCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_ComidCreateCmd_no_templates(t *testing.T) {
	cmd := NewComidCreateCmd()

	// no args

	err := cmd.Execute()
	assert.EqualError(t, err, "no templates supplied")
}

func Test_ComidCreateCmd_no_files_found(t *testing.T) {
	cmd := NewComidCreateCmd()

	args := []string{
		"--template=unknown",
		"--template-dir=unsure",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no files found")
}

func Test_ComidCreateCmd_template_with_invalid_json(t *testing.T) {
	var err error

	cmd := NewComidCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "invalid.json", []byte("..."), 0644)
	require.NoError(t, err)

	args := []string{
		"--template=invalid.json",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, "error decoding template from invalid.json: invalid character '.' looking for beginning of value")
}

func Test_ComidCreateCmd_template_with_invalid_comid(t *testing.T) {
	var err error

	cmd := NewComidCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "bad-comid.json", []byte("{}"), 0644)
	require.NoError(t, err)

	args := []string{
		"--template=bad-comid.json",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, "error validating template bad-comid.json: tag-identity validation failed: empty tag-id")
}

func Test_ComidCreateCmd_template_from_file_to_default_dir(t *testing.T) {
	var err error

	cmd := NewComidCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "ok.json", []byte(comid.PSARefValJSONTemplate), 0644)
	require.NoError(t, err)

	args := []string{
		"--template=ok.json",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.NoError(t, err)

	expectedFileName := "ok.cbor"

	_, err = fs.Stat(expectedFileName)
	assert.NoError(t, err)
}

func Test_ComidCreateCmd_template_from_dir_to_custom_dir(t *testing.T) {
	var err error

	cmd := NewComidCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "testdir/ok.json", []byte(comid.PSARefValJSONTemplate), 0644)
	require.NoError(t, err)

	args := []string{
		"--template-dir=testdir",
		"--output-dir=testdir",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.NoError(t, err)

	expectedFileName := "testdir/ok.cbor"

	_, err = fs.Stat(expectedFileName)
	assert.NoError(t, err)
}
