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
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	LoggerObj *zap.SugaredLogger
	once      sync.Once
)

// 初始化 Logger 的函数
func initLogger() {
	// 配置 zap 的日志等级和输出格式
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel), // 设置日志等级
		Development:      false,                                // 非开发模式
		Encoding:         "console",                            // 输出格式（json 或 console）
		OutputPaths:      []string{"stdout"},                   // 输出目标
		ErrorOutputPaths: []string{"stderr"},                   // 错误输出目标
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder, // 日志等级大写
			EncodeTime:     zapcore.ISO8601TimeEncoder,  // 时间格式
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	logger, err := config.Build(zap.AddCallerSkip(1)) // 创建基础 Logger
	if err != nil {
		panic(err)
	}

	// 创建 SugaredLogger
	LoggerObj = logger.Sugar()
}

// 确保 LoggerObj 只初始化一次
// 使用 init() 函数来替代 syncOnce 是个方法，而且不需要额外的代码来进行 check
// 但是这样损失了更多的 DIY 特性，比如可以在初始化的时候传入参数
func ensureLoggerInit() {
	once.Do(func() {
		initLogger()
	})
}

func Fatalf(template string, args ...interface{}) {
	ensureLoggerInit()
	LoggerObj.Fatalf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	ensureLoggerInit()
	LoggerObj.Errorf(template, args...)
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
	LoggerObj.Fatal(args)
}

func Info(args ...interface{}) {
	ensureLoggerInit()
	LoggerObj.Info(args)
}

func Error(args ...interface{}) {
	ensureLoggerInit()
	LoggerObj.Error(args)
}

func Debug(args ...interface{}) {
	ensureLoggerInit()
	LoggerObj.Debug(args)
}
