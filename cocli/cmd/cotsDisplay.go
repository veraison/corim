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
	"github.com/veraison/corim/corim"
	"github.com/veraison/corim/cots"
)

var (
	cotsDisplayCorimFile *string
	cotsDisplayShowTags  *bool
)

var cotsDisplayCmd = NewCotsDisplayCmd()

func NewCotsDisplayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "display",
		Short: "display the content of an unsigned CoTS as JSON",
		Long: `display the content of an unsigned CoTS as JSON

	Display the contents of an unsigned CoTS unsigned-cots.cbor 
	
	  cocli cots display --file unsigned-cots.cbor

	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkCotsDisplayArgs(); err != nil {
				return err
			}

			if err := cotsDisplay(*cotsDisplayCorimFile, *cotsDisplayShowTags); err != nil {
				return err
			}

			return nil
		},
	}

	cotsDisplayCorimFile = cmd.Flags().StringP("file", "f", "", "an unsigned CoTS file (in CBOR format)")
	cotsDisplayShowTags = cmd.Flags().BoolP("show-tags", "v", false, "display embedded tags")

	return cmd
}

func checkCotsDisplayArgs() error {
	if cotsDisplayCorimFile == nil || *cotsDisplayCorimFile == "" {
		return errors.New("no CoTS supplied")
	}

	return nil
}

func cotsDisplay(signedCotsFile string, showTags bool) error {
	var (
		cotsData []byte
		metaJSON []byte
		cotsJSON []byte
		err      error
		s        corim.SignedCorim
	)

	if cotsData, err = afero.ReadFile(fs, signedCotsFile); err != nil {
		return fmt.Errorf("error loading signed CoTS from %s: %w", signedCotsFile, err)
	}

	if err = s.FromCOSE(cotsData); err != nil {
		return fmt.Errorf("error decoding signed CoTS from %s: %w", signedCotsFile, err)
	}

	if metaJSON, err = json.MarshalIndent(&s.Meta, "", "  "); err != nil {
		return fmt.Errorf("error decoding CoRIM Meta from %s: %w", signedCotsFile, err)
	}

	fmt.Println("Meta:")
	fmt.Println(string(metaJSON))

	if cotsJSON, err = json.MarshalIndent(&s.UnsignedCorim, "", "  "); err != nil {
		return fmt.Errorf("error decoding unsigned CoTS from %s: %w", signedCotsFile, err)
	}

	fmt.Println("Cots:")
	fmt.Println(string(cotsJSON))

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

			if bytes.Equal(cborTag, cots.CotsTag) {
				if err = printCots(cborData, hdr); err != nil {
					fmt.Printf(">> skipping malformed CoMID tag at index %d: %v\n", i, err)
				}
			} else {
				fmt.Printf(">> unmatched CBOR tag: %x\n", cborTag)
			}
		}
	}
	return nil
}

func printCots(cbor []byte, heading string) error {
	return printJSONFromCBOR(&cots.ConciseTaStores{}, cbor, heading)
}

func init() {
	cotsCmd.AddCommand(cotsDisplayCmd)
}
