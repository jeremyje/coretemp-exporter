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

package common

import (
	_ "embed"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed testdata/example.json
	exampleJSON []byte
	//go:embed testdata/example.yaml
	exampleYAML []byte
)

func TestMarshalJSON(t *testing.T) {
	data := &HardwareInfo{}
	if err := json.Unmarshal(exampleJSON, data); err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff([]int{0, 1, 2, 3}, data.Load); diff != "" {
		t.Errorf("json.Unmarshal() mismatch (-want +got):\n%s", diff)
	}
	m, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(exampleJSON, m); diff != "" {
		t.Logf("json.Marshal\n------------\n%s", string(m))
		t.Errorf("json.Marshal() mismatch (-want +got):\n%s", diff)
	}
}

func TestMarshalYAML(t *testing.T) {
	want := []byte(strings.ReplaceAll(string(exampleYAML), "\r", ""))
	data := &HardwareInfo{}
	if err := yaml.Unmarshal(want, data); err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff([]int{0, 1, 2, 3}, data.Load); diff != "" {
		t.Errorf("yaml.Unmarshal() mismatch (-want +got):\n%s", diff)
	}
	m, err := yaml.Marshal(data)
	if err != nil {
		t.Error(err)
	}
	m = []byte(strings.ReplaceAll(string(m), "\r", ""))
	if diff := cmp.Diff(want, m); diff != "" {
		t.Logf("yaml.Marshal\n------------\n%s", string(m))
		t.Errorf("yaml.Marshal() mismatch (-want +got):\n%s", diff)
	}
}
