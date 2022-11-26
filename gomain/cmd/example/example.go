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
	"context"
	"net/http"

	"github.com/jeremyje/gomain"
)

func main() {
	gomain.Run(appMain, gomain.Config{
		ServiceName:        "Example",
		ServiceDescription: "Example Service Description",
		Command:            "",
	})
}

func appMain(waitFunc func()) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte(req.URL.Path))
	})
	s := &http.Server{
		Handler: mux,
		Addr:    "localhost:0",
	}

	go func() {
		waitFunc()
		ctx := context.Background()
		s.Shutdown(ctx)
	}()

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
