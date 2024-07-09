// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/veraison/corim/comid"
)

var (
	comidCreateFiles     []string
	comidCreateDirs      []string
	comidCreateOutputDir string
)

var comidCreateCmd = NewComidCreateCmd()

func NewComidCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create one or more CBOR-encoded CoMID(s) from the supplied JSON template(s)",
		Long: `create one or more CBOR-encoded CoMID(s) from the supplied JSON template(s)

	Create CoMIDs from templates t1.json and t2.json, plus any template found in
	the templates/ directory.  Save them to the current working directory.
	
		cocli comid create --template=t1.json \
	    			--template=t2.json \
	    			--template-dir=templates
	  
	Create one CoMID from template t3.json and save it to the comids/ directory.
	Note that the output directory must exist.
	
		cocli comid create --template=t3.json --output-dir=comids

	Note: since the output file is deterministically generated from the template
	file name, all the template file names (when from different directories)
	MUST be different.
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkComidCreateArgs(); err != nil {
				return err
			}

			filesList := filesList(comidCreateFiles, comidCreateDirs, ".json")
			if len(filesList) == 0 {
				return errors.New("no files found")
			}

			errs := 0
			for _, tmplFile := range filesList {
				cborFile, err := templateToCBOR(tmplFile, comidCreateOutputDir)
				if err != nil {
					fmt.Printf(">> creation failed for %q: %v\n", cborFile, err)
					errs++
					continue
				}
				fmt.Printf(">> created %q from %q\n", cborFile, tmplFile)
			}

			if errs != 0 {
				return fmt.Errorf("%d/%d creations(s) failed", errs, len(filesList))
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(
		&comidCreateFiles, "template", "t", []string{}, "a CoMID template file (in JSON format)",
	)

	cmd.Flags().StringArrayVarP(
		&comidCreateDirs, "template-dir", "T", []string{}, "a directory containing CoMID template files",
	)

	cmd.Flags().StringVarP(
		&comidCreateOutputDir, "output-dir", "o", ".", "directory where the created files are stored",
	)

	return cmd
}

func checkComidCreateArgs() error {
	if len(comidCreateFiles) == 0 && len(comidCreateDirs) == 0 {
		return errors.New("no templates supplied")
	}
	return nil
}

func templateToCBOR(tmplFile, outputDir string) (string, error) {
	var (
		tmplData, cborData []byte
		cborFile           string
		c                  comid.Comid
		err                error
	)

	if tmplData, err = afero.ReadFile(fs, tmplFile); err != nil {
		return "", fmt.Errorf("error loading template from %s: %w", tmplFile, err)
	}

	if err = c.FromJSON(tmplData); err != nil {
		return "", fmt.Errorf("error decoding template from %s: %w", tmplFile, err)
	}

	if err = c.Valid(); err != nil {
		return "", fmt.Errorf("error validating template %s: %w", tmplFile, err)
	}

	cborData, err = c.ToCBOR()
	if err != nil {
		return "", fmt.Errorf("error encoding template %s to CBOR: %w", tmplFile, err)
	}

	cborFile = makeFileName(outputDir, tmplFile, ".cbor")

	err = afero.WriteFile(fs, cborFile, cborData, 0644)
	if err != nil {
		return "", fmt.Errorf("error saving CBOR file %s: %w", cborFile, err)
	}

	return cborFile, nil
}

func init() {
	comidCmd.AddCommand(comidCreateCmd)
}
