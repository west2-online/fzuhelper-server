// utils/logger.go
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 全局 Logger 变量
var LoggerObj *zap.SugaredLogger

func LoggerInit() {
	// 配置 zap 的日志等级和输出格式
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel), // 设置日志等级
		Development:      false,                               // 非开发模式
		Encoding:         "console",                           // 输出格式（json 或 console）
		OutputPaths:      []string{"stdout"},                  // 输出目标
		ErrorOutputPaths: []string{"stderr"},                  // 错误输出目标
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
	var err error
	logger, err := config.Build() // 创建基础 Logger
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // 确保日志缓冲区被刷新

	// 创建 SugaredLogger
	LoggerObj = logger.Sugar()
}
