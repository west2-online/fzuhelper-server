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

	sentinel "github.com/alibaba/sentinel-golang/api"
	sentinelbase "github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

const (
	rateLimitMessage    = "服务器当前处于请求高峰，请稍后再试"
	qpsStatIntervalInMs = 1000
	qpmStatIntervalInMs = 60000
)

// InterfaceRateLimit 接口限流函数
func InterfaceRateLimit() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		state := getRateLimitRuntimeState()
		if !state.Enabled || !state.InterfaceEnabled {
			c.Next(ctx)
			return
		}

		resource := string(c.Method()) + ":" + string(c.Path())
		entry, blockErr := sentinel.Entry(resource, sentinel.WithTrafficType(sentinelbase.Inbound))
		if blockErr != nil {
			logger.Warnf("frequent requests have been rejected by the gateway. resource: %s, clientIP: %v\n", resource, c.ClientIP())
			pack.RespError(c, errno.InternalServiceError.WithMessage(rateLimitMessage))
			c.Abort()
			return
		}

		defer entry.Exit()
		c.Next(ctx)
	}
}

// BuildRateLimitRules 构建限流规则
func BuildRateLimitRules() ([]*flow.Rule, error) {
	rules, _, err := BuildRateLimitRulesFromConfig(rateLimitConfigFromRuntime())
	return rules, err
}

// BuildRateLimitRulesFromConfig 从配置对象校验并构建 Sentinel 限流规则。
func BuildRateLimitRulesFromConfig(cfg *rateLimitConfig) ([]*flow.Rule, rateLimitRuntimeState, error) {
	state := rateLimitRuntimeState{}
	if cfg == nil || !cfg.Enabled {
		return nil, state, nil
	}
	state.Enabled = true
	state.InterfaceEnabled = cfg.Interface.Enabled

	rules := make([]*flow.Rule, 0)

	if cfg.Global.Enabled {
		if cfg.Global.Resource == "" {
			return nil, state, errno.InternalServiceError.WithMessage("ratelimit.global.resource is empty")
		}
		if cfg.Global.QPS <= 0 {
			return nil, state, errno.InternalServiceError.WithMessage("ratelimit.global.qps must be greater than 0")
		}
		rules = append(rules, qpsRule(cfg.Global.Resource, cfg.Global.QPS))
	}

	if cfg.Interface.Enabled {
		seen := make(map[string]struct{})
		for _, item := range cfg.Interface.Rules {
			if item.Resource == "" {
				return nil, state, errno.InternalServiceError.WithMessage("ratelimit.interface.rules.resource is empty")
			}
			if _, ok := seen[item.Resource]; ok {
				return nil, state, errno.InternalServiceError.WithMessage("duplicate rate limit resource: " + item.Resource)
			}
			seen[item.Resource] = struct{}{}

			if item.QPS > 0 && item.QPM > 0 {
				return nil, state, errno.InternalServiceError.WithMessage("qps and qpm cannot both be set: " + item.Resource)
			}
			if item.QPS <= 0 && item.QPM <= 0 {
				return nil, state, errno.InternalServiceError.WithMessage("qps or qpm must be set: " + item.Resource)
			}

			if item.QPS > 0 {
				rules = append(rules, qpsRule(item.Resource, item.QPS))
				continue
			}
			rules = append(rules, qpmRule(item.Resource, item.QPM))
		}
	}

	state.RuleCount = len(rules)
	return rules, state, nil
}

// InitSentinelRateLimit 初始化 Sentinel 并加载启动时的限流规则。
func InitSentinelRateLimit() (int, error) {
	if err := sentinel.InitDefault(); err != nil {
		return 0, errno.InternalServiceError.WithMessage("init sentinel failed").WithError(err)
	}

	count, err := ReloadSentinelRateLimit(context.Background(), rateLimitConfigFromRuntime(), "startup")
	if err != nil {
		return 0, err
	}

	return count, nil
}

// RespRateLimitError 返回统一的限流错误响应。
func RespRateLimitError(c *app.RequestContext) {
	pack.RespError(c, errno.InternalServiceError.WithMessage(rateLimitMessage))
}

// rateLimitConfigFromRuntime 从全局运行时配置复制限流配置。
func rateLimitConfigFromRuntime() *rateLimitConfig {
	cfg := config.RateLimit
	if cfg == nil {
		return nil
	}

	rules := make([]rateLimitRuleConfig, 0, len(cfg.Interface.Rules))
	for _, item := range cfg.Interface.Rules {
		rules = append(rules, rateLimitRuleConfig{
			Resource:    item.Resource,
			QPS:         item.QPS,
			QPM:         item.QPM,
			Description: item.Description,
		})
	}

	return &rateLimitConfig{
		Enabled: cfg.Enabled,
		Global: rateLimitGlobalConfig{
			Enabled:  cfg.Global.Enabled,
			Resource: cfg.Global.Resource,
			QPS:      cfg.Global.QPS,
		},
		Interface: rateLimitInterfaceConfig{
			Enabled: cfg.Interface.Enabled,
			Rules:   rules,
		},
	}
}

// 构建QPS限流规则
func qpsRule(resource string, threshold float64) *flow.Rule {
	return &flow.Rule{
		Resource:               resource,
		Threshold:              threshold,
		TokenCalculateStrategy: flow.Direct,
		ControlBehavior:        flow.Reject,
		StatIntervalInMs:       qpsStatIntervalInMs,
	}
}

// 构建QPM限流规则
func qpmRule(resource string, threshold float64) *flow.Rule {
	return &flow.Rule{
		Resource:               resource,
		Threshold:              threshold,
		TokenCalculateStrategy: flow.Direct,
		ControlBehavior:        flow.Reject,
		StatIntervalInMs:       qpmStatIntervalInMs,
	}
}
