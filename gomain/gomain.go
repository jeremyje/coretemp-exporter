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
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jeremyje/gomain/internal"
)

type MainCtx interface {
	Wait()
}

type MainFunc func(func()) error

type Config struct {
	ServiceName        string
	ServiceDescription string
	Command            string
}

func Run(f MainFunc, cfg Config) {
	platformRun(f, cfg)
}

func getTerminalSignals() []os.Signal {
	return []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGABRT}
}

func runInteractive(f MainFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, getTerminalSignals()...)
	defer func() {
		signal.Stop(sigCh)
		close(sigCh)
	}()
	runInteractiveInternal(f, sigCh)
}

func runInteractiveInternal(f MainFunc, sigCh chan os.Signal) {
	mainErrCh := make(chan error, 1)

	mc := internal.NewRunCtx()
	defer mc.Close()

	go func() {
		mainErrCh <- f(mc.Wait)
		close(mainErrCh)
	}()

	select {
	case err := <-mainErrCh:
		handleError(err)
		return
	case sig := <-sigCh:
		if handleSignal(sig) {
			signal.Stop(sigCh)
			mc.Kill()
		}
	}
}

func handleError(err error) {
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
}

func handleSignalBase(sig os.Signal) bool {
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		return true
	case syscall.SIGKILL, syscall.SIGABRT:
		logStackDump()
		return true
	default:
		return false
	}
}
