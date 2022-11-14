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

package lmsensors

import (
	_ "embed"
	"fmt"
	"testing"
)

var (
	//go:embed testdata/example.json
	exampleJSON []byte
)

func ExampleNew() {
	info, err := New().Get()
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
	fmt.Printf("GetCoreTempInfo: %+v", info)
}

func TestGet(t *testing.T) {
	data, err := fromJSON(exampleJSON)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.M) == 0 {
		t.Error("result was empty")
	}
}
