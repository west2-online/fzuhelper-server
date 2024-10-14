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
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	LoggerObj    *logrus.Logger
	ErrLoggerObj *logrus.Logger
	once         sync.Once
)

// initLogger 初始化两个日志对象，一个用于普通日志输出，一个用于错误日志输出。
func initLogger() {
	logger := getBaseLogger(os.Stdout)
	errLogger := getBaseLogger(os.Stderr)

	ErrLoggerObj = errLogger
	LoggerObj = logger
}

func getBaseLogger(output io.Writer) *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true) // 设置日志记录器在日志条目中包含调用者信息。
	logger.SetFormatter(getNewFormatter())
	logger.SetLevel(logrus.DebugLevel) // 调整日志等级
	logger.SetOutput(output)
	return logger
}

// 确保 LoggerObj 只初始化一次
func ensureLoggerInit() {
	once.Do(func() {
		initLogger()
	})
}

func Fatalf(template string, args ...interface{}) {
	ensureLoggerInit()
	ErrLoggerObj.Fatalf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	ensureLoggerInit()
	ErrLoggerObj.Errorf(template, args...)
}

func Infof(template string, args ...interface{}) {
	ensureLoggerInit()
	LoggerObj.Infof(template, args...)
}

func Debugf(template string, args ...interface{}) {
	ensureLoggerInit()
	LoggerObj.Debugf(template, args...)
}

func Fatal(args ...interface{}) {
	ensureLoggerInit()
	ErrLoggerObj.Fatal(args...)
}

func Error(args ...interface{}) {
	ensureLoggerInit()
	ErrLoggerObj.Error(args...)
}

func Info(args ...interface{}) {
	ensureLoggerInit()
	LoggerObj.Info(args...)
}

func Debug(args ...interface{}) {
	ensureLoggerInit()
	LoggerObj.Debug(args...)
}
