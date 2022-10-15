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

// Package coretemp is a Go library for interacting with GetCoreTempInfo.dll.
// You can get the DLL from: https://www.alcpu.com/CoreTemp/main_data/CoreTempSDK.zip
package coretemp

type CoreTempInfo struct {
	Load         []int
	TJMax        []int
	CoreCount    int
	CPUCount     int
	Temperature  []float32
	VID          float32
	CPUSpeed     float32
	FSBSpeed     float32
	Multiplier   float32
	CPUName      string
	Fahrenheit   bool
	DeltaToTJMax bool
}
