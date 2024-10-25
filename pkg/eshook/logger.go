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

package eshook

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"go.uber.org/zap"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/client"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// InitLoggerWithHook 初始化带有EsHook的logger
// index: 索引的名字
func InitLoggerWithHook(index string) {
	if config.Elasticsearch == nil {
		return
	}

	c, err := client.NewEsClient()
	if err != nil {
		panic(err)
	}

	if !client.Connected(c) {
		logger.Warn("es not worked!")
		return
	}

	hook := NewElasticHook(c, config.Elasticsearch.Host, index)
	v := logger.DefaultLogger(zap.Hooks(hook.Fire))
	klog.SetLogger(v)
}
