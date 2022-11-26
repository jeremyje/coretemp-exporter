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

package gomain

import (
	"testing"
)

func TestHandleSignalBase(t *testing.T) {
	for _, tc := range handleSignalTestCases {
		tc := tc
		t.Run(tc.input.String(), func(t *testing.T) {
			t.Parallel()
			got := handleSignal(tc.input)
			if got != tc.want {
				t.Fatalf("expected: %t, got: %t", tc.want, got)
			}
		})
	}
}
