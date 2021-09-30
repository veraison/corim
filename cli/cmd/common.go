// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"path/filepath"

	"github.com/spf13/afero"
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
