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

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"unsafe"

	"github.com/jeremyje/coretemp-exporter/drivers/common"
	pb "github.com/jeremyje/coretemp-exporter/proto"
	"golang.org/x/sys/windows"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	globalFnGetCoreTempInfo *windows.LazyProc
	globalLock              sync.RWMutex
)

const (
	dllNameGetCoreTempInfoDLL            = "GetCoreTempInfo.dll"
	dllFuncfnGetCoreTempInfoAlt          = "fnGetCoreTempInfoAlt"
	dllURIGetCoreTempInfoDLLDownloadPage = "https://www.alcpu.com/CoreTemp/developers.html"
	dllURIGetCoreTempInfoDLL             = "https://www.alcpu.com/CoreTemp/main_data/CoreTempSDK.zip"
	vcRuntimeDownloadPage                = "https://learn.microsoft.com/en-US/cpp/windows/latest-supported-vc-redist"

	enableAutoDownload = false
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

func getCoreTempInfo() (*pb.MachineMetrics, error) {
	rawInfo, err := getCoreTempInfoAlt()
	if err != nil {
		return nil, err
	}

	coreCount := int(rawInfo.uiCoreCnt)

	temps := float64List(rawInfo.fTemp[:], coreCount)

	if byteToBool(rawInfo.ucFahrenheit) {
		for i := 0; i < len(temps); i++ {
			temps[i] = fToC(temps[i])
		}
	}

	return &pb.MachineMetrics{
		Name:      common.Hostname(),
		Timestamp: timestamppb.Now(),
		Device: []*pb.DeviceMetrics{
			{
				Name:        cleanString(string(rawInfo.sCPUName[:])),
				Kind:        "cpu",
				Temperature: common.Average(temps),
				Cpu: &pb.CpuDeviceMetrics{
					Load:            int32List(rawInfo.uiLoad[:], coreCount),
					Temperature:     temps,
					NumCores:        int32(coreCount),
					FrequencyMhz:    float64(rawInfo.fCPUSpeed),
					FsbFrequencyMhz: float64(rawInfo.fFSBSpeed),
				},
			},
		},
	}, nil
}

func fToC(f float64) float64 {
	return (f - 32) * 5 / 9
}

func byteToBool(b byte) bool {
	return b != 0
}

func int32List[T uint32 | int32](input []T, size int) []int32 {
	result := make([]int32, size)
	for i := 0; i < int(size); i++ {
		result[i] = int32(input[i])
	}
	return result
}

func float64List[T float32 | float64 | int | uint32](input []T, size int) []float64 {
	result := make([]float64, size)
	for i := 0; i < int(size); i++ {
		result[i] = float64(input[i])
	}
	return result
}

func cleanString(input string) string {
	return strings.TrimSpace(strings.Trim(input, string("\x00")))
}

type coreTempSDKError struct {
	msg string
	err error
}

func (d coreTempSDKError) Error() string {
	return d.msg
}

func wrapDLLError(err error) error {
	dir, errWD := os.Getwd()
	if errWD != nil {
		dir = "."
	}
	return coreTempSDKError{
		msg: fmt.Sprintf("Make sure that '%s' is in directory '%s'. And the version is at least 1.2.0.0. You can download the DLL from '%s'. You also need VC++ Runtime 10 from '%s'. Error= %s", dllNameGetCoreTempInfoDLL, dir, dllURIGetCoreTempInfoDLLDownloadPage, vcRuntimeDownloadPage, err.Error()),
		err: err,
	}
}

func wrapCallError(err error) error {
	return coreTempSDKError{
		msg: fmt.Sprintf("The DLL was found but the call failed. Make sure Core Temp is running in the background. Error= %s", err.Error()),
		err: err,
	}
}

func getFnGetCoreTempInfo() (*windows.LazyProc, error) {
	globalLock.RLock()
	dllLoaded := globalFnGetCoreTempInfo == nil
	globalLock.RUnlock()

	if !dllLoaded {
		ensureCoreTempDLL()
	}

	globalLock.Lock()
	defer globalLock.Unlock()

	if globalFnGetCoreTempInfo == nil {
		coretempDLL := windows.NewLazyDLL(dllNameGetCoreTempInfoDLL)
		globalFnGetCoreTempInfo = coretempDLL.NewProc(dllFuncfnGetCoreTempInfoAlt)
	}
	if err := globalFnGetCoreTempInfo.Find(); err != nil {
		globalFnGetCoreTempInfo = nil
		return nil, wrapDLLError(err)
	}

	return globalFnGetCoreTempInfo, nil
}

func getCoreTempInfoAlt() (*coreTempSharedDataEx, error) {
	fnGetCoreTempInfo, err := getFnGetCoreTempInfo()
	if err != nil {
		return nil, wrapCallError(err)
	}

	rawInfo := &coreTempSharedDataEx{}
	r1, _, err := fnGetCoreTempInfo.Call(uintptr(unsafe.Pointer(rawInfo)))

	if r1 != 1 {
		return nil, err
	}
	return rawInfo, nil
}

func ensureCoreTempDLL() error {
	stat, err := os.Stat(dllNameGetCoreTempInfoDLL)
	if err == nil {
		if stat.Size() > 5000 {
			return nil
		}
	}
	if !enableAutoDownload {
		return nil
	}
	zipData, err := common.DownloadFile(dllURIGetCoreTempInfoDLL)
	if err != nil {
		return fmt.Errorf("cannot download '%s', err= %w", dllURIGetCoreTempInfoDLL, err)
	}
	globalLock.Lock()
	defer globalLock.Unlock()
	if err := common.UnzipFile(zipData, dllNameGetCoreTempInfoDLL, dllNameGetCoreTempInfoDLL); err != nil {
		return fmt.Errorf("cannot extract '%s' from zip file '%s', err= %w", dllNameGetCoreTempInfoDLL, dllURIGetCoreTempInfoDLL, err)
	}
	return nil
}

type coreTempSDKDriver struct {
}

func (d *coreTempSDKDriver) Get() (*pb.MachineMetrics, error) {
	return getCoreTempInfo()
}

func newDriver() common.Driver {
	return &coreTempSDKDriver{}
}
