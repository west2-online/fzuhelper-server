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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/kitex/pkg/klog"
	kitexzap "github.com/kitex-contrib/obs-opentelemetry/logging/zap"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

var (
	loggerObj         *kitexzap.Logger
	stdErrFileHandler *os.File // 全局变量，避免被 GC 回收
	logFileHandler    *os.File
)

type Logger struct {
	*kitexzap.Logger
}

const (
	permission = 0o755 // 用户具有读/写/执行权限，组用户和其它用户具有读写权限
)

// Init 将会依据服务名在日志目录(参考 constants 设置的常量)下构建相应的标准输出和日志输出
func Init(service string, level int64) {
	var err error
	var pwd string

	// 获取当前目录
	if pwd, err = getCurrentDirectory(); err != nil {
		panic(err)
	}

	logPath := fmt.Sprintf("%s/%s/%s.log", pwd, constants.LogFilePath, service)
	stderrPath := fmt.Sprintf("%s/%s/%s_stderr.log", pwd, constants.LogFilePath, service)

	logFileHandler = checkAndOpenFile(logPath)
	stdErrFileHandler = checkAndOpenFile(stderrPath)

	// 这个系统调用在某些场合（含 tmux 或者跨平台）上并不好用，还是用 golang 内置的好
	// if err = syscall.Dup2(int(stdErrFileHandler.Fd()), int(os.Stderr.Fd())); err != nil {
	// 	panic(fmt.Sprintf("dup2 stderr failed: %v", err))
	// }
	os.Stderr = stdErrFileHandler // 直接替换标准错误输出，需要注意的是，这里不替换 os.Stdout

	loggerObj = DefaultLogger()
	klog.SetLogger(loggerObj)
	klog.SetOutput(logFileHandler)
	klog.SetLevel(klog.Level(level))
	logger.Infof("logger init success, log level: %d", level)
}

// getCurrentDirectory 会返回当前运行的目录
func getCurrentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(dir, "\\", "/"), nil
}

func checkAndOpenFile(path string) *os.File {
	var err error
	var handler *os.File
	if err = os.MkdirAll(filepath.Dir(path), permission); err != nil {
		panic(err)
	}

	handler, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, permission)
	if err != nil {
		panic(err)
	}
	runtime.SetFinalizer(handler, func(fd *os.File) {
		if err := fd.Close(); err != nil {
			logger.Infof(fmt.Sprintf("close file failed %v", err))
		}
	})
	return handler
}

func (l *Logger) GetLoggerObj() *kitexzap.Logger {
	return loggerObj
}

func (l *Logger) Printf(template string, args ...interface{}) {
	l.Infof(template, args...)
}

func GetLogger() *Logger {
	return &Logger{loggerObj}
}
