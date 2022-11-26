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
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	waitForeverFuncs = map[string]MainFunc{
		"waitForeverFunc": func(wait func()) error {
			wait()
			return nil
		},
		"waitForeverFailingFunc": func(wait func()) error {
			wait()
			return fmt.Errorf("failed")
		},
	}

	immediateMainFuncs = map[string]MainFunc{
		"immediateReturnFunc": func(wait func()) error {
			return nil
		},
		"immediateFailFunc": func(wait func()) error {
			return fmt.Errorf("failed")
		},
	}
)

func getAllMainFuncs() map[string]MainFunc {
	mains := map[string]MainFunc{}
	for k, v := range waitForeverFuncs {
		mains[k] = v
	}
	for k, v := range immediateMainFuncs {
		mains[k] = v
	}
	return mains
}

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

func TestRunInteractiveInternal(t *testing.T) {
	for mainName, mainFunc := range getAllMainFuncs() {
		mainFunc := mainFunc
		for _, tc := range handleSignalTestCases {
			tc := tc

			t.Run(fmt.Sprintf("%s - %s", mainName, tc.input.String()), func(t *testing.T) {
				t.Parallel()
				sigCh := make(chan os.Signal, 1)

				go func() {
					time.Sleep(time.Millisecond * 100)
					sigCh <- tc.input
				}()

				runInteractiveInternal(mainFunc, sigCh)
			})
		}
	}
}

func TestRunInteractiveInternalAllSignals(t *testing.T) {
	for mainName, mainFunc := range getAllMainFuncs() {
		mainFunc := mainFunc
		for _, signal := range getAllSignals() {
			signal := signal

			t.Run(fmt.Sprintf("%s - %s", mainName, signal.String()), func(t *testing.T) {
				t.Parallel()
				sigCh := make(chan os.Signal, 1)
				var m sync.Mutex
				closed := false
				defer func() {
					m.Lock()
					closed = true
					m.Unlock()
					close(sigCh)
				}()

				go func() {
					time.Sleep(time.Millisecond * 100)
					m.Lock()
					if !closed {
						sigCh <- signal
					}
					m.Unlock()
				}()

				runInteractiveInternal(mainFunc, sigCh)
			})
		}
	}
}

func TestRunInteractiveAllSignals(t *testing.T) {
	for mainName, mainFunc := range immediateMainFuncs {
		mainFunc := mainFunc
		for _, signal := range getAllSignals() {
			signal := signal
			t.Run(fmt.Sprintf("%s - %s", mainName, signal), func(t *testing.T) {
				t.Parallel()
				runInteractive(mainFunc)
			})

			t.Run(fmt.Sprintf("%s - %s", mainName, signal), func(t *testing.T) {
				t.Parallel()
				Run(mainFunc, Config{})
			})
		}
	}
}
