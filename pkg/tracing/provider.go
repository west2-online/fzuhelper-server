/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tracing

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	uptraceDSNKey = "uptrace-dsn"
)

func NewOtelProvider(serviceName string, endpoint string, uptraceDSN string) func(context.Context) error {
	ctx := context.Background()

	res := getResource(ctx, serviceName)

	// includes tracing and metics provider, from kitexotel api
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(serviceName),
		provider.WithExportEndpoint(endpoint),
		provider.WithHeaders(map[string]string{
			uptraceDSNKey: uptraceDSN,
		}),
		provider.WithResource(res),
		provider.WithInsecure())

	// manually written logger provider
	lp := newOtelLoggerProvider(serviceName, endpoint, uptraceDSN)

	// return shutdown func
	return func(ctx context.Context) error {
		var err error

		if err = p.Shutdown(ctx); err != nil {
			return err
		}

		if err = lp.Shutdown(ctx); err != nil {
			otel.Handle(err) // handle by otel
		}

		return err
	}
}

// newOtelLoggerProvider 手动初始化 LoggerProvider
func newOtelLoggerProvider(serviceName string, endpoint string, uptraceDSN string) *sdklog.LoggerProvider {
	ctx := context.Background()

	res := getResource(ctx, serviceName)

	// log exporter
	logExp, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithHeaders(map[string]string{
			uptraceDSNKey: uptraceDSN,
		}),
		otlploggrpc.WithInsecure())
	if err != nil {
		klog.Fatalf("failed to create otlp log exporter: %s", err)
		return nil
	}

	// log processor
	bp := sdklog.NewBatchProcessor(logExp)

	// logger provider
	lp := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(bp))

	global.SetLoggerProvider(lp)

	return lp
}

// getResource 一个 tracing/metrics/logging 通用的 Resource
func getResource(ctx context.Context, serviceName string) *resource.Resource {
	// 参见 https://github.com/kitex-contrib/obs-opentelemetry/blob/main/provider/provider.go 下 newResource()
	res, err := resource.New(ctx,
		resource.WithHost(),
		resource.WithFromEnv(),
		resource.WithProcessPID(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName)), // service.name
	)
	if err != nil {
		return resource.Default()
	}
	return res
}
