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
	"os"
	"reflect"
	"time"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

const (
	rateLimitRemoteFallbackPoll = 10 * time.Second
	rateLimitRemoteProvider     = "etcd3"
)

// etcd状态下监听配置
func watchRateLimitRemoteConfig() {
	etcdAddr, err := getRateLimitEtcdAddr()
	if err != nil {
		logger.Errorf("rate limit reloader: init remote config failed: %v", err)
		return
	}

	for {
		remoteViper, err := newRateLimitRemoteViper(etcdAddr)
		if err != nil {
			logger.Errorf("rate limit reloader: init remote config failed: %v", err)
			time.Sleep(rateLimitWatchRetry)
			continue
		}

		watchRateLimitRemoteChanges(remoteViper)
		time.Sleep(rateLimitWatchRetry)
	}
}

// getRateLimitEtcdAddr 获取用于读取远程配置的 etcd 地址。
func getRateLimitEtcdAddr() (string, error) {
	etcdAddr := ""
	if config.Etcd != nil {
		etcdAddr = config.Etcd.Addr
	}
	if etcdAddr == "" {
		etcdAddr = os.Getenv("ETCD_ADDR")
	}
	if etcdAddr == "" {
		return "", errno.InternalServiceError.WithMessage("rate limit remote config etcd addr is empty")
	}

	return etcdAddr, nil
}

// newRateLimitRemoteViper 创建只用于读取限流远程配置的 Viper 实例。
func newRateLimitRemoteViper(etcdAddr string) (*viper.Viper, error) {
	remoteViper := viper.New()
	if err := remoteViper.AddRemoteProvider(rateLimitRemoteProvider, etcdAddr, rateLimitRemotePath); err != nil {
		return nil, err
	}
	remoteViper.SetConfigType(rateLimitRemoteConfigTyp)
	return remoteViper, nil
}

// watchRateLimitRemoteChanges 定时从远程配置读取限流规则并在变化时重载。
func watchRateLimitRemoteChanges(remoteViper *viper.Viper) {
	ticker := time.NewTicker(rateLimitRemoteFallbackPoll)
	defer ticker.Stop()

	last, err := reloadRateLimitRemoteViper(remoteViper)
	if err != nil {
		logger.Errorf("rate limit reloader: read remote config failed: %v", err)
		return
	}

	for {
		<-ticker.C
		next, err := readRateLimitRemoteViper(remoteViper)
		if err != nil {
			logger.Errorf("rate limit reloader: watch remote config failed: %v", err)
			continue
		}
		// 通过反射递归比较配置结构体是否一致。
		if reflect.DeepEqual(next, last) {
			continue
		}

		if _, err := ReloadSentinelRateLimit(context.Background(), &next, "etcd"); err != nil {
			logger.Errorf("rate limit reloader: watch remote config failed: %v", err)
			continue
		}
		logger.Infof("rate limit reloader: remote config changed")
		last = next
	}
}

// readRateLimitRemoteViper 从 Viper 远程配置读取并解析 rate-limit 节点。
func readRateLimitRemoteViper(remoteViper *viper.Viper) (rateLimitConfig, error) {
	if err := remoteViper.WatchRemoteConfig(); err != nil {
		if err = remoteViper.ReadRemoteConfig(); err != nil {
			return rateLimitConfig{}, errno.InternalServiceError.WithMessage("read rate limit remote config failed").WithError(err)
		}
	}

	cfg := new(rateLimitConfig)
	if err := remoteViper.UnmarshalKey("rate-limit", cfg); err != nil {
		return rateLimitConfig{}, errno.InternalServiceError.WithMessage("unmarshal rate limit remote config failed").WithError(err)
	}

	return *cfg, nil
}

// reloadRateLimitRemoteViper 从远程配置读取 rate-limit 节点并重载限流规则。
func reloadRateLimitRemoteViper(remoteViper *viper.Viper) (rateLimitConfig, error) {
	cfg, err := readRateLimitRemoteViper(remoteViper)
	if err != nil {
		return rateLimitConfig{}, err
	}
	if _, err = ReloadSentinelRateLimit(context.Background(), &cfg, "etcd"); err != nil {
		return rateLimitConfig{}, errno.InternalServiceError.WithMessage("reload rate limit remote config failed").WithError(err)
	}

	return cfg, nil
}
