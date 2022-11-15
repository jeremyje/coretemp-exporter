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
	"net/http"

	pb "github.com/jeremyje/coretemp-exporter/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type metricsSink struct {
	CPUCoreTemperature asyncfloat64.Gauge
	CPUCoreLoad        asyncint64.Gauge
	CPUInfoPollCount   syncfloat64.Counter
	CPUFrequency       asyncfloat64.Gauge
	CPUFSBFrequency    asyncfloat64.Gauge
	lastValue          *pb.MachineMetrics
}

func (m *metricsSink) ObserveAsync(ctx context.Context) {
	m.Observe(ctx, m.lastValue)
}

func (m *metricsSink) Observe(ctx context.Context, mm *pb.MachineMetrics) {
	if mm == nil {
		return
	}

	m.lastValue = mm

	for _, device := range mm.GetDevice() {
		attrs := []attribute.KeyValue{
			attribute.Key("hostname").String(mm.GetName()),
			attribute.Key("name").String(device.GetName()),
			attribute.Key("kind").String(device.GetKind()),
		}
		curAttrs := attrs

		if device.GetCpu() != nil {
			cpuMetrics := device.GetCpu()
			for core, tempC := range cpuMetrics.GetTemperature() {
				m.CPUCoreTemperature.Observe(ctx, tempC, append(curAttrs, attribute.Int("core", core))...)
			}
			m.CPUInfoPollCount.Add(ctx, 1, curAttrs...)

			m.CPUFrequency.Observe(ctx, cpuMetrics.GetFrequencyMhz()*1000*1000, append(
				curAttrs,
				attribute.Int("core_count", int(cpuMetrics.GetNumCores())),
			)...)

			m.CPUFSBFrequency.Observe(ctx, cpuMetrics.GetFsbFrequencyMhz()*1000*1000, append(
				curAttrs,
				attribute.Int("core_count", int(cpuMetrics.GetNumCores())),
			)...)

			for core, load := range cpuMetrics.GetLoad() {
				m.CPUCoreLoad.Observe(ctx, int64(load), append(curAttrs, attribute.Int("core", core))...)
			}
		}
	}

}

func newMetricsSink(ctx context.Context) (*metricsSink, http.Handler, error) {
	prometheusExporter := otelprom.New()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(prometheusExporter))
	meter := provider.Meter("github.com/jeremyje/coretemp-exporter")
	sink, err := newMetrics(meter)
	if err != nil {
		return nil, nil, err
	}
	registry := prometheus.NewRegistry()
	registry.Register(collectors.NewBuildInfoCollector())
	registry.Register(collectors.NewGoCollector())
	registry.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
		ReportErrors: true,
	}))
	registry.Register(prometheusExporter.Collector)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	return sink, h, nil
}

func newMetrics(meter metric.Meter) (*metricsSink, error) {
	cpuCoreTemperature, err := meter.AsyncFloat64().Gauge("cpu_core_temperature", instrument.WithDescription("Temperature of a CPU Core in Celcius"), instrument.WithUnit("C"))
	if err != nil {
		return nil, err
	}
	cpuCoreLoad, err := meter.AsyncInt64().Gauge("cpu_core_load", instrument.WithDescription("CPU Load percentage (0-100)"), instrument.WithUnit("C"))
	if err != nil {
		return nil, err
	}
	cpuInfoPollCount, err := meter.SyncFloat64().Counter("cpu_core_poll", instrument.WithDescription("Number of times the CPU temperature has been polled."))
	if err != nil {
		return nil, err
	}
	cpuFrequency, err := meter.AsyncFloat64().Gauge("cpu_frequency", instrument.WithDescription("CPU Core Frequency"))
	if err != nil {
		return nil, err
	}
	cpuFSBFrequency, err := meter.AsyncFloat64().Gauge("cpu_fsb_frequency", instrument.WithDescription("CPU Front Side Bus Frequency"))
	if err != nil {
		return nil, err
	}

	sink := &metricsSink{
		CPUCoreTemperature: cpuCoreTemperature,
		CPUCoreLoad:        cpuCoreLoad,
		CPUInfoPollCount:   cpuInfoPollCount,
		CPUFrequency:       cpuFrequency,
		CPUFSBFrequency:    cpuFSBFrequency,
	}

	meter.RegisterCallback([]instrument.Asynchronous{cpuCoreTemperature, cpuCoreLoad, cpuFrequency, cpuFSBFrequency}, func(ctx context.Context) {
		sink.ObserveAsync(ctx)
	})

	return sink, nil
}
