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

// Package common contains all core structures and API for CoreTemp.
package common

import "time"

// HardwareInfo includes CPU and motherboard information about the host machine.
type HardwareInfo struct {
	// Load as a percentage [0-100].
	Load []int `json:"load" yaml:"load"`
	// TJMax is the thermal junction maximum of the CPU. Once the temperature hits this limit the CPU will thermal throttle.
	TJMax []float64 `json:"tjMax" yaml:"tjMax"`
	// CoreCount is the total number of cores across all CPUs on the machine.
	CoreCount int `json:"coreCount" yaml:"coreCount"`
	// CPUCount is the number of CPUs on the machine.
	CPUCount int `json:"cpuCount" yaml:"cpuCount"`
	// TemperatureCelcius is the temperature of each core expressed in Celcius.
	TemperatureCelcius []float64 `json:"temperatureCelcius" yaml:"temperatureCelcius"`
	// VID is the voltage requested by the CPU.
	VID float64 `json:"vid" yaml:"vid"`
	// CPUSpeed is the clock frequency of the CPU.
	CPUSpeed float64 `json:"cpuSpeed" yaml:"cpuSpeed"`
	// FSBSpeed is the clock frequency of the front side bus.
	FSBSpeed float64 `json:"fsbSpeed" yaml:"fsbSpeed"`
	// Multiplier is the front side bus multipler of the CPU.
	Multiplier float64 `json:"multiplier" yaml:"multiplier"`
	// CPUName is the name and model of the CPU.
	CPUName string `json:"cpuName" yaml:"cpuName"`
	// Timestamp is the time that the sample was taken.
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
}

type Driver interface {
	Get() (*HardwareInfo, error)
}

func NotSupported(err error) Driver {
	return &notSupportedDriver{
		err: err,
	}
}

type notSupportedDriver struct {
	err error
}

func (n *notSupportedDriver) Get() (*HardwareInfo, error) {
	return nil, n.err
}
