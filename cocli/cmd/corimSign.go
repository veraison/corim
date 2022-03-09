// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/veraison/corim/corim"
	cose "github.com/veraison/go-cose"
)

var (
	corimSignCorimFile  *string
	corimSignKeyFile    *string
	corimSignOutputFile *string
	corimSignMetaFile   *string
)

var corimSignCmd = NewCorimSignCmd()

func NewCorimSignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign",
		Short: "create a signed CoRIM from an unsigned, CBOR-encoded CoRIM using the supplied key",
		Long: `create a signed CoRIM from an unsigned, CBOR-encoded CoRIM using the supplied key

	Sign the unsigned CoRIM unsigned-corim.cbor using the key in JWK format from
	file key.jwk and save the resulting COSE Sign1 to signed-corim.cbor.  Read
	the relevant CorimMeta information from file meta.json.
	
	  cocli corim sign  --file=unsigned-corim.cbor \
					--key=key.jwk \
					--meta=meta.json \
					--output=signed-corim.cbor
	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkCorimSignArgs(); err != nil {
				return err
			}

			// checkCorimSignArgs makes sure corimSignCorimFile is not nil
			coseFile, err := sign(*corimSignCorimFile, *corimSignKeyFile,
				*corimSignMetaFile, corimSignOutputFile)
			if err != nil {
				return err
			}
			fmt.Printf(">> %q signed and saved to %q\n", *corimSignCorimFile, coseFile)

			return nil
		},
	}

	corimSignCorimFile = cmd.Flags().StringP("file", "f", "", "an unsigned CoRIM file (in CBOR format)")
	corimSignMetaFile = cmd.Flags().StringP("meta", "m", "", "CoRIM Meta file (in JSON format)")
	corimSignKeyFile = cmd.Flags().StringP("key", "k", "", "signing key in JWK format")
	corimSignOutputFile = cmd.Flags().StringP("output", "o", "", "name of the generated COSE Sign1 file")

	return cmd
}

func checkCorimSignArgs() error {
	if corimSignCorimFile == nil || *corimSignCorimFile == "" {
		return errors.New("no CoRIM supplied")
	}

	if corimSignKeyFile == nil || *corimSignKeyFile == "" {
		return errors.New("no key supplied")
	}

	if corimSignMetaFile == nil || *corimSignMetaFile == "" {
		return errors.New("no CoRIM Meta supplied")
	}

	return nil
}

func sign(unsignedCorimFile, keyFile, metaFile string, outputFile *string) (string, error) {
	var (
		unsignedCorimCBOR []byte
		signedCorimCBOR   []byte
		metaJSON          []byte
		keyJWK            []byte
		err               error
		signedCorimFile   string
		c                 corim.UnsignedCorim
		m                 corim.Meta
		signer            *cose.Signer
	)

	if unsignedCorimCBOR, err = afero.ReadFile(fs, unsignedCorimFile); err != nil {
		return "", fmt.Errorf("error loading unsigned CoRIM from %s: %w", unsignedCorimFile, err)
	}

	if err = c.FromCBOR(unsignedCorimCBOR); err != nil {
		return "", fmt.Errorf("error decoding unsigned CoRIM from %s: %w", unsignedCorimFile, err)
	}

	if err = c.Valid(); err != nil {
		return "", fmt.Errorf("error validating CoRIM: %w", err)
	}

	if metaJSON, err = afero.ReadFile(fs, metaFile); err != nil {
		return "", fmt.Errorf("error loading CoRIM Meta from %s: %w", metaFile, err)
	}

	if err = m.FromJSON(metaJSON); err != nil {
		return "", fmt.Errorf("error decoding CoRIM Meta from %s: %w", metaFile, err)
	}

	if err = m.Valid(); err != nil {
		return "", fmt.Errorf("error validating CoRIM Meta: %w", err)
	}

	if keyJWK, err = afero.ReadFile(fs, keyFile); err != nil {
		return "", fmt.Errorf("error loading signing key from %s: %w", keyFile, err)
	}

	if signer, err = corim.SignerFromJWK(keyJWK); err != nil {
		return "", fmt.Errorf("error loading signing key from %s: %w", keyFile, err)
	}

	s := corim.SignedCorim{
		UnsignedCorim: c,
		Meta:          m,
	}

	signedCorimCBOR, err = s.Sign(signer)
	if err != nil {
		return "", fmt.Errorf("error signing CoRIM: %w", err)
	}

	if outputFile == nil || *outputFile == "" {
		signedCorimFile = "signed-" + unsignedCorimFile
	} else {
		signedCorimFile = *outputFile
	}

	err = afero.WriteFile(fs, signedCorimFile, signedCorimCBOR, 0644)
	if err != nil {
		return "", fmt.Errorf("error saving signed CoRIM to file %s: %w", signedCorimFile, err)
	}

	return signedCorimFile, nil
}

func init() {
	corimCmd.AddCommand(corimSignCmd)
}
