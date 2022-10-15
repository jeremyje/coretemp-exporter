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

package coretemp

import (
	"log"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// https://www.alcpu.com/CoreTemp/developers.html

type coreTempSharedDataEx struct {
	// Original structure (CoreTempSharedData)
	uiLoad      [256]uint32
	uiTjMax     [128]uint32
	uiCoreCnt   uint32
	uiCPUCnt    uint32
	fTemp       [256]float32
	fVID        float32
	fCPUSpeed   float32
	fFSBSpeed   float32
	fMultiplier float32
	sCPUName    [100]byte
	// If ucFahrenheit is true, the temperature is reported in Fahrenheit.
	ucFahrenheit byte
	// If ucDeltaToTjMax is true, the temperature reported represents the distance from TjMax.
	ucDeltaToTjMax byte

	// uiStructVersion = 2

	// If ucTdpSupported is true, processor TDP information in the uiTdp array is valid.
	ucTdpSupported byte
	// If ucPowerSupported is true, processor power consumption information in the fPower array is valid.
	ucPowerSupported byte
	uiStructVersion  uint32
	uiTdp            [128]uint32
	fPower           [128]float32
	fMultipliers     [256]float32
}

var (
	fnGetCoreTempInfo *windows.LazyProc
)

func init() {
	coretempDLL := windows.NewLazyDLL("GetCoreTempInfo.dll")
	log.Printf("DLL: %s, %t", coretempDLL.Name, coretempDLL.System)
	fnGetCoreTempInfo = coretempDLL.NewProc("fnGetCoreTempInfoAlt")
}

func GetCoreTempInfo() (*CoreTempInfo, error) {
	data, err := getCoreTempInfoAlt()
	if err != nil {
		return nil, err
	}

	coreCount := int(data.uiCoreCnt)

	if byteToBool(data.ucTdpSupported) {

	}

	return &CoreTempInfo{
		Load:         intList(data.uiLoad[:], coreCount),
		TJMax:        intList(data.uiTjMax[:], int(data.uiCPUCnt)),
		CoreCount:    coreCount,
		Temperature:  float32List(data.fTemp[:], coreCount),
		VID:          data.fVID,
		CPUSpeed:     data.fCPUSpeed,
		FSBSpeed:     data.fFSBSpeed,
		Multiplier:   data.fMultiplier,
		CPUName:      cleanString(string(data.sCPUName[:])),
		Fahrenheit:   byteToBool(data.ucFahrenheit),
		DeltaToTJMax: byteToBool(data.ucDeltaToTjMax),
	}, nil
}

func byteToBool(b byte) bool {
	return b != 0
}

func intList[T uint32 | int32](input []T, size int) []int {
	result := make([]int, size)
	for i := 0; i < int(size); i++ {
		result[i] = int(input[i])
	}
	return result
}

func float32List[T float32 | float64](input []T, size int) []float32 {
	result := make([]float32, size)
	for i := 0; i < int(size); i++ {
		result[i] = float32(input[i])
	}
	return result
}

func cleanString(input string) string {
	return strings.TrimSpace(strings.Trim(input, string("\x00")))
}

func getCoreTempInfoAlt() (*coreTempSharedDataEx, error) {
	if err := fnGetCoreTempInfo.Find(); err != nil {
		return nil, err
	}

	data := &coreTempSharedDataEx{}
	r1, _, err := fnGetCoreTempInfo.Call(uintptr(unsafe.Pointer(data)))

	if r1 != 1 {
		return nil, err
	}
	return data, nil
}
