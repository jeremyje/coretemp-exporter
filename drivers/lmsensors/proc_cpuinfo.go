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

package lmsensors

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	cpuinfoFile = "/proc/cpuinfo"
)

type ProcCPUInfo struct {
	VendorId     string  `json:"vendor_id"`
	CPUName      string  `json:"cpu_name"`
	CoreId       string  `json:"core_id"`
	FrequencyMhz float64 `json:"frequency"`
}

func readCPUInfo() (*ProcCPUInfo, error) {
	data, err := os.ReadFile(cpuinfoFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read '%s', err= %w", cpuinfoFile, err)
	}
	return parseCPUInfo(data)
}

func parseCPUInfo(consoleOut []byte) (*ProcCPUInfo, error) {
	all := parseCpuInfoConsoleMaps(consoleOut)

	for _, m := range all {
		freqString := m["cpu MHz"]
		freq, err := strconv.ParseFloat(freqString, 64)
		if err == nil {
			freq *= 1000 * 1000
		}
		return &ProcCPUInfo{
			VendorId:     m["vendor_id"],
			CPUName:      m["model name"],
			FrequencyMhz: freq,
		}, nil
	}

	return nil, fmt.Errorf("cannot get CPU information from '%s'", string(consoleOut))
}

func parseCpuInfoConsoleMaps(consoleOut []byte) []map[string]string {
	all := []map[string]string{}
	m := map[string]string{}

	scanner := bufio.NewScanner(strings.NewReader(string(consoleOut)))
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		parts := strings.SplitN(text, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if strings.ToLower(key) == "processor" {
				if len(m) > 0 {
					all = append(all, m)
				}
				m = map[string]string{}
			}
			m[key] = value
		}
	}

	if len(m) > 0 {
		all = append(all, m)
	}
	return all
}
