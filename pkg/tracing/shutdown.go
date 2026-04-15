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

package tracing

import (
	"context"

	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// ProviderShutdown 封装 OtelProvider 的原生关闭函数
func ProviderShutdown(shutdown func(context.Context) error, logTemplate string) func() {
	return ProviderShutdownWithContext(shutdown, context.Background(), logTemplate)
}

// ProviderShutdownWithContext 封装 OtelProvider 的原生关闭函数，支持自定义ctx
func ProviderShutdownWithContext(shutdown func(context.Context) error, ctx context.Context, logTemplate string) func() {
	return func() {
		if err := shutdown(ctx); err != nil {
			logger.Fatalf(logTemplate, err)
		}
	}
}
