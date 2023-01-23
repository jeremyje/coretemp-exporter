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
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	pb "github.com/jeremyje/coretemp-exporter/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type ConvertCSVArgs struct {
	InputFile  []string
	OutputFile string
}

func ConvertCSV(args *ConvertCSVArgs) error {
	os.Remove(args.OutputFile)
	fp, err := os.Create(args.OutputFile)
	if err != nil {
		return fmt.Errorf("cannot create output file '%s', %w", args.OutputFile, err)
	}

	cw := csv.NewWriter(fp)
	cw.Write([]string{
		"ts_year", "ts_month", "ts_day", "ts_hour", "ts_min", "ts_second", // Timestamp
		"active_cores", "total_cores", "frequency", "avg_load", // Basic CPU metrics
		"temperature"}) // CPU Temperature
	defer cw.Flush()
	return scanNdJson(args.InputFile, func(item *pb.MachineMetrics) error {
		if len(item.GetDevice()) == 0 {
			return nil
		}
		t := item.GetTimestamp().AsTime()

		cpuMetrics := item.GetDevice()[0]
		activeCores := 0
		totalLoad := 0
		for _, coreLoad := range cpuMetrics.Cpu.GetLoad() {
			if coreLoad > 0 {
				activeCores++
			}
			totalLoad += int(coreLoad)
		}
		numCores := len(cpuMetrics.Cpu.GetLoad())
		avgLoad := float64(totalLoad) / float64(numCores)

		return cw.Write([]string{
			strconv.Itoa(t.Year()), strconv.Itoa(int(t.Month())), strconv.Itoa(t.Day()), strconv.Itoa(t.Hour()), strconv.Itoa(t.Minute()), strconv.Itoa(t.Second()),
			strconv.Itoa(activeCores), strconv.Itoa(numCores), fmt.Sprintf("%f", cpuMetrics.GetCpu().GetFrequencyMhz()), fmt.Sprintf("%f", avgLoad),
			fmt.Sprintf("%f", cpuMetrics.GetTemperature()),
		})
	})
}

func scanNdJson(inputFiles []string, consumerFunc func(line *pb.MachineMetrics) error) error {
	for _, inputFile := range inputFiles {
		fp, err := os.Open(inputFile)
		if err != nil {
			return fmt.Errorf("cannot open '%s', %w", inputFile, err)
		}
		defer fp.Close()

		ln := 0
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			ln++
			if ln%10000 == 0 {
				log.Printf("Line: %s:%d", inputFile, ln)
			}
			line := scanner.Bytes()
			if len(line) < 10 {
				continue
			}
			mm := &pb.MachineMetrics{}
			err := protojson.Unmarshal(line, mm)
			if err != nil || mm == nil || len(mm.GetName()) == 0 {
				return fmt.Errorf("cannot read line '%s:%d', %w", inputFile, ln, err)
			}
			if err := consumerFunc(mm); err != nil {
				return err
			}
		}
	}
	return nil
}
