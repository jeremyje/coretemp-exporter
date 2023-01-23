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

// Package lmsensors is a Go library reading lm-sensors data.
package lmsensors

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/jeremyje/coretemp-exporter/drivers/common"
	pb "github.com/jeremyje/coretemp-exporter/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func New() common.Driver {
	return newDriver()
}

func fromJSON(out []byte) (*lmsensorData, error) {
	m := map[string]any{}
	if err := json.Unmarshal(sanitizeSensorData(out), &m); err != nil {
		return nil, err
	}
	return &lmsensorData{
		M: m,
	}, nil
}

func sanitizeSensorData(out []byte) []byte {
	data := ""
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(text, "ERROR") {
			data += scanner.Text() + "\n"
		}
	}
	return []byte(data)
}

type lmsensorData struct {
	M map[string]any
}

type LMsensor struct {
	Adapter      string `json:"Adapter"`
	Temperatures map[string]*LMSensorTemperature
}

type LMSensorTemperature struct {
	M map[string]float64
}

func parseLmsensorsOutput(out []byte) (*pb.MachineMetrics, error) {
	data, err := fromJSON(out)
	if err != nil {
		return nil, err
	}

	cpuName := "Unknown CPU"
	frequency := float64(0.0)

	// cpuInfo should not be read outside of this scope.
	{
		cpuInfo, err := readCPUInfo()
		if err == nil {
			cpuName = cpuInfo.CPUName
			frequency = cpuInfo.FrequencyMhz
		}
	}

	temperatures := []float64{}
	load := []int32{}
	for sensorID, sensorDetail := range data.M {
		if strings.Contains(sensorID, "coretemp") {
			concreteSensorDetail, ok := sensorDetail.(map[string]any)
			if ok {
				adapterName := concreteSensorDetail["Adapter"]
				if adapterName == "" {
					adapterName = sensorID
				}
				keys := []string{}
				for detailName := range concreteSensorDetail {
					if strings.Contains(detailName, "Package") {
						continue
					}
					keys = append(keys, detailName)
				}
				sort.Strings(keys)
				for _, detailName := range keys {
					maybeTempDetail := concreteSensorDetail[detailName]
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
					FrequencyMhz: frequency,
				},
			}},
	}, nil
}
