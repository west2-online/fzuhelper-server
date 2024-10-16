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
	"os"

	kitexzap "github.com/kitex-contrib/obs-opentelemetry/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Enc zapcore.Encoder
	Ws  zapcore.WriteSyncer
	lvl zapcore.Level
}

func NewLogger(lvl zapcore.Level, cfg Config, options ...zap.Option) *kitexzap.Logger {
	if cfg.Enc == nil {
		cfg.Enc = defaultEnc()
	}
	if cfg.Ws == nil {
		cfg.Ws = defaultWs()
	}
	cfg.lvl = lvl

	var ops []kitexzap.Option
	ops = append(ops, kitexzap.WithCoreEnc(cfg.Enc))
	ops = append(ops, kitexzap.WithCoreWs(cfg.Ws))
	ops = append(ops, kitexzap.WithCoreLevel(zap.NewAtomicLevelAt(cfg.lvl)))
	ops = append(ops, kitexzap.WithZapOptions(options...))
	return kitexzap.NewLogger(ops...)
}

func DefaultLogger(options ...zap.Option) *kitexzap.Logger {
	var ops []kitexzap.Option
	ops = append(ops, kitexzap.WithCoreEnc(defaultEnc()))
	ops = append(ops, kitexzap.WithCoreWs(defaultWs()))
	ops = append(ops, kitexzap.WithCoreLevel(zap.NewAtomicLevelAt(defaultLvl())))
	ops = append(ops, kitexzap.WithZapOptions(options...))
	return kitexzap.NewLogger(ops...)
}

func defaultEnc() zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
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
	}

	return zapcore.NewConsoleEncoder(cfg)
}

func defaultWs() zapcore.WriteSyncer {
	return os.Stdout
}

func defaultLvl() zapcore.Level {
	return zapcore.DebugLevel
}
