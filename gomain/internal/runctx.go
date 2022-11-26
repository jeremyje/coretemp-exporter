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
	"os"
	"sync"
)

type RunCtx struct {
	sync.RWMutex
	waitCh chan os.Signal
	closed bool
}

func NewRunCtx() *RunCtx {
	return &RunCtx{
		waitCh: make(chan os.Signal, 1),
		closed: false,
	}
}

func (mc *RunCtx) Kill() {
	mc.RLock()
	if !mc.closed {
		mc.waitCh <- os.Kill
	} else {
		mc.RUnlock()
		return
	}
	mc.RUnlock()
}

func (mc *RunCtx) Wait() {
	<-mc.waitCh
}

func (mc *RunCtx) Close() {
	doClose := false
	mc.Lock()
	if !mc.closed {
		doClose = true
		mc.closed = true
	}

	mc.Unlock()
	if doClose {
		close(mc.waitCh)
	}
}
