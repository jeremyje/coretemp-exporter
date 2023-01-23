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

package internal

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	//go:embed testdata/cputemps.ndjson
	cputempsNdjson []byte
	//go:embed testdata/cputemps.csv
	cputempsCsv []byte
)

func TestConvertCSV(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.ndjson")
	output := filepath.Join(dir, "output.csv")
	if err := os.WriteFile(input, cputempsNdjson, 0664); err != nil {
		t.Fatalf("cannot write '%s', %s", input, err)
	}

	if err := ConvertCSV(&ConvertCSVArgs{
		InputFile:  []string{input},
		OutputFile: output,
	}); err != nil {
		t.Fatalf("cannot write csv '%s', %s", output, err)
	}

	actual, err := os.ReadFile(output)
	if err != nil {
		t.Errorf("cannot read back '%s', %s", output, err)
	}
	if diff := cmp.Diff(string(actual), string(cputempsCsv)); diff != "" {
		t.Errorf(diff)
	}
}
