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

package testing

import (
	"testing"
)

func TestTestMainAsync(t *testing.T) {
	ready := make(chan int)
	m := func(waitFunc func()) error {
		ready <- 1
		waitFunc()
		return nil
	}
	close := Main(m)
	val := <-ready
	if val != 1 {
		t.Errorf("got: %d, want: 1", val)
	}
	if err := close(); err != nil {
		t.Error(err)
	}
}
