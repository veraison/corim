// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
)

func filesList(files, dirs []string, ext string) []string {
	var l []string

	for _, file := range files {
		if _, err := fs.Stat(file); err == nil {
			l = append(l, file)
		}
	}

	for _, dir := range dirs {
		filesInfo, err := afero.ReadDir(fs, dir)
		if err != nil {
			continue
		}

		for _, fileInfo := range filesInfo {
			if !fileInfo.IsDir() && filepath.Ext(fileInfo.Name()) == ext {
				l = append(l, filepath.Join(dir, fileInfo.Name()))
			}
		}
	}

	return l
}

type FromCBORLoader interface {
	FromCBOR([]byte) error
}

func printJSONFromCBOR(fcl FromCBORLoader, cbor []byte, heading string) error {
	var (
		err error
		j   []byte
	)

	if err = fcl.FromCBOR(cbor); err != nil {
		return fmt.Errorf("CBOR decoding failed: %w", err)
	}

	indent := "  "
	if j, err = json.MarshalIndent(fcl, "", indent); err != nil {
		return fmt.Errorf("JSON encoding failed: %w", err)
	}

	fmt.Println(heading)
	fmt.Println(string(j))

	return nil
}

func printComid(cbor []byte, heading string) error {
	return printJSONFromCBOR(&comid.Comid{}, cbor, heading)
}

func printCoswid(cbor []byte, heading string) error {
	return printJSONFromCBOR(&swid.SoftwareIdentity{}, cbor, heading)
}

func makeFileName(dirName, baseName, ext string) string {
	return filepath.Join(
		dirName,
		filepath.Base(
			strings.TrimSuffix(
				baseName,
				filepath.Ext(baseName),
			),
		)+ext,
	)
}
