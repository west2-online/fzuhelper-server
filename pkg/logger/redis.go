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
	"context"
	"net"
	"time"

	kitexzap "github.com/kitex-contrib/obs-opentelemetry/logging/zap"
	"github.com/redis/go-redis/v9"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

type RedisLogger struct {
	*kitexzap.Logger
}

func (l *RedisLogger) Printf(ctx context.Context, template string, args ...interface{}) {
	l.Infof(template, args...)
}

func (l *RedisLogger) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (l *RedisLogger) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now().UnixMilli()

		if err := next(ctx, cmd); err != nil {
			return err
		}

		consume := time.Now().UnixMilli() - start
		if consume >= constants.RedisSlowQuery {
			Warnf("slowly redis query. consume %d microsecond, query: %s", consume, cmd.String())
		}

		return nil
	}
}

func (l *RedisLogger) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}

func GetRedisLogger() *RedisLogger {
	return &RedisLogger{loggerObj}
}
