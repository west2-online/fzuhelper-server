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
	"go.uber.org/zap/zapcore"
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
	ctxLogger.log(zapcore.DebugLevel, msg, fields...)
}

func (ctxLogger *ContextLogger) Debugf(template string, args ...interface{}) {
	ctxLogger.logf(zapcore.DebugLevel, template, args...)
}

func (ctxLogger *ContextLogger) Info(msg string, fields ...zap.Field) {
	ctxLogger.log(zapcore.InfoLevel, msg, fields...)
}

func (ctxLogger *ContextLogger) Infof(template string, args ...interface{}) {
	ctxLogger.logf(zapcore.InfoLevel, template, args...)
}

func (ctxLogger *ContextLogger) Warn(msg string, fields ...zap.Field) {
	ctxLogger.log(zapcore.WarnLevel, msg, fields...)
}

func (ctxLogger *ContextLogger) Warnf(template string, args ...interface{}) {
	ctxLogger.logf(zapcore.WarnLevel, template, args...)
}

func (ctxLogger *ContextLogger) Error(msg string, fields ...zap.Field) {
	ctxLogger.log(zapcore.ErrorLevel, msg, fields...)
}

func (ctxLogger *ContextLogger) Errorf(template string, args ...interface{}) {
	ctxLogger.logf(zapcore.ErrorLevel, template, args...)
}

func (ctxLogger *ContextLogger) Fatal(msg string, fields ...zap.Field) {
	ctxLogger.log(zapcore.FatalLevel, msg, fields...)
}

func (ctxLogger *ContextLogger) Fatalf(template string, args ...interface{}) {
	ctxLogger.logf(zapcore.FatalLevel, template, args...)
}

func (ctxLogger *ContextLogger) log(lvl zapcore.Level, msg string, fields ...zap.Field) {
	// 注入 span 信息
	control.log(lvl, msg,
		append(fields, ctxLogger.spanFields...)...)

	if lvl >= errorSpanLevel {
		markSpanError(ctxLogger.ctx, msg)
	}
}

func (ctxLogger *ContextLogger) logf(lvl zapcore.Level, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	// 注入 span 信息
	control.log(lvl, msg, ctxLogger.spanFields...)

	if lvl >= errorSpanLevel {
		markSpanError(ctxLogger.ctx, msg)
	}
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
