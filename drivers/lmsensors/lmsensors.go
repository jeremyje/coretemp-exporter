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

// Package lmsensors is a Go library reading lm-sensors data.
package lmsensors

import (
	"bufio"
	"encoding/json"
	"strings"

	"github.com/jeremyje/coretemp-exporter/drivers/common"
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
