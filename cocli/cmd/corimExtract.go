// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/veraison/corim/corim"
	"github.com/veraison/corim/cots"
)

var (
	corimExtractCorimFile *string
	corimExtractOutputDir *string
)

var corimExtractCmd = NewCorimExtractCmd()

func NewCorimExtractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract",
		Short: "extract, as-is, CoSWIDs CoMIDs, CoTS found in a CoRIM and save them to disk",
		Long: `extract, as-is, CoSWIDs and CoMIDs, CoTS found in a CoRIM and save them to disk

	Extract the contents of the signed CoRIM signed-corim.cbor to the current
	directory
	
	  cocli corim extract --file=signed-corim.cbor

	Extract the contents of the signed CoRIM yet-another-signed-corim.cbor and
	store them to directory my-dir.  Note that my-dir must exist.
	
	  cocli corim extract --file=yet-another-signed-corim.cbor \
	    				--output-dir=my-dir
	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkCorimExtractArgs(); err != nil {
				return err
			}

			return extract(*corimExtractCorimFile, corimExtractOutputDir)
		},
	}

	corimExtractCorimFile = cmd.Flags().StringP("file", "f", "", "a signed CoRIM file (in CBOR format)")
	corimExtractOutputDir = cmd.Flags().StringP("output-dir", "o", ".", "folder to which CoSWIDs, CoMIDs, CoTSs are saved")

	return cmd
}

func checkCorimExtractArgs() error {
	if corimExtractCorimFile == nil || *corimExtractCorimFile == "" {
		return errors.New("no CoRIM supplied")
	}

	return nil
}

func extract(signedCorimFile string, outputDir *string) error {
	var (
		signedCorimCBOR []byte
		err             error
		s               corim.SignedCorim
		baseDir         string
	)

	if signedCorimCBOR, err = afero.ReadFile(fs, signedCorimFile); err != nil {
		return fmt.Errorf("error loading signed CoRIM from %s: %w", signedCorimFile, err)
	}

	if err = s.FromCOSE(signedCorimCBOR); err != nil {
		return fmt.Errorf("error decoding signed CoRIM from %s: %w", signedCorimFile, err)
	}

	baseDir = "."
	if outputDir != nil {
		baseDir = *outputDir
	}

	for i, e := range s.UnsignedCorim.Tags {
		var (
			outputFile string
		)

		// need at least 3 bytes for the tag and 1 for the smallest bstr
		if len(e) < 3+1 {
			fmt.Printf(">> skipping malformed tag at index %d\n", i)
			continue
		}

		// split tag from data
		cborTag, cborData := e[:3], e[3:]

		if bytes.Equal(cborTag, corim.ComidTag) {
			outputFile = filepath.Join(baseDir, fmt.Sprintf("%06d-comid.cbor", i))

			if err = afero.WriteFile(fs, outputFile, cborData, 0644); err != nil {
				fmt.Printf(">> error saving CoMID tag at index %d: %v\n", i, err)
			}
		} else if bytes.Equal(cborTag, corim.CoswidTag) {
			outputFile = filepath.Join(baseDir, fmt.Sprintf("%06d-coswid.cbor", i))

			if err = afero.WriteFile(fs, outputFile, cborData, 0644); err != nil {
				fmt.Printf(">> error saving CoSWID tag at index %d: %v\n", i, err)
			}
		} else if bytes.Equal(cborTag, cots.CotsTag) {
			outputFile = filepath.Join(baseDir, fmt.Sprintf("%06d-cots.cbor", i))

			if err = afero.WriteFile(fs, outputFile, cborData, 0644); err != nil {
				fmt.Printf(">> error saving CoTS tag at index %d: %v\n", i, err)
			}
		} else {
			fmt.Printf(">> unmatched CBOR tag: %x\n", cborTag)
		}
	}

	return nil
}

func init() {
	corimCmd.AddCommand(corimExtractCmd)
}
