// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/veraison/corim/cots"
)

var (
	cotsCreateCtsDirs    []string
	cotsCreateCtsFiles   []string
	cotsCreateOutputFile *string
)

var cotsCreateCmd = NewCotsCreateCmd()

func NewCotsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a CBOR-encoded concise-ta-stores instance from the supplied concise-ta-store objects",
		Long: `create a CBOR-encoded concise-ta-stores instance from the supplied concise-ta-store objects,

	Create a concise-ta-stores containing stores from the cts1.cbor and cts2.cbor files.

	  cocli cots createCts --ctsFile=cts1.cbor --ctsFile=cts2.cbor
	 
	Create a concise-ta-stores containing stores from the cts directory.

	  cocli cots createCts --cts=cts
	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkcotsCreateArgs(); err != nil {
				return err
			}

			ctsFilesList := filesList(cotsCreateCtsFiles, cotsCreateCtsDirs, ".cbor")

			if len(ctsFilesList) == 0 {
				return errors.New("no CTS files found")
			}

			// checkCorimCreateArgs makes sure cotsCreateCotsFile is not nil
			cborFile, err := cotsTemplateToCBOR(ctsFilesList, cotsCreateOutputFile)
			if err != nil {
				return err
			}
			fmt.Printf(">> created %q\n", cborFile)

			return nil
		},
	}
	cmd.Flags().StringArrayVarP(
		&cotsCreateCtsDirs, "cts", "c", []string{}, "a directory containing binary CBOR-encoded concise-ta-store-map files",
	)
	cmd.Flags().StringArrayVarP(
		&cotsCreateCtsFiles, "ctsfile", "", []string{}, "a CBOR-encoded concise-ta-store-map file",
	)

	cotsCreateOutputFile = cmd.Flags().StringP("output", "o", "", "name of the generated (unsigned) CoTS file")

	return cmd
}

func checkcotsCreateArgs() error {
	if len(cotsCreateCtsFiles)+len(cotsCreateCtsDirs) == 0 {
		return errors.New("no CTS files or folders supplied")
	}

	return nil
}

func cotsTemplateToCBOR(ctsFiles []string, outputFile *string) (string, error) {
	var (
		err     error
		ctsCBOR []byte
		ctsFile string
	)

	tastores := cots.ConciseTaStores{}

	for _, ctsFile := range ctsFiles {
		var (
			ctsData []byte
			cts     cots.ConciseTaStore
		)

		ctsData, err = afero.ReadFile(fs, ctsFile)
		if err != nil {
			return "", fmt.Errorf("error loading CTS from %s: %w", ctsFile, err)
		}

		err = cts.FromCBOR(ctsData)
		if err != nil {
			return "", fmt.Errorf("failed to parse as CBOR from %s: %w", ctsFile, err)
		}
		tastores.AddConciseTaStores(cts)
	}

	// check the result
	if err = tastores.Valid(); err != nil {
		return "", fmt.Errorf("error validating COTS: %w", err)
	}

	ctsCBOR, err = tastores.ToCBOR()
	if err != nil {
		return "", fmt.Errorf("error encoding CoTS to CBOR: %w", err)
	}

	if outputFile == nil || *outputFile == "" {
		ctsFile = makeFileName("", "cts", ".cbor")
	} else {
		ctsFile = *outputFile
	}

	err = afero.WriteFile(fs, ctsFile, ctsCBOR, 0644)
	if err != nil {
		return "", fmt.Errorf("error saving CoTS to file %s: %w", ctsFile, err)
	}

	return ctsFile, nil
}

func init() {
	cotsCmd.AddCommand(cotsCreateCmd)
}
