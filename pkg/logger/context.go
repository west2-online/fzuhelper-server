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
	"fmt"

	"go.uber.org/zap"
)

func WithCtx(ctx context.Context) *ContextLogger {
	spanFields := extractSpanContext(ctx)

	ctxLogger := &ContextLogger{
		ctx:        ctx,
		spanFields: spanFields,
	}
	return ctxLogger
}

type ContextLogger struct {
	ctx        context.Context
	spanFields []zap.Field
}

func (ctxLogger *ContextLogger) Debug(msg string, fields ...zap.Field) {
	control.debug(msg,
		append(fields, ctxLogger.spanFields...)...)
}

func (ctxLogger *ContextLogger) Debugf(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	control.debug(msg,
		ctxLogger.spanFields...)
}

func (ctxLogger *ContextLogger) Info(msg string, fields ...zap.Field) {
	control.info(msg,
		append(fields, ctxLogger.spanFields...)...)
}

func (ctxLogger *ContextLogger) Infof(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	control.info(msg,
		ctxLogger.spanFields...)
}

func (ctxLogger *ContextLogger) Warn(msg string, fields ...zap.Field) {
	control.warn(msg,
		append(fields, ctxLogger.spanFields...)...)
}

func (ctxLogger *ContextLogger) Warnf(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	control.warn(msg,
		ctxLogger.spanFields...)
}

func (ctxLogger *ContextLogger) Error(msg string, fields ...zap.Field) {
	control.error(msg,
		append(fields, ctxLogger.spanFields...)...)
}

func (ctxLogger *ContextLogger) Errorf(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	control.error(msg,
		ctxLogger.spanFields...)
}

func (ctxLogger *ContextLogger) Fatal(msg string, fields ...zap.Field) {
	control.fatal(msg,
		append(fields, ctxLogger.spanFields...)...)
}

func (ctxLogger *ContextLogger) Fatalf(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	control.fatal(msg,
		ctxLogger.spanFields...)
}

// ----------------------------------------
// 以下是 *WithCtx(ctx, ...) 系列函数实现
// ----------------------------------------

func DebugWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithCtx(ctx).Debug(msg, fields...)
}

func DebugfWithCtx(ctx context.Context, template string, args ...interface{}) {
	WithCtx(ctx).Debugf(template, args...)
}

func InfoWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithCtx(ctx).Info(msg, fields...)
}

func InfofWithCtx(ctx context.Context, template string, args ...interface{}) {
	WithCtx(ctx).Infof(template, args...)
}

func WarnWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithCtx(ctx).Warn(msg, fields...)
}

func WarnfWithCtx(ctx context.Context, template string, args ...interface{}) {
	WithCtx(ctx).Warnf(template, args...)
}

func ErrorWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithCtx(ctx).Error(msg, fields...)
}

func ErrorfWithCtx(ctx context.Context, template string, args ...interface{}) {
	WithCtx(ctx).Errorf(template, args...)
}

func FatalWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithCtx(ctx).Fatal(msg, fields...)
}

func FatalfWithCtx(ctx context.Context, template string, args ...interface{}) {
	WithCtx(ctx).Fatalf(template, args...)
}
