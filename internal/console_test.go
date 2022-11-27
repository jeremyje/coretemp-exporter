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

package internal

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	pb "github.com/jeremyje/coretemp-exporter/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestConsoleToText(t *testing.T) {
	wellKnownTimestamp, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		input *pb.MachineMetrics
		want  string
	}{
		{
			input: nil,
			want:  "<empty>",
		},
		{
			input: &pb.MachineMetrics{},
			want:  "<empty>",
		},
		{
			input: &pb.MachineMetrics{
				Name: "machine-name",
				Device: []*pb.DeviceMetrics{{
					Name:        "some-processor",
					Kind:        "cpu",
					Temperature: 45,
					Cpu: &pb.CpuDeviceMetrics{
						Load:            []int32{1, 2, 3, 4},
						Temperature:     []float64{50, 40, 50, 40},
						NumCores:        4,
						FrequencyMhz:    1000,
						FsbFrequencyMhz: 100,
					},
				}},
				Timestamp: timestamppb.New(wellKnownTimestamp),
			},
			want: `{"name":"machine-name","device":[{"name":"some-processor","kind":"cpu","temperature":45,"cpu":{"load":[1,2,3,4],"temperature":[50,40,50,40],"numCores":4,"frequencyMhz":1000,"fsbFrequencyMhz":100}}],"timestamp":"2006-01-02T15:04:05Z"}`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%+v", tc.input), func(t *testing.T) {
			t.Parallel()

			c := &consoleSink{}
			// https://stackoverflow.com/questions/72359452/proto-unmarshal-test-fails-inconsistently
			got := strings.ReplaceAll(c.toText(tc.input), " ", "")
			want := strings.ReplaceAll(tc.want, " ", "")
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("\n\nWANT:\n\n%s\n\n", tc.want)
				t.Errorf("\n\nGOT:\n\n%s\n\n", got)
				t.Errorf("consoleSink.toText() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
