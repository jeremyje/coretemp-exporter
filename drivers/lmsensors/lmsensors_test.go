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

package lmsensors

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	pb "github.com/jeremyje/coretemp-exporter/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	//go:embed testdata/sensors.json
	sensorsJSON []byte
	//go:embed testdata/sensors_nuc2.json
	sensorsNuc2JSON []byte
)

func ExampleNew() {
	info, err := New().Get()
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
	fmt.Printf("GetCoreTempInfo: %+v", info)
}

func TestGet(t *testing.T) {
	data, err := fromJSON(sensorsJSON)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.M) == 0 {
		t.Error("result was empty")
	}
}

func TestParseLmsensorsOutput(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  *pb.MachineMetrics
	}{
		{
			name:  "sensors.json",
			input: sensorsJSON,
			want: &pb.MachineMetrics{
				Device: []*pb.DeviceMetrics{
					{
						Name:        "",
						Kind:        "cpu",
						Temperature: 36.5,
						Cpu: &pb.CpuDeviceMetrics{
							NumCores:    2,
							Temperature: []float64{30, 43},
						},
					},
				},
			},
		},
		{
			name:  "sensors_nuc2.json",
			input: sensorsNuc2JSON,
			want: &pb.MachineMetrics{
				Device: []*pb.DeviceMetrics{
					{
						Name:        "",
						Kind:        "cpu",
						Temperature: 45,
						Cpu: &pb.CpuDeviceMetrics{
							NumCores:    2,
							Temperature: []float64{45, 45},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseLmsensorsOutput(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform(), protocmp.IgnoreFields(&pb.MachineMetrics{}, "name", "timestamp"), protocmp.IgnoreFields(&pb.DeviceMetrics{}, "name"), protocmp.IgnoreFields(&pb.CpuDeviceMetrics{}, "frequency_mhz")); diff != "" {
				t.Errorf("protojson.Unmarshal() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
