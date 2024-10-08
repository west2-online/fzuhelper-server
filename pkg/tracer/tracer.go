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

package tracer

import (
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"

	"github.com/west2-online/fzuhelper-server/config"
)

func InitJaeger(service string) {
	cfg := &jaegerconfig.Configuration{
		Disabled: false,
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: config.Jaeger.Addr,
		},
	}

	cfg.ServiceName = service

	tracer, _, err := cfg.NewTracer(
		jaegerconfig.Logger(jaeger.StdLogger),
		jaegerconfig.ZipkinSharedRPCSpan(true),
	)
	if err != nil {
		panic(fmt.Sprintf("cannot init jaeger: %v\n", err))
	}

	opentracing.SetGlobalTracer(tracer)
}
