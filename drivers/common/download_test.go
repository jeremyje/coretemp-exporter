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
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadFromZip(t *testing.T) {
	zipFile, err := DownloadFile("https://www.alcpu.com/CoreTemp/main_data/CoreTempSDK.zip")
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	dllPath := filepath.Join(dir, "GetCoreTempInfo.dll")
	if err := UnzipFile(zipFile, "GetCoreTempInfo.dll", dllPath); err != nil {
		t.Fatal(err)
	}

	stat, err := os.Stat(dllPath)
	if err != nil {
		t.Fatal(err)
	}
	if stat.Size() < 5000 {
		t.Errorf("GetCoreTempInfo.dll is not at least 5KB, size= %d", stat.Size())
	}
}
