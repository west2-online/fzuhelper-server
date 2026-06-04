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

package mw

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/alibaba/sentinel-golang/core/flow"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

type rateLimitRuleConfig struct {
	Resource    string  `mapstructure:"resource"`
	QPS         float64 `mapstructure:"qps"`
	QPM         float64 `mapstructure:"qpm"`
	Description string  `mapstructure:"description"`
}

type rateLimitGlobalConfig struct {
	Enabled  bool    `mapstructure:"enabled"`
	Resource string  `mapstructure:"resource"`
	QPS      float64 `mapstructure:"qps"`
}

type rateLimitInterfaceConfig struct {
	Enabled bool                  `mapstructure:"enabled"`
	Rules   []rateLimitRuleConfig `mapstructure:"rules"`
}

type rateLimitConfig struct {
	Enabled   bool                     `mapstructure:"enabled"`
	Global    rateLimitGlobalConfig    `mapstructure:"global"`
	Interface rateLimitInterfaceConfig `mapstructure:"interface"`
}

type rateLimitRuntimeState struct {
	Enabled          bool
	InterfaceEnabled bool
	RuleCount        int
}

var (
	rateLimitState    atomic.Value
	rateLimitReloadMu sync.Mutex
)

// ReloadSentinelRateLimit 校验配置并原子性地重载 Sentinel 限流规则。
func ReloadSentinelRateLimit(ctx context.Context, cfg *rateLimitConfig, source string) (int, error) {
	rateLimitReloadMu.Lock()
	defer rateLimitReloadMu.Unlock()

	rules, state, err := BuildRateLimitRulesFromConfig(cfg)
	if err != nil {
		return 0, errno.InternalServiceError.WithMessage("build rate limit rules failed").WithError(err)
	}

	updated, err := flow.LoadRules(rules)
	if err != nil {
		return 0, errno.InternalServiceError.WithMessage("load rate limit rules failed").WithError(err)
	}

	storeRateLimitRuntimeState(state)
	logger.Infof("rate limit reload success, source=%s, rules=%d, updated=%v", source, state.RuleCount, updated)

	return state.RuleCount, nil
}

// getRateLimitRuntimeState 读取当前限流开关和规则数量状态。
func getRateLimitRuntimeState() rateLimitRuntimeState {
	state := rateLimitState.Load()
	if state == nil {
		return rateLimitRuntimeState{}
	}

	return state.(rateLimitRuntimeState)
}

// storeRateLimitRuntimeState 保存当前限流运行状态供请求中间件读取。
func storeRateLimitRuntimeState(state rateLimitRuntimeState) {
	rateLimitState.Store(state)
}
