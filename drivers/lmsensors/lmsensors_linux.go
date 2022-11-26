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

	"github.com/jeremyje/coretemp-exporter/drivers/common"
	pb "github.com/jeremyje/coretemp-exporter/proto"
)

type lmsensorsDriver struct {
}

func (d *lmsensorsDriver) Get() (*pb.MachineMetrics, error) {
	out, err := exec.Command("sensors", "-j").Output()
	if err != nil {
		return nil, fmt.Errorf("cannot run 'sensors' command, is it installed or running in a VM?\nout= %s\nerr= %w", out, err)
	}
	return parseLmsensorsOutput(out)
}

func newDriver() common.Driver {
	return &lmsensorsDriver{}
}
