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

package logger

import (
	"context"

	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// extractSpanContext 从传入的 ctx 中取出 spanContext
func extractSpanContext(ctx context.Context) []zap.Field {
	// ctx nil值检查
	if ctx == nil {
		return nil
	}
	// 调用 otel sdk 取 spanContext
	spanCtx := oteltrace.SpanContextFromContext(ctx)
	// 不合法，丢
	if !spanCtx.IsValid() {
		return nil
	}

	return []zap.Field{
		zap.String("trace_id", spanCtx.TraceID().String()),
		zap.String("span_id", spanCtx.SpanID().String()),
	}
}
