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

package coretempsdk

import (
	_ "embed"
	"fmt"
	"testing"
)

func ExampleNew() {
	mm, err := New().Get()
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
	fmt.Printf("GetCoreTempInfo: %+v", mm)
}

func TestNew(t *testing.T) {
	driver := New()
	if driver == nil {
		t.Error("driver should not be nil")
	}
}

func TestGet(t *testing.T) {
	mm, err := New().Get()
	if err != nil && mm != nil {
		t.Error("MachineMetrics and error should NOT both be set at the same time.")
	}
	if err == nil && mm == nil {
		t.Error("MachineMetrics and error should NOT both be nil")
	}
}
