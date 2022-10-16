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
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/jeremyje/coretempsdk-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric"
)

type Args struct {
	Endpoint string
	Interval time.Duration
}

func Run(args *Args) error {
	ctx := context.Background()
	prom := prometheus.New()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = os.Getenv("COMPUTERNAME")
	}

	provider := metric.NewMeterProvider(metric.WithReader(prom))
	meter := provider.Meter("github.com/jeremyje/coretemp-exporter")

	attrs := []attribute.KeyValue{
		attribute.Key("hostname").String(hostname),
	}

	cpuCoreTemperature, err := meter.AsyncFloat64().Gauge("cpu_core_temperature", instrument.WithDescription("Temperature of CPU Cores"), instrument.WithUnit("C"))
	if err != nil {
		return err
	}
	cpuInfoPollCount, err := meter.SyncFloat64().Counter("cpu_core_poll", instrument.WithDescription("Temperature of CPU Cores"))
	if err != nil {
		return err
	}

	ticker := time.NewTicker(args.Interval)
	done := make(chan bool)
	defer func() {
		done <- true
		close(done)
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				info, err := coretempsdk.GetCoreTempInfo()
				if err != nil {
					log.Printf("ERROR: %s", err)
				}

				curAttrs := append(attrs, attribute.String("model", info.CPUName))

				cpuCoreTemperature.Observe(ctx, float64(info.TemperatureCelcius[0]), curAttrs...)
				cpuInfoPollCount.Add(ctx, 1, curAttrs...)
				prom.ForceFlush(ctx)
			}
		}
	}()

	return serve(ctx, args, promhttp.Handler())
}

func serve(ctx context.Context, args *Args, h http.Handler) error {
	addr := args.Endpoint
	s := &http.Server{
		Addr:    addr,
		Handler: h,
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	errCh := make(chan error)
	log.Printf("Serving on %s", addr)
	go func() {
		errCh <- s.Serve(lis)
	}()

	stopCtx, _ := signal.NotifyContext(ctx, os.Interrupt)
	select {
	case <-errCh:
	case <-stopCtx.Done():
		lis.Close()
	}
	return nil
}
