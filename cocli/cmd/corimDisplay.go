// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/veraison/corim/v2/corim"
	"github.com/veraison/corim/v2/cots"
)

var (
	corimDisplayCorimFile *string
	corimDisplayShowTags  *bool
)

var corimDisplayCmd = NewCorimDisplayCmd()

func NewCorimDisplayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "display",
		Short: "display the content of a CoRIM as JSON",
		Long: `display the content of a CoRIM as JSON

	Display the contents of the signed CoRIM signed-corim.cbor 
	
	  cocli corim display --file signed-corim.cbor

	Display the contents of the signed CoRIM yet-another-signed-corim.cbor and
	also unpack any embedded CoMID, CoSWID and CoTS
	
	  cocli corim display --file yet-another-signed-corim.cbor --show-tags
	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkCorimDisplayArgs(); err != nil {
				return err
			}

			if err := display(*corimDisplayCorimFile, *corimDisplayShowTags); err != nil {
				return err
			}

			return nil
		},
	}

	corimDisplayCorimFile = cmd.Flags().StringP("file", "f", "", "a signed CoRIM file (in CBOR format)")
	corimDisplayShowTags = cmd.Flags().BoolP("show-tags", "v", false, "display embedded tags")

	return cmd
}

func checkCorimDisplayArgs() error {
	if corimDisplayCorimFile == nil || *corimDisplayCorimFile == "" {
		return errors.New("no CoRIM supplied")
	}

	return nil
}

func display(signedCorimFile string, showTags bool) error {
	var (
		signedCorimCBOR []byte
		metaJSON        []byte
		corimJSON       []byte
		err             error
		s               corim.SignedCorim
	)

	if signedCorimCBOR, err = afero.ReadFile(fs, signedCorimFile); err != nil {
		return fmt.Errorf("error loading signed CoRIM from %s: %w", signedCorimFile, err)
	}

	if err = s.FromCOSE(signedCorimCBOR); err != nil {
		return fmt.Errorf("error decoding signed CoRIM from %s: %w", signedCorimFile, err)
	}

	if metaJSON, err = json.MarshalIndent(&s.Meta, "", "  "); err != nil {
		return fmt.Errorf("error decoding CoRIM Meta from %s: %w", signedCorimFile, err)
	}

	fmt.Println("Meta:")
	fmt.Println(string(metaJSON))

	if corimJSON, err = json.MarshalIndent(&s.UnsignedCorim, "", "  "); err != nil {
		return fmt.Errorf("error decoding unsigned CoRIM from %s: %w", signedCorimFile, err)
	}

	fmt.Println("Corim:")
	fmt.Println(string(corimJSON))

	if showTags {
		fmt.Println("Tags:")
		for i, e := range s.UnsignedCorim.Tags {
			// need at least 3 bytes for the tag and 1 for the smallest bstr
			if len(e) < 3+1 {
				fmt.Printf(">> skipping malformed tag at index %d\n", i)
				continue
			}

			// split tag from data
			cborTag, cborData := e[:3], e[3:]

			hdr := fmt.Sprintf(">> [ %d ]", i)

			if bytes.Equal(cborTag, corim.ComidTag) {
				if err = printComid(cborData, hdr); err != nil {
					fmt.Printf(">> skipping malformed CoMID tag at index %d: %v\n", i, err)
				}
			} else if bytes.Equal(cborTag, corim.CoswidTag) {
				if err = printCoswid(cborData, hdr); err != nil {
					fmt.Printf(">> skipping malformed CoSWID tag at index %d: %v\n", i, err)
				}
			} else if bytes.Equal(cborTag, cots.CotsTag) {
				if err = printCots(cborData, hdr); err != nil {
					fmt.Printf(">> skipping malformed CoTS tag at index %d: %v\n", i, err)
				}
			} else {
				fmt.Printf(">> unmatched CBOR tag: %x\n", cborTag)
			}
		}
	}

	return nil
}

func init() {
	corimCmd.AddCommand(corimDisplayCmd)
}
