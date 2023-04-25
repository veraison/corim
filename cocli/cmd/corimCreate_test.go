// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CorimCreateCmd_unknown_argument(t *testing.T) {
	cmd := NewCorimCreateCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_CorimCreateCmd_no_templates(t *testing.T) {
	cmd := NewCorimCreateCmd()

	// no args

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoRIM template supplied")
}

func Test_CorimCreateCmd_no_files_found(t *testing.T) {
	cmd := NewCorimCreateCmd()

	args := []string{
		"--template=unknown.json",
		"--comid=unsure.cbor",
		"--comid-dir=somedir",
		"--coswid=what.cbor",
		"--coswid-dir=someotherdir",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoMID, CoSWID or CoTS files found")
}

func Test_CorimCreateCmd_no_tag_files(t *testing.T) {
	cmd := NewCorimCreateCmd()

	args := []string{
		"--template=unknown.json",
		// no --co{m,sw}id,{cots}-{dir,}
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoMID, CoSWID or CoTS files or folders supplied")
}

func Test_CorimCreateCmd_template_not_found(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "ignored-comid.cbor", []byte{}, 0644)
	require.NoError(t, err)

	args := []string{
		"--template=nonexistent.json",
		"--comid=ignored-comid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, "error loading template from nonexistent.json: open nonexistent.json: file does not exist")
}

func Test_CorimCreateCmd_template_with_invalid_json(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "invalid.json", []byte("..."), 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ignored-comid.cbor", []byte{}, 0644)
	require.NoError(t, err)

	args := []string{
		"--template=invalid.json",
		"--comid=ignored-comid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, "error decoding template from invalid.json: invalid character '.' looking for beginning of value")
}

func Test_CorimCreateCmd_with_a_bad_comid(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "min-tmpl.json", minimalCorimTemplate, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "bad-comid.cbor", badCBOR, 0644)
	require.NoError(t, err)

	args := []string{
		"--template=min-tmpl.json",
		"--comid=bad-comid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, `error loading CoMID from bad-comid.cbor: cbor: unexpected "break" code`)
}

func Test_CorimCreateCmd_with_an_invalid_comid(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "min-tmpl.json", minimalCorimTemplate, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "invalid-comid.cbor", invalidComid, 0644)
	require.NoError(t, err)

	args := []string{
		"--template=min-tmpl.json",
		"--comid=invalid-comid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, `error adding CoMID from invalid-comid.cbor (check its validity using the "comid validate" sub-command)`)
}

func Test_CorimCreateCmd_with_a_bad_coswid(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "min-tmpl.json", minimalCorimTemplate, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "bad-coswid.cbor", badCBOR, 0644)
	require.NoError(t, err)

	args := []string{
		"--template=min-tmpl.json",
		"--coswid=bad-coswid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, `error loading CoSWID from bad-coswid.cbor: cbor: unexpected "break" code`)
}

func Test_CorimCreateCmd_with_an_invalid_cots(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "min-tmpl.json", minimalCorimTemplate, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "invalid-cots.cbor", invalidCots, 0644)
	require.NoError(t, err)

	args := []string{
		"--template=min-tmpl.json",
		"--cots=invalid-cots.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, `error adding CoTS from invalid-cots.cbor`)
}

func Test_CorimCreateCmd_with_a_bad_cots(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "min-tmpl.json", minimalCorimTemplate, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "bad-cots.cbor", badCBOR, 0644)
	require.NoError(t, err)

	args := []string{
		"--template=min-tmpl.json",
		"--cots=bad-cots.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, `error loading CoTS from bad-cots.cbor: cbor: unexpected "break" code`)
}

func Test_CorimCreateCmd_successful_comid_coswid_and_cots_from_file(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "min-tmpl.json", minimalCorimTemplate, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "coswid.cbor", testCoswid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "comid.cbor", testComid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "cots.cbor", testCots, 0644)
	require.NoError(t, err)

	args := []string{
		"--template=min-tmpl.json",
		"--coswid=coswid.cbor",
		"--comid=comid.cbor",
		"--cots=cots.cbor",
		"--output=corim.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.NoError(t, err)

	_, err = fs.Stat("corim.cbor")
	assert.NoError(t, err)
}

func Test_CorimCreateCmd_successful_comid_coswid_and_cots_from_dir(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "min-tmpl.json", minimalCorimTemplate, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "coswid/1.cbor", testCoswid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "comid/1.cbor", testComid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "cots/1.cbor", testCots, 0644)
	require.NoError(t, err)

	args := []string{
		"--template=min-tmpl.json",
		"--coswid-dir=coswid",
		"--comid-dir=comid",
		"--cots-dir=cots",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.NoError(t, err)

	_, err = fs.Stat("min-tmpl.cbor")
	assert.NoError(t, err)
}
