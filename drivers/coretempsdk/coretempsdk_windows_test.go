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

//go:build windows
// +build windows

package coretempsdk

import "testing"

func TestGetCoreTempInfo(t *testing.T) {
	info, err := getCoreTempInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", info)
}

func TestGetCoreTempInfoAlt(t *testing.T) {
	info, err := getCoreTempInfoAlt()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", info)
}
