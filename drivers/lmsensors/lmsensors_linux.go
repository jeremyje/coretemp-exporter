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
	"os/exec"
	"strconv"
	"strings"

	"github.com/jeremyje/coretemp-exporter/drivers/common"
)

type lmsensorsDriver struct {
}

func (d *lmsensorsDriver) Get() (*common.HardwareInfo, error) {
	out, err := exec.Command("sensors", "-j").Output()
	if err != nil {
		return nil, fmt.Errorf("cannot run 'sensors' command, is it installed?\nout= %s\nerr= %w", out, err)
	}
	data, err := fromJSON(out)
	if err != nil {
		return nil, err
	}
	r := &common.HardwareInfo{}
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
								r.TemperatureCelcius = append(r.TemperatureCelcius, s)
							}
						}
					}
				}
			}
		}
	}

	return r, nil
}

func newDriver() common.Driver {
	return &lmsensorsDriver{}
}
