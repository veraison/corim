// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/veraison/corim/cots"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/veraison/corim/corim"
)

var (
	cotsCorimCreateCorimFile  *string
	cotsCorimCreateCotsFile   *string
	cotsCorimCreateOutputFile *string
	cotsCorimCreateCotsFiles  []string
)

var cotsCorimCreateCmd = NewCotsCorimCreateCmd()

func NewCotsCorimCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createCorim",
		Short: "create a CBOR-encoded CoRIM from the supplied JSON template and CoTS",
		Long: `create a CBOR-encoded CoRIM from the supplied JSON template and CoTS,

	Create a CoRIM from template t1.json, adding CoMIDs found in the comid/
	directory and CoSWIDs found in the coswid/ directory.  Since no explicit
	output file is set, the (unsigned) CoRIM is saved to the current directory
	with tag-id as basename and a .cbor extension.

	  cocli cots createCorim --template=t1.json --cots=cots.cbor
	 
	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkCotsCorimCreateArgs(); err != nil {
				return err
			}

			cotsCorimCreateCotsFiles = append(cotsCorimCreateCotsFiles, *cotsCorimCreateCotsFile)
			cotsFilesList := filesList(cotsCorimCreateCotsFiles, nil, ".cbor")

			if len(cotsFilesList) == 0 {
				return errors.New("no CoTS files found")
			}

			// checkCorimCreateArgs makes sure cotsCorimCreateCorimFile is not nil
			cborFile, err := cotsCorimTemplateToCBOR(*cotsCorimCreateCorimFile,
				*cotsCorimCreateCotsFile, cotsCorimCreateOutputFile)
			if err != nil {
				return err
			}
			fmt.Printf(">> created %q from %q\n", cborFile, *cotsCorimCreateCorimFile)

			return nil
		},
	}

	cotsCorimCreateCorimFile = cmd.Flags().StringP("template", "t", "", "a CoRIM template file (in JSON format)")

	cotsCorimCreateCotsFile = cmd.Flags().StringP("cots", "c", "", "a CoTS file (in CBOR format)")

	cotsCorimCreateOutputFile = cmd.Flags().StringP("output", "o", "", "name of the generated (unsigned) CoRIM file")

	return cmd
}

func checkCotsCorimCreateArgs() error {
	if cotsCorimCreateCorimFile == nil || *cotsCorimCreateCorimFile == "" {
		return errors.New("no CoRIM template supplied")
	}
	if cotsCorimCreateCotsFile == nil || *cotsCorimCreateCotsFile == "" {
		return errors.New("no CoTS supplied")
	}

	return nil
}

func cotsCorimTemplateToCBOR(tmplFile string, cotsFile string, outputFile *string) (string, error) {
	var (
		tmplData, corimCBOR []byte
		c                   corim.UnsignedCorim
		corimFile           string
		err                 error
	)

	if tmplData, err = afero.ReadFile(fs, tmplFile); err != nil {
		return "", fmt.Errorf("error loading template from %s: %w", tmplFile, err)
	}

	if err = c.FromJSON(tmplData); err != nil {
		return "", fmt.Errorf("error decoding template from %s: %w", tmplFile, err)
	}

	// append CoMID(s)
	var (
		comidCBOR []byte
		m         cots.ConciseTaStores
	)

	comidCBOR, err = afero.ReadFile(fs, cotsFile)
	if err != nil {
		return "", fmt.Errorf("error loading CoTS from %s: %w", cotsFile, err)
	}

	err = m.FromCBOR(comidCBOR)
	if err != nil {
		return "", fmt.Errorf("error loading CoTS from %s: %w", cotsFile, err)
	}

	if c.AddCots(m) == nil {
		return "", fmt.Errorf(
			"error adding CoTS from %s (check its validity using the %q sub-command)",
			cotsFile, "comid validate",
		)
	}

	c.SetID(uuid.New().String())

	// check the result
	if err = c.Valid(); err != nil {
		return "", fmt.Errorf("error validating CoRIM: %w", err)
	}

	corimCBOR, err = c.ToCBOR()
	if err != nil {
		return "", fmt.Errorf("error encoding CoRIM to CBOR: %w", err)
	}

	if outputFile == nil || *outputFile == "" {
		corimFile = makeFileName("", tmplFile, ".cbor")
	} else {
		corimFile = *outputFile
	}

	err = afero.WriteFile(fs, corimFile, corimCBOR, 0644)
	if err != nil {
		return "", fmt.Errorf("error saving CoRIM to file %s: %w", corimFile, err)
	}

	return corimFile, nil
}

func init() {
	cotsCmd.AddCommand(cotsCorimCreateCmd)
}
