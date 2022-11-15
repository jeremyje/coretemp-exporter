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

import (
	"os"

	pb "github.com/jeremyje/coretemp-exporter/proto"
)

type Driver interface {
	Get() (*pb.MachineMetrics, error)
}

func NotSupported(err error) Driver {
	return &notSupportedDriver{
		err: err,
	}
}

type notSupportedDriver struct {
	err error
}

func (n *notSupportedDriver) Get() (*pb.MachineMetrics, error) {
	return nil, n.err
}

func Hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return name
}

func Average(val []float64) float64 {
	size := len(val)
	if size == 0 {
		return float64(0.0)
	}
	total := float64(0.0)
	for _, v := range val {
		total += v
	}
	return total / float64(size)
}
