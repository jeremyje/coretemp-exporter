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
	"os"

	"github.com/jeremyje/coretemp-exporter/drivers/common"
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

func getDefaultAttributes() []attribute.KeyValue {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = os.Getenv("COMPUTERNAME")
	}

	return []attribute.KeyValue{
		attribute.Key("hostname").String(hostname),
	}
}

type metricsSink struct {
	CPUCoreTemperature    asyncfloat64.Gauge
	CPUCoreLoad           asyncint64.Gauge
	CPUThermalJunctionMax asyncfloat64.Gauge
	CPUInfoPollCount      syncfloat64.Counter
	CPUSpeed              asyncfloat64.Gauge
	CPUMultiplier         asyncfloat64.Gauge
	lastValue             *common.HardwareInfo
}

func (m *metricsSink) ObserveAsync(ctx context.Context) {
	m.Observe(ctx, m.lastValue)
}

func (m *metricsSink) Observe(ctx context.Context, info *common.HardwareInfo) {
	if info == nil {
		return
	}

	m.lastValue = info
	attrs := getDefaultAttributes()
	curAttrs := append(attrs, attribute.String("model", info.CPUName))

	for core, tempC := range info.TemperatureCelcius {
		m.CPUCoreTemperature.Observe(ctx, tempC, append(curAttrs, attribute.Int("core", core))...)
	}
	m.CPUInfoPollCount.Add(ctx, 1, curAttrs...)

	m.CPUSpeed.Observe(ctx, info.CPUSpeed*1000*1000, append(
		curAttrs,
		attribute.Int("cpu_count", info.CPUCount),
		attribute.Int("core_count", info.CoreCount),
		attribute.Float64("fsb_speed", info.FSBSpeed),
		attribute.Int("voltage", int(info.VID)),
	)...)

	m.CPUMultiplier.Observe(ctx, info.Multiplier, append(
		curAttrs,
		attribute.Int("cpu_count", info.CPUCount),
		attribute.Int("core_count", info.CoreCount),
		attribute.Float64("fsb_speed", info.FSBSpeed),
		attribute.Int("voltage", int(info.VID)),
	)...)

	for core, load := range info.Load {
		m.CPUCoreLoad.Observe(ctx, int64(load), append(curAttrs, attribute.Int("core", core))...)
	}
	for core, tjMax := range info.TJMax {
		m.CPUThermalJunctionMax.Observe(ctx, tjMax, append(curAttrs, attribute.Int("core", core))...)
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
	cpuThermalJunctionMax, err := meter.AsyncFloat64().Gauge("cpu_tj_max", instrument.WithDescription("CPU Load percentage (0-100)"), instrument.WithUnit("C"))
	if err != nil {
		return nil, err
	}
	cpuInfoPollCount, err := meter.SyncFloat64().Counter("cpu_core_poll", instrument.WithDescription("Number of times the CPU temperature has been polled."))
	if err != nil {
		return nil, err
	}
	cpuSpeed, err := meter.AsyncFloat64().Gauge("cpu_speed", instrument.WithDescription("CPU Core Speeds"))
	if err != nil {
		return nil, err
	}
	cpuMultiplier, err := meter.AsyncFloat64().Gauge("cpu_multiplier", instrument.WithDescription("FSB Multiplier for the CPU."))
	if err != nil {
		return nil, err
	}

	sink := &metricsSink{
		CPUCoreTemperature:    cpuCoreTemperature,
		CPUCoreLoad:           cpuCoreLoad,
		CPUThermalJunctionMax: cpuThermalJunctionMax,
		CPUInfoPollCount:      cpuInfoPollCount,
		CPUSpeed:              cpuSpeed,
		CPUMultiplier:         cpuMultiplier,
	}

	meter.RegisterCallback([]instrument.Asynchronous{cpuCoreTemperature, cpuCoreLoad, cpuThermalJunctionMax, cpuSpeed, cpuMultiplier}, func(ctx context.Context) {
		sink.ObserveAsync(ctx)
	})

	return sink, nil
}
