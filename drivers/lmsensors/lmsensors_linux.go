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

//go:build linux
// +build linux

package lmsensors

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jeremyje/coretemp-exporter/drivers/common"
	pb "github.com/jeremyje/coretemp-exporter/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type lmsensorsDriver struct {
}

func (d *lmsensorsDriver) Get() (*pb.MachineMetrics, error) {
	out, err := exec.Command("sensors", "-j").Output()
	if err != nil {
		return nil, fmt.Errorf("cannot run 'sensors' command, is it installed?\nout= %s\nerr= %w", out, err)
	}
	data, err := fromJSON(out)
	if err != nil {
		return nil, err
	}
	cpuName := "Unknown CPU"
	cpuInfo, err := readCPUInfo()
	if err == nil {
		cpuName = cpuInfo.CPUName
	}

	temperatures := []float64{}
	load := []int32{}
	for sensorID, sensorDetail := range data.M {
		concreteSensorDetail, ok := sensorDetail.(map[string]any)
		if ok {
			adapterName := concreteSensorDetail["Adapter"]
			if adapterName == "" {
				adapterName = sensorID
			}

			for _, maybeTempDetail := range concreteSensorDetail {
				if concreteTempDetail, ok := maybeTempDetail.(map[string]any); ok {
					for name, value := range concreteTempDetail {
						if strings.Contains(name, "input") {
							if s, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64); err == nil {
								temperatures = append(temperatures, s)
							}
						}
					}
				}
			}
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return &pb.MachineMetrics{
		Name:      hostname,
		Timestamp: timestamppb.Now(),
		Device: []*pb.DeviceMetrics{
			{
				Name:        cpuName,
				Kind:        "cpu",
				Temperature: common.Average(temperatures),
				Cpu: &pb.CpuDeviceMetrics{
					Load:         load,
					Temperature:  temperatures,
					NumCores:     int32(len(temperatures)),
					FrequencyMhz: cpuInfo.FrequencyMhz,
				},
			}},
	}, nil
}

func newDriver() common.Driver {
	return &lmsensorsDriver{}
}
