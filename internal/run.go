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
	"time"

	"github.com/jeremyje/coretemp-exporter/drivers"
	"github.com/jeremyje/gomain"
)

type Args struct {
	Endpoint              string
	Interval              time.Duration
	Log                   string
	Console               bool
	ServiceControlCommand string
}

func Run(args *Args) {
	gomain.Run(func(wait func()) error {
		return run(args, wait)
	}, gomain.Config{
		ServiceName:        "coretemp-exporter",
		ServiceDescription: "Reports CPU Core Temperatures to Prometheus",
		Command:            args.ServiceControlCommand,
	})
}

func run(args *Args, wait func()) error {
	var handler http.Handler
	sinks := []HardwareDataSink{}
	ctx := context.Background()
	handler = http.NewServeMux()

	if args.Endpoint != "" {
		metrics, promHandler, err := newMetricsSink(ctx)
		if err != nil {
			return err
		}
		sinks = append(sinks, metrics)
		handler = promHandler
	}

	if args.Console {
		sinks = append(sinks, &consoleSink{})
	}

	if args.Log != "" {
		fs, err := newFileSink(args.Log)
		if err != nil {
			return err
		}

		sinks = append(sinks, fs)
	}

	ms := newMultiSink(sinks...)

	ticker := time.NewTicker(args.Interval)
	done := make(chan bool)
	go func() {
		ctx := context.Background()

		d := drivers.New()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				info, err := d.Get()
				if err != nil {
					log.Printf("ERROR: %s", err)
				}

				ms.Observe(ctx, info)
			}
		}
	}()

	addr := args.Endpoint
	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("Serving on %s", addr)

	go func() {
		wait()
		ctx := context.Background()
		s.Shutdown(ctx)
		ticker.Stop()
		done <- true
		close(done)
	}()

	return s.Serve(lis)
}
