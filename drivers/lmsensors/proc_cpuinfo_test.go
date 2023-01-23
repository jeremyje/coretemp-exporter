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

package lmsensors

import (
	_ "embed"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	//go:embed testdata/proc_cpuinfo.txt
	cpuInfoTXT []byte
)

func TestParseCPUInfo(t *testing.T) {
	data, err := parseCPUInfo(cpuInfoTXT)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff("Intel(R) Celeron(R) CPU  N3050  @ 1.60GHz", data.CPUName); diff != "" {
		t.Errorf(".CPUName mismatch (-want +got):\n%s", diff)
	}
}
