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

package gomain

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExePath(t *testing.T) {
	ep := exePath()
	if ep == "" {
		t.Error("exePath() should not be empty")
	}
}

func TestExePathFromPath(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		input string
		want  string
	}{
		{input: "util_test", want: filepath.Join(dir, "util_test")},
		{input: "util_test.exe", want: filepath.Join(dir, "util_test.exe")},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			got := exePathFromPath(tc.input)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
