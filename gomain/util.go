// Copyright 2022 Jeremy Edwards
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

//go:build windows
// +build windows

package gomain

import (
	"fmt"
	"os"
	"path/filepath"
)

func exePath() (string, error) {
	return exePathFromPath(os.Args[0])
}

func exePathFromPath(prog string) (string, error) {
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", fmt.Errorf("cannot get the absolute path of '%s', err= %w", prog, err)
	}
	if fi, statErr := os.Stat(p); statErr == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"

		if fi, statErr := os.Stat(p); statErr == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is directory", p)
		}
	}
	return "", err

}
