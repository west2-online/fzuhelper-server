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
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

const (
	rateLimitConfigPath      = "config/config.yaml"
	rateLimitRemotePath      = "/config"
	rateLimitRemoteConfigTyp = "yaml"
	rateLimitReloadDebounce  = 500 * time.Millisecond
	rateLimitWatchRetry      = 5 * time.Second
)

var rateLimitReloaderRun atomic.Bool

// 启动限流热更新配置循环
func StartRateLimitReloader() {
	if !rateLimitReloaderRun.CompareAndSwap(false, true) {
		return
	}

	// 判断是不是k8s环境
	if os.Getenv("DEPLOY_ENV") == "k8s" {
		logger.Infof("rate limit reloader start, mode=k8s")
		go watchRateLimitConfigFile(rateLimitConfigPath)
		return
	}

	// 判断是不是etcd环境
	logger.Infof("rate limit reloader start, mode=etcd")
	go watchRateLimitRemoteConfig()
}

// k8s状态下监听配置文件
func watchRateLimitConfigFile(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Errorf("rate limit reloader: new file watcher failed: %v", err)
		return
	}
	defer watcher.Close()

	// 用来遍历整个config目录
	dir := filepath.Dir(path)
	if err = watcher.Add(dir); err != nil {
		logger.Errorf("rate limit reloader: watch config file failed: %v", err)
		return
	}

	// 初始化一个定时器用来监听事件
	var timer *time.Timer
	var timerC <-chan time.Time
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()

	for {
		select {
		/*
			第一个事件就是确定这个目录文件有没有发生变化。
			如果有变化就进行一个判断，是我们想要的配置文件就发送更新信号。
			在这个期间如果又有新更新的话，就将定时器重置等待500ms再进行热更新。
		*/
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if !isRateLimitConfigEvent(path, event) {
				continue
			}
			if timer == nil {
				timer = time.NewTimer(rateLimitReloadDebounce)
				timerC = timer.C
				continue
			}
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(rateLimitReloadDebounce)
		case <-timerC:
			// 第二个事件就是加载配置文件
			timerC = nil
			timer = nil
			if _, err = reloadRateLimitFromFile(context.Background(), path); err != nil {
				logger.Errorf("rate limit reloader: reload config file failed: %v", err)
			}
		case err, ok := <-watcher.Errors:
			// 第三个事件就是发生了报错之后停止循环
			if !ok {
				return
			}
			logger.Errorf("rate limit reloader: watch config file error: %v", err)
		}
	}
}

// k8s模式下面重新加载配置文件
func reloadRateLimitFromFile(ctx context.Context, path string) (int, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return 0, errno.InternalServiceError.WithMessage("read rate limit config file failed").WithError(err)
	}

	cfg, err := parseRateLimitConfig(raw)
	if err != nil {
		return 0, errno.InternalServiceError.WithMessage("parse rate limit config file failed").WithError(err)
	}

	return ReloadSentinelRateLimit(ctx, cfg, "file")
}

// 解析限流配置文件，这里只对rate-limit节点进行解析
func parseRateLimitConfig(raw []byte) (*rateLimitConfig, error) {
	if len(raw) == 0 {
		return nil, errors.New("rate limit config is empty")
	}

	parser := viper.New()
	parser.SetConfigType(rateLimitRemoteConfigTyp)
	if err := parser.ReadConfig(bytes.NewReader(raw)); err != nil {
		return nil, err
	}

	cfg := new(rateLimitConfig)
	if err := parser.UnmarshalKey("rate-limit", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// 确定一下是不是我们想要的那个限流节点的配置文件发生了改变
func isRateLimitConfigEvent(path string, event fsnotify.Event) bool {
	if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) == 0 {
		return false
	}

	name := filepath.Base(event.Name)
	if name == filepath.Base(path) {
		return true
	}

	return strings.HasPrefix(name, "..data")
}
