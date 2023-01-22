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

package main

import (
	"flag"
	"log"
	"strings"

	"github.com/jeremyje/coretemp-exporter/internal"
)

var (
	oldFilesFlag = flag.String("old", "cputemps.log,cputemps.ndjson", "Comma separated list of old files.")
	newFileFlag  = flag.String("new", "new.ndjson", "The new ndjson file")
)

func main() {
	if err := internal.ConvertV1ToV2(&internal.ConvertArgs{
		OldFiles: strings.Split(*oldFilesFlag, ","),
		NewFile:  *newFileFlag,
	}); err != nil {
		log.Fatal(err)
	}
}
