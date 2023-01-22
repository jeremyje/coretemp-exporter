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

package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jeremyje/coretemp-exporter/drivers/common"
	pb "github.com/jeremyje/coretemp-exporter/proto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ConvertArgs struct {
	OldFiles []string
	NewFile  string
}

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

func ConvertV1ToV2(args *ConvertArgs) error {
	os.Remove(args.NewFile)
	fp, err := os.Create(args.NewFile)
	if err != nil {
		return err
	}
	defer fp.Close()

	ln := 0
	for _, filename := range args.OldFiles {
		log.Printf("OPEN %s", filename)
		in, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer in.Close()
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			ln++
			if ln%10000 == 0 {
				log.Printf("Line: %s:%d", filename, ln)
			}
			pba := &pb.MachineMetrics{}
			line := scanner.Bytes()
			if len(line) < 10 {
				continue
			}
			err := json.Unmarshal(line, &pba)
			if err != nil || pba == nil || len(pba.GetName()) == 0 {
				info := &HardwareInfo{}
				err := json.Unmarshal(line, &info)
				if err != nil {
					return err
				}
				pba = convertToProto(info)
				if len(pba.Name) == 0 {
					if err != nil {
						return fmt.Errorf("empty name: %+v", pba)
					}
				}
			}

			newLine, err := protojson.Marshal(pba)
			if err != nil {
				return err
			}
			_, err = fp.Write(newLine)
			if err != nil {
				return err
			}
			_, err = fp.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func copyInt(i []int) []int32 {
	a := []int32{}
	for _, v := range i {
		a = append(a, int32(v))
	}
	return a
}

func convertToProto(info *HardwareInfo) *pb.MachineMetrics {
	return &pb.MachineMetrics{
		Name: "quartz",
		Device: []*pb.DeviceMetrics{
			{
				Name:        info.CPUName,
				Kind:        "cpu",
				Temperature: common.Average(info.TemperatureCelcius),
				Cpu: &pb.CpuDeviceMetrics{
					Load:            copyInt(info.Load),
					Temperature:     info.TemperatureCelcius,
					NumCores:        int32(info.CoreCount),
					FrequencyMhz:    info.CPUSpeed,
					FsbFrequencyMhz: info.FSBSpeed,
				},
			},
		},
		Timestamp: timestamppb.New(info.Timestamp),
	}
}
