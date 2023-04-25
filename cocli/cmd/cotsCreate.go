// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/veraison/corim/cots"
)

var (
	cotsCreateLanguage          *string
	cotsCreateTagID             *string
	cotsCreateTagUUIDStr        *string
	cotsCreateTagUUID           *bool
	cotsCreateTagVersion        *uint
	cotsCreateCtsEnvFile        *string
	cotsCreateCtsPermClaimsFile *string
	cotsCreateCtsExclClaimsFile *string
	cotsCreateCtsPurposes       []string
	cotsCreateCtsTaDirs         []string
	cotsCreateCtsTaFiles        []string
	cotsCreateCtsCaDirs         []string
	cotsCreateCtsCaFiles        []string
	cotsCreateCtsOutputFile     *string
)

var cotsCreateCtsCmd = NewCotsCreateCtsCmd()

func NewCotsCreateCtsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a CBOR-encoded concise-ta-store-map instance from the supplied JSON environment template, JSON constraints templates, purposes, and TAs/CAs",
		Long: `create a CBOR-encoded concise-ta-store-map instance from the supplied JSON environment template, JSON constraints templates, purposes, and TAs/CAs,

	Create a concise-ta-store-map from template env-template.json and TAs from tas directory.  Since no explicit output file is set, 
	the result is saved to the current directory with the env-template as basename and a .cbor extension.

	  cocli cots create --environment=env-template.json --tas=tas_dir
	 
	Create a concise-ta-store-map from env-template.json, claims-template.json, TAs from tas_dir directory, CAs from cas_dir directory, and 
	with eat and corim purposes asserted. The result is saved to cots.cbor.

	  cocli cots create --environment=env-template.json \
	                   --purpose=eat \
	                   --purpose=corim \
	                   --permclaims=claims-template.json \
	                   --tas=tas_dir \
	                   --cas=cas_dir \
	                   --output=cots.cbor
	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkctsCreateCtsArgs(); err != nil {
				return err
			}

			certFilesList := filesList(cotsCreateCtsTaFiles, cotsCreateCtsTaDirs, ".der")
			taiFilesList := filesList(cotsCreateCtsTaFiles, cotsCreateCtsTaDirs, ".ta")
			spkiFilesList := filesList(cotsCreateCtsTaFiles, cotsCreateCtsTaDirs, ".spki")
			tasFilesList := append(certFilesList, taiFilesList...)
			tasFilesList = append(tasFilesList, spkiFilesList...)
			casFilesList := filesList(cotsCreateCtsCaFiles, cotsCreateCtsCaDirs, ".der")

			if len(tasFilesList) == 0 {
				return errors.New("no TA files found")
			}

			cborFile, err := ctsTemplateToCBOR(*cotsCreateLanguage, *cotsCreateTagID, *cotsCreateTagUUID, *cotsCreateTagUUIDStr, cotsCreateTagVersion, *cotsCreateCtsEnvFile, *cotsCreateCtsPermClaimsFile, *cotsCreateCtsExclClaimsFile, cotsCreateCtsPurposes,
				tasFilesList, casFilesList, cotsCreateCtsOutputFile)
			if err != nil {
				return err
			}
			fmt.Printf(">> created %q\n", cborFile)

			return nil
		},
	}

	cotsCreateLanguage = cmd.Flags().StringP("language", "l", "", "language tag")
	cotsCreateTagUUIDStr = cmd.Flags().StringP("uuid-str", "", "", "string representation of a UUID to use as tag ID (mutually exclusive from --uuid and --id)")
	cotsCreateTagUUID = cmd.Flags().BoolP("uuid", "", false, "boolean indicating a random UUID value should be used as tag ID (mutually exclusive from --id and --uuid-str)")
	cotsCreateTagID = cmd.Flags().StringP("id", "", "", "string value containing a tag ID value (mutually exclusive from --uuid and --uuid-str)")
	cotsCreateTagVersion = cmd.Flags().UintP("tag-version", "", 0, "integer value indicating version of tag identity (ignored if neither --uuid nor --id are supplied)")
	cotsCreateCtsEnvFile = cmd.Flags().StringP("environment", "e", "", "an environment template file (in JSON format)")
	cotsCreateCtsPermClaimsFile = cmd.Flags().StringP("permclaims", "p", "", "a permitted claims template file (in JSON format)")
	cotsCreateCtsExclClaimsFile = cmd.Flags().StringP("exclclaims", "x", "", "an excluded claims template file (in JSON format)")

	cmd.Flags().StringArrayVarP(
		&cotsCreateCtsPurposes, "purpose", "u", []string{}, "string value indicating purpose: cots,corim,comid,coswid,eat,certificate",
	)

	cmd.Flags().StringArrayVarP(
		&cotsCreateCtsTaDirs, "tas", "t", []string{}, "a directory containing binary DER-encoded trust anchor files",
	)
	cmd.Flags().StringArrayVarP(
		&cotsCreateCtsTaFiles, "tafile", "", []string{}, "a DER-encoded trust anchor file",
	)

	cmd.Flags().StringArrayVarP(
		&cotsCreateCtsCaDirs, "cas", "c", []string{}, "a directory containing binary DER-encoded X.509 CA certificate files",
	)
	cmd.Flags().StringArrayVarP(
		&cotsCreateCtsCaFiles, "cafile", "", []string{}, "a DER-encoded certificate file",
	)

	cotsCreateCtsOutputFile = cmd.Flags().StringP("output", "o", "", "name of the generated (unsigned) CoTS file")

	return cmd
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func checkctsCreateCtsArgs() error {
	if cotsCreateCtsEnvFile == nil || *cotsCreateCtsEnvFile == "" {
		return errors.New("no environment template supplied")
	}

	if (*cotsCreateTagUUID != false && *cotsCreateTagID != "") || (*cotsCreateTagUUID != false && *cotsCreateTagUUIDStr != "") || (*cotsCreateTagUUIDStr != "" && *cotsCreateTagID != "") {
		return errors.New("only one of --uuid, --uuid-str and --id can be used at the same time")
	}

	if *cotsCreateTagUUIDStr != "" && !IsValidUUID(*cotsCreateTagUUIDStr) {
		return errors.New("--uuid-str does not contain a valid UUID")
	}

	if len(cotsCreateCtsTaFiles)+len(cotsCreateCtsTaDirs) == 0 {
		return errors.New("no TA files or folders supplied")
	}

	return nil
}

func ctsTemplateToCBOR(language string, tagID string, genUUID bool, uuidStr string, version *uint, envFile string, permClaimsFile string, exclClaimsFile string, purposes, taFiles, caFiles []string, outputFile *string) (string, error) {
	var (
		envData        []byte
		env            cots.EnvironmentGroups
		permClaimsData []byte
		permClaims     cots.EatCWTClaim
		exclClaimsData []byte
		exclClaims     cots.EatCWTClaim
		err            error
		ctsCBOR        []byte
		ctsFile        string
	)

	cts := cots.ConciseTaStore{}

	if envData, err = afero.ReadFile(fs, envFile); err != nil {
		return "", fmt.Errorf("error loading template from %s: %w", envFile, err)
	}

	if err = env.FromJSON(envData); err != nil {
		return "", fmt.Errorf("error decoding template from %s: %w", envFile, err)
	}

	cts.Environments = env

	if language != "" {
		cts.Language = &language
	}

	if tagID != "" {
		cts.SetTagIdentity(tagID, version)
	} else if genUUID != false {
		u := uuid.New()
		b, _ := u.MarshalBinary()
		cts.SetTagIdentity(b, version)
	} else if uuidStr != "" {
		u, _ := uuid.Parse(uuidStr)
		b, _ := u.MarshalBinary()
		cts.SetTagIdentity(b, version)
	}

	if permClaimsFile != "" {
		if permClaimsData, err = afero.ReadFile(fs, permClaimsFile); err != nil {
			return "", fmt.Errorf("error loading template from %s: %w", permClaimsFile, err)
		}

		if err = permClaims.FromJSON(permClaimsData); err != nil {
			return "", fmt.Errorf("error decoding template from %s: %w", permClaimsFile, err)
		}
		cts.AddPermClaims(permClaims)
	}
	if exclClaimsFile != "" {
		if exclClaimsData, err = afero.ReadFile(fs, exclClaimsFile); err != nil {
			return "", fmt.Errorf("error loading template from %s: %w", exclClaimsFile, err)
		}

		if err = exclClaims.FromJSON(exclClaimsData); err != nil {
			return "", fmt.Errorf("error decoding template from %s: %w", exclClaimsFile, err)
		}
		cts.AddExclClaims(exclClaims)
	}
	if 0 != len(purposes) {
		cts.Purposes = purposes
	}

	k := cots.TasAndCas{}
	cts.Keys = &k
	for _, taFile := range taFiles {
		var (
			tadata      []byte
			trustAnchor cots.TrustAnchor
		)

		tadata, err = afero.ReadFile(fs, taFile)
		if err != nil {
			return "", fmt.Errorf("error loading TA from %s: %w", taFile, err)
		}
		if ".der" == filepath.Ext(taFile) {
			trustAnchor.Format = cots.TaFormatCertificate
		}
		if ".spki" == filepath.Ext(taFile) {
			trustAnchor.Format = cots.TaFormatSubjectPublicKeyInfo
		}
		if ".ta" == filepath.Ext(taFile) {
			trustAnchor.Format = cots.TaFormatTrustAnchorInfo
		}
		trustAnchor.Data = tadata
		cts.Keys.Tas = append(cts.Keys.Tas, trustAnchor)
	}

	for _, caFile := range caFiles {
		var (
			cadata []byte
		)

		cadata, err = afero.ReadFile(fs, caFile)
		if err != nil {
			return "", fmt.Errorf("error loading CA from %s: %w", caFile, err)
		}
		cts.Keys.Cas = append(cts.Keys.Cas, cadata)
	}

	// check the result
	if err = cts.Valid(); err != nil {
		return "", fmt.Errorf("error validating CoTS: %w", err)
	}

	ctsCBOR, err = cts.ToCBOR()
	if err != nil {
		return "", fmt.Errorf("error encoding CoTS to CBOR: %w", err)
	}

	if outputFile == nil || *outputFile == "" {
		ctsFile = makeFileName("", envFile, ".cbor")
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
	cotsCmd.AddCommand(cotsCreateCtsCmd)
}
