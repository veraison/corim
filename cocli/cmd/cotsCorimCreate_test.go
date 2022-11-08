// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CotsCreateCorimCmd_unknown_argument(t *testing.T) {
	cmd := NewCotsCorimCreateCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_CotsCreateCorimCmd_no_templates(t *testing.T) {
	cmd := NewCotsCorimCreateCmd()

	// no args

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoRIM template supplied")
}

func Test_CotsCreateCorimCmd_no_files_found(t *testing.T) {
	cmd := NewCotsCorimCreateCmd()

	args := []string{
		"--template=unknown.json",
		"--cots=unsure.cbor",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoTS files found")
}

func Test_CotsCreateCorimCmd_no_tag_files(t *testing.T) {
	cmd := NewCotsCorimCreateCmd()

	args := []string{
		"--template=unknown.json",
		// no --co{m,sw}id-{dir,}
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoTS supplied")
}

func Test_CotsCreateCorimCmd_template_not_found(t *testing.T) {
	var err error

	cmd := NewCotsCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "ignored-comid.cbor", []byte{}, 0644)
	require.NoError(t, err)

	args := []string{
		"--template=nonexistent.json",
		"--cots=ignored-comid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, "error loading template from nonexistent.json: open nonexistent.json: file does not exist")
}

func Test_CotsCreateCorimCmd_template_with_invalid_json(t *testing.T) {
	var err error

	cmd := NewCotsCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "invalid.json", []byte("..."), 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ignored-comid.cbor", []byte{}, 0644)
	require.NoError(t, err)

	args := []string{
		"--template=invalid.json",
		"--cots=ignored-comid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, "error decoding template from invalid.json: invalid character '.' looking for beginning of value")
}

func Test_CotsCreateCorimCmd_with_a_bad_cots(t *testing.T) {
	var err error

	cmd := NewCotsCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "min-tmpl.json", minimalCorimTemplate, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "bad-comid.cbor", badCBOR, 0644)
	require.NoError(t, err)

	args := []string{
		"--template=min-tmpl.json",
		"--cots=bad-comid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, `error loading CoTS from bad-comid.cbor: cbor: unexpected "break" code`)
}

func Test_CotsCreateCorimCmd_with_an_invalid_cots(t *testing.T) {
	var err error

	cmd := NewCotsCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "min-tmpl.json", minimalCorimTemplate, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "invalid-cots.cbor", invalidComid, 0644)
	require.NoError(t, err)

	args := []string{
		"--template=min-tmpl.json",
		"--cots=invalid-cots.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, `error loading CoTS from invalid-cots.cbor: cbor: cannot unmarshal map into Go value of type cots.ConciseTaStores`)
}

func Test_CotsCreateCorimCmd_successful_cots_from_file(t *testing.T) {
	var err error

	cmd := NewCotsCorimCreateCmd()

	fs = afero.NewOsFs()

	args := []string{
		"--template=../data/templates/meta-full.json",
		"--cots=../data/cots/cots.cbor",
		"--output=corim.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.NoError(t, err)

	_, err = fs.Stat("corim.cbor")
	assert.NoError(t, err)
	e := os.Remove("corim.cbor")
	assert.NoError(t, e)
}
