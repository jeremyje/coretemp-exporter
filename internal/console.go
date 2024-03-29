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
	"fmt"
	"log"

	pb "github.com/jeremyje/coretemp-exporter/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type consoleSink struct {
}

func (s *consoleSink) Observe(ctx context.Context, info *pb.MachineMetrics) {
	log.Println(s.toText(info))
}

func (s *consoleSink) toText(info *pb.MachineMetrics) string {
	txt := ""

	m := &protojson.MarshalOptions{
		Multiline:       false,
		Indent:          "",
		AllowPartial:    false,
		EmitUnpopulated: false,
		UseProtoNames:   false,
	}

	if raw, err := m.Marshal(info); err == nil {
		txt = string(raw)
	} else {
		txt = fmt.Sprintf("%+v", info)
	}
	if txt == "{}" {
		return "<empty>"
	}
	return txt
}
