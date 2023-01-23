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

package common

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	pb "github.com/jeremyje/coretemp-exporter/proto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed testdata/example.json
	exampleJSON []byte
	//go:embed testdata/example.yaml
	exampleYAML []byte
)

func TestBasic(t *testing.T) {
	wantJSON := cleanJSON(exampleJSON)
	wantYAML := noLF(exampleYAML)
	ts, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}
	want := &pb.MachineMetrics{
		Name: "test",
		Device: []*pb.DeviceMetrics{
			{
				Name:        "vermeer",
				Kind:        "cpu",
				Temperature: 37.325,
				Cpu: &pb.CpuDeviceMetrics{
					Load:            []int32{0, 1, 2, 3},
					Temperature:     []float64{1, 4, 65.4, 78.9},
					NumCores:        4,
					FrequencyMhz:    5000.2,
					FsbFrequencyMhz: 100.4,
				},
			},
		},
		Timestamp: timestamppb.New(ts),
	}

	gotProto := &pb.MachineMetrics{}
	if err := protojson.Unmarshal(wantJSON, gotProto); err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(want, gotProto, protocmp.Transform()); diff != "" {
		t.Errorf("protojson.Unmarshal() mismatch (-want +got):\n%s", diff)
	}
	if gotJSON, err := protojson.Marshal(want); err != nil {
		t.Error(err)
	} else {
		gotJSON = cleanJSON(gotJSON)
		if diff := cmp.Diff(wantJSON, gotJSON); diff != "" {
			t.Logf("protojson.Marshal\n============\n%s", string(gotJSON))
			t.Errorf("protojson.Marshal() mismatch (-want +got):\n%s", diff)
		}
	}
	if gotYAML, err := yaml.Marshal(want); err != nil {
		t.Error(err)
	} else {
		gotYAML = noLF(gotYAML)
		if diff := cmp.Diff(wantYAML, gotYAML); diff != "" {
			t.Logf("yaml.Marshal\n============\n%s", string(gotYAML))
			t.Errorf("yaml.Marshal() mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	data := &pb.MachineMetrics{}
	want := cleanJSON(exampleJSON)
	if err := protojson.Unmarshal(want, data); err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff([]int32{0, 1, 2, 3}, data.GetDevice()[0].GetCpu().GetLoad()); diff != "" {
		t.Errorf("protojson.Unmarshal() mismatch (-want +got):\n%s", diff)
	}
	m, err := protojson.Marshal(data)
	if err != nil {
		t.Error(err)
	}

	m = cleanJSON(m)
	if diff := cmp.Diff(want, m); diff != "" {
		t.Logf("protojson.Marshal\n============\n%s", string(m))
		t.Errorf("protojson.Marshal() mismatch (-want +got):\n%s", diff)
	}
}

func TestMarshalYAML(t *testing.T) {
	want := noLF(exampleYAML)
	data := &pb.MachineMetrics{}
	if err := yaml.Unmarshal(want, data); err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff([]int32{0, 1, 2, 3}, data.GetDevice()[0].GetCpu().GetLoad()); diff != "" {
		t.Errorf("yaml.Unmarshal() mismatch (-want +got):\n%s", diff)
	}
	m, err := yaml.Marshal(data)
	if err != nil {
		t.Error(err)
	}
	m = noLF(m)
	if diff := cmp.Diff(want, m); diff != "" {
		t.Logf("yaml.Marshal() GOT\n============\n%s", string(m))
		t.Logf("yaml.Marshal() WANT\n============\n%s", string(want))
		t.Errorf("yaml.Marshal() mismatch (-want +got):\n%s", diff)
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		input []float64
		want  float64
	}{
		{
			input: []float64{},
			want:  0.0,
		},
		{
			input: []float64{1.0},
			want:  1.0,
		},
		{
			input: []float64{1.0, 3.0},
			want:  2.0,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%v => %f", tc.input, tc.want), func(t *testing.T) {
			t.Parallel()

			got := Average(tc.input)
			if tc.want != got {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func noLF(data []byte) []byte {
	return []byte(strings.ReplaceAll(string(data), "\r", ""))
}

func cleanJSON(data []byte) []byte {
	return []byte(strings.ReplaceAll(string(noLF(data)), " ", ""))
}
