// Copyright 2023 Jeremy Edwards
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package common contains all core structures and API for CoreTemp.
package common

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func DownloadFile(url string) ([]byte, error) {
	buf := &bytes.Buffer{}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func UnzipFile(buf []byte, pattern string, outFile string) error {
	bufReader := bytes.NewReader(buf)
	z, err := zip.NewReader(bufReader, bufReader.Size())
	if err != nil {
		return err
	}
	for _, entry := range z.File {
		if strings.Contains(entry.Name, pattern) {
			r, err := entry.Open()
			if err != nil {
				return fmt.Errorf("cannot open file '%s' within zip file. err= %w", entry.Name, err)
			}
			fp, err := os.Create(outFile)
			if err != nil {
				return fmt.Errorf("cannot create file '%s', err= %w", outFile, err)
			}
			if _, err := io.Copy(fp, r); err != nil {
				return fmt.Errorf("cannot file contents to file '%s', err= %w", outFile, err)
			}
			return nil
		}
	}
	return fmt.Errorf("cannot find a file with name '%s' in zip file", pattern)
}
