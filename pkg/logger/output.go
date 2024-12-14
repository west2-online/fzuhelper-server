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
	"github.com/cloudwego/kitex/pkg/klog"
)

func Debug(args ...interface{}) {
	klog.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	klog.Debugf(template, args...)
}

func Info(args ...interface{}) {
	klog.Info(args...)
}

func Infof(template string, args ...interface{}) {
	klog.Infof(template, args...)
}

func Warn(args ...interface{}) {
	klog.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	klog.Warnf(template, args...)
}

func Error(args ...interface{}) {
	klog.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	klog.Errorf(template, args...)
}

func Fatal(args ...interface{}) {
	klog.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	klog.Fatalf(template, args...)
}

// LErrorf Equals Errorf less one stack
func LErrorf(template string, args ...interface{}) {
	loggerObj.Errorf(template, args...)
}
