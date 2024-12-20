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
	"io"

	"github.com/cloudwego/kitex/pkg/klog"
)

type KlogLogger struct{}

func GetKlogLogger() *KlogLogger {
	return &KlogLogger{}
}

func (l *KlogLogger) Trace(v ...interface{}) {
	Debug(v...)
}

func (l *KlogLogger) Debug(v ...interface{}) {
	Debug(v...)
}

func (l *KlogLogger) Info(v ...interface{}) {
	Info(v...)
}

func (l *KlogLogger) Notice(v ...interface{}) {
	Info(v...)
}

func (l *KlogLogger) Warn(v ...interface{}) {
	Warn(v...)
}

func (l *KlogLogger) Error(v ...interface{}) {
	Error(v...)
}

func (l *KlogLogger) Fatal(v ...interface{}) {
	Fatal(v...)
}

func (l *KlogLogger) Tracef(format string, v ...interface{}) {
	Debugf(format, v...)
}

func (l *KlogLogger) Debugf(format string, v ...interface{}) {
	Debugf(format, v...)
}

func (l *KlogLogger) Infof(format string, v ...interface{}) {
	Infof(format, v...)
}

func (l *KlogLogger) Noticef(format string, v ...interface{}) {
	Infof(format, v...)
}

func (l *KlogLogger) Warnf(format string, v ...interface{}) {
	Warnf(format, v...)
}

func (l *KlogLogger) Errorf(format string, v ...interface{}) {
	Errorf(format, v...)
}

func (l *KlogLogger) Fatalf(format string, v ...interface{}) {
	Fatalf(format, v...)
}

func (l *KlogLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	Debugf(format, v...)
}

func (l *KlogLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	Debugf(format, v...)
}

func (l *KlogLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	Infof(format, v...)
}

func (l *KlogLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	Infof(format, v...)
}

func (l *KlogLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	Warnf(format, v...)
}

func (l *KlogLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	Errorf(format, v...)
}

func (l *KlogLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	Fatalf(format, v...)
}

func (l *KlogLogger) SetLevel(klog.Level) {
}

func (l *KlogLogger) SetOutput(io.Writer) {
}
