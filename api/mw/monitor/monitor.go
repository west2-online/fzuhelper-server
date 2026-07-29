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
	apiMonitorStop      = func(context.Context) {}
)

// StartAPIMonitor 启动 API 监控后台检查，并返回用于优雅停止监控的函数。
// 停止函数会等待后台检查 goroutine 退出，便于服务关闭时完成资源清理。
func StartAPIMonitor(cfg MonitorConfig) func(context.Context) {
	// 监控进程生命周期与服务进程一致，只允许初始化一次。
	apiMonitorStartOnce.Do(func() {
		apiMonitorInstance = newAPIMonitor(cfg)
		if !apiMonitorInstance.enabled() {
			// 未启用监控时不启动后台任务，返回默认的空操作停止函数。
			return
		}

		interval := apiMonitorInstance.checkInterval()
		if interval <= 0 {
			// 无效的检查间隔不能用于创建 ticker，避免启动时 panic。
			return
		}

		// cancel 用于通知后台任务退出，done 用于确认检查 goroutine 已结束。
		ctx, cancel := context.WithCancel(context.Background())
		checkDone := make(chan struct{})
		go func() {
			defer close(checkDone)

			ticker := time.NewTicker(interval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					// 定期清理滑动窗口并检查各路由的错误率。
					apiMonitorInstance.check()
				case <-ctx.Done():
					// 服务关闭时优先响应取消信号，不再执行新的检查。
					return
				}
			}
		}()

		apiMonitorStop = func(ctx context.Context) {
			cancel()
			// 等待检查循环退出，避免服务关闭后仍残留后台任务。
			select {
			case <-checkDone:
			case <-ctx.Done():
				return
			}
		}
	})

	// API 服务会在 OnShutdown hook 中调用该函数。
	return apiMonitorStop
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
