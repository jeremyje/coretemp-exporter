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

package drivers

import (
	"fmt"
	"log"

	"github.com/yusufpapurcu/wmi"
)

type Sensor struct {
	Identifier string
	Name       string
	Index      int
	InstanceId int
	Max        float32
	Min        float32
	Parent     string
	ProcessId  string
	SensorType string
	Value      float32
}

/*
Identifier       : /lpc/nct6798d/temperature/4
Index            : 4
InstanceId       : 3928
Max              : 42
Min              : 42
Name             : Temperature #4
Parent           : /lpc/nct6798d
ProcessId        : 4f21c82c-c97a-4965-8637-d4e61b3a20f4
SensorType       : Temperature
Value            : 42
PSComputerName   : QUARTZ
*/
func Wmi() string {
	var dst []Sensor
	q := wmi.CreateQuery(&dst, "")
	err := wmi.QueryNamespace(q, &dst, "root\\OpenHardwareMonitor")
	if err != nil {
		log.Fatal(err)
	}

	s := ""
	for i, v := range dst {
		s += fmt.Sprintf("\n%d - %+v", i, v)
	}
	return s
}
