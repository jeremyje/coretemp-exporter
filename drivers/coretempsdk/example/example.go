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

// Package main is an example for using the Core Temp SDK.
package main

import (
	"encoding/json"
	"log"

	"github.com/jeremyje/coretemp-exporter/drivers/coretempsdk"
)

func main() {
	d := coretempsdk.New()
	info, err := d.Get()
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}
	log.Printf("Hostname: %s", info.Name)
	log.Printf("CPU: %s", info.GetDevice()[0].Name)
	log.Printf("Temperatures: %v", info.GetDevice()[0].GetTemperature())
	data, err := json.Marshal(info)
	if err != nil {
		log.Printf("%+v", info)
	} else {
		log.Printf("%s", data)
	}
}
