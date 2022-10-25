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
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

func serveAsync(ctx context.Context, args *Args, h http.Handler) (func(), func(), error) {
	addr := args.Endpoint
	s := &http.Server{
		Addr:    addr,
		Handler: h,
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return func() {}, func() {}, err
	}
	errCh := make(chan error)
	log.Printf("Serving on %s", addr)

	go func() {
		errCh <- s.Serve(lis)
		close(errCh)
	}()

	return func() {
			stopCtx, _ := signal.NotifyContext(ctx, os.Interrupt)
			select {
			case <-errCh:
			case <-stopCtx.Done():
				lis.Close()
			}
		}, func() {
			lis.Close()
		}, err
}
