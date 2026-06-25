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

package monitor

import (
	"context"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// MonitorConfig 是 API 监控的配置参数。
type MonitorConfig struct {
	Enabled       bool
	Window        time.Duration
	CheckInterval time.Duration
	Threshold     float64
	MinRequests   int64
	Cooldown      time.Duration
	Blacklist     map[string]struct{}
}

var (
	apiMonitorInstance  *apiMonitor
	apiMonitorStartOnce sync.Once
)

func StartAPIMonitor(cfg MonitorConfig) {
	apiMonitorStartOnce.Do(func() {
		apiMonitorInstance = newAPIMonitor(cfg)
		if !apiMonitorInstance.enabled() {
			return
		}

		go func() {
			for {
				time.Sleep(apiMonitorInstance.checkInterval())
				apiMonitorInstance.check()
			}
		}()
	})
}

func APIMonitorMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if apiMonitorInstance == nil || !apiMonitorInstance.enabled() {
			c.Next(ctx)
			return
		}

		c.Next(ctx)

		now := time.Now()
		event := buildRequestEvent(
			routeName(c),
			c.Response.Body(),
			getTraceIDFromContext(ctx),
			now,
		)
		apiMonitorInstance.record(event)
	}
}

func getTraceIDFromContext(ctx context.Context) string {
	spanCtx := oteltrace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return ""
	}
	return spanCtx.TraceID().String()
}

func routeName(c *app.RequestContext) string {
	if route := c.FullPath(); route != "" {
		return route
	}
	return string(c.Path())
}

// MarkAPIMonitorPanic 在 panic 恢复时记录一条错误事件，供滑动窗口统计。
func MarkAPIMonitorPanic(ctx context.Context, c *app.RequestContext) {
	if apiMonitorInstance == nil || !apiMonitorInstance.enabled() {
		return
	}
	event := requestEvent{
		route:     routeName(c),
		errorCode: errno.InternalServiceErrorCode,
		traceID:   getTraceIDFromContext(ctx),
		timestamp: time.Now(),
	}
	apiMonitorInstance.record(event)
}
