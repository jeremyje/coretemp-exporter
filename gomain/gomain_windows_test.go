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

//go:build windows
// +build windows

package gomain

import (
	"os"
	"syscall"
)

var (
	handleSignalTestCases = []struct {
		input os.Signal
		want  bool
	}{
		{input: syscall.SIGINT, want: true},
		{input: syscall.SIGTERM, want: true},
		{input: syscall.SIGKILL, want: true},
		{input: syscall.SIGABRT, want: true},
	}
)

func getAllSignals() []os.Signal {
	return []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGILL,
		syscall.SIGTRAP,
		syscall.SIGABRT,
		syscall.SIGBUS,
		syscall.SIGFPE,
		syscall.SIGKILL,
		syscall.SIGSEGV,
		syscall.SIGPIPE,
		syscall.SIGALRM,
		syscall.SIGTERM,
	}
}
