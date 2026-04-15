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

	"github.com/kitex-contrib/obs-opentelemetry/provider"
)

func NewOtelProvider(serviceName string, endpoint string) func(context.Context) error {
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(serviceName),
		provider.WithExportEndpoint(endpoint),
		provider.WithInsecure())
	return shutdownFunc(p)
}

func shutdownFunc(p provider.OtelProvider) func(context.Context) error {
	return func(ctx context.Context) error {
		return p.Shutdown(ctx)
	}
}
