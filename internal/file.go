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
	"context"
	"log"
	"os"

	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/jeremyje/coretemp-exporter/proto"
)

var (
	newlineAsByte []byte = []byte("\n")
)

type fileSink struct {
	fp *os.File
}

func newFileSink(name string) (*fileSink, error) {
	fp, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &fileSink{
		fp: fp,
	}, nil
}

func (s *fileSink) Observe(ctx context.Context, info *pb.MachineMetrics) {
	line, err := protojson.Marshal(info)
	if err == nil {
		if _, err := s.fp.Write(line); err != nil {
			log.Printf("ERROR: %s", err)
		}
		if _, err := s.fp.Write(newlineAsByte); err != nil {
			log.Printf("ERROR: %s", err)
		}
	}
}
