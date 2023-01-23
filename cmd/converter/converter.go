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
	"fmt"
	"log"
	"strings"

	"github.com/jeremyje/coretemp-exporter/internal"
)

var (
	inputFilesFlag = flag.String("input", "cputemps.log,cputemps.ndjson", "Comma separated list of old files.")
	outputFileFlag = flag.String("output", "new.ndjson", "The new ndjson file")
	modeFlag       = flag.String("mode", "csv", "Conversion mode (update, csv)")
)

func main() {
	flag.Parse()
	var err error
	switch *modeFlag {
	case "csv":
		err = internal.ConvertCSV(&internal.ConvertCSVArgs{
			InputFile:  strings.Split(*inputFilesFlag, ","),
			OutputFile: *outputFileFlag,
		})
	case "update":
		err = internal.ConvertV1ToV2(&internal.ConvertArgs{
			InputFiles: strings.Split(*inputFilesFlag, ","),
			OutputFile: *outputFileFlag,
		})
	default:
		err = fmt.Errorf("mode '%s' is not supported", *modeFlag)
	}
	if err != nil {
		log.Fatal(err)
	}
}
