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
	"context"

	pb "github.com/jeremyje/coretemp-exporter/proto"
)

type HardwareDataSink interface {
	Observe(ctx context.Context, info *pb.MachineMetrics)
}

type multiSink struct {
	sinks []HardwareDataSink
}

func (m *multiSink) Observe(ctx context.Context, info *pb.MachineMetrics) {
	for _, s := range m.sinks {
		s.Observe(ctx, info)
	}
}
func newMultiSink(s ...HardwareDataSink) *multiSink {
	return &multiSink{
		sinks: s,
	}
}
