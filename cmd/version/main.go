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

package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/internal/version"
	"github.com/west2-online/fzuhelper-server/kitex_gen/version/versionservice"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

var (
	serviceName = constants.VersionServiceName
	clientSet   *base.ClientSet
	taskQueue   taskqueue.TaskQueue
)

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())
	// eshook.InitLoggerWithHook(serviceName)
	clientSet = base.NewClientSet(
		base.WithDBClient(),
		base.WithRedisClient(constants.RedisDBVersion),
	)
	taskQueue = taskqueue.NewBaseTaskQueue()
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("Version: etcd registry failed, error: %v", err)
	}
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("Version: get available port failed: %v", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("Version: listen addr failed %v", err)
	}

	svr := versionservice.NewServer(
		version.NewVersionService(clientSet),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: serviceName,
		}),
		server.WithMuxTransport(),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithLimit(&limit.Option{
			MaxConnections: constants.MaxConnections,
			MaxQPS:         constants.MaxQPS,
		}),
	)
	taskQueue.AddSchedule(constants.VersionVisitedTaskKey, taskqueue.ScheduleQueueTask{
		Execute: syncVersionVisitDailyTask,
		GetScheduleTime: func() time.Duration {
			now := time.Now().In(constants.ChinaTZ)
			nextRun := time.Date(now.Year(), now.Month(), now.Day(), constants.VersionVisitRefreshHour, constants.VersionVisitRefreshMinute, 0, 0, time.Local)
			if !now.Before(nextRun) {
				nextRun = nextRun.Add(constants.ONE_DAY)
			}
			// 动态处理到下一跳的刷新时间
			return nextRun.Sub(now)
		},
	})
	taskQueue.Start()

	if err = svr.Run(); err != nil {
		logger.Fatalf("Version: server run failed: %v", err)
	}
}

func syncVersionVisitDailyTask() error {
	ctx := context.Background()
	now := time.Now().Add(-1 * constants.ONE_DAY).In(constants.ChinaTZ)
	key := now.Format("2006-01-02")
	if !clientSet.CacheClient.IsKeyExist(ctx, key) {
		err := clientSet.CacheClient.Version.CreateVisitKey(ctx, key)
		if err != nil {
			return fmt.Errorf("version.TaskQueue: set visits key error: %w", err)
		}
		return nil
	}
	visits, err := clientSet.CacheClient.Version.GetVisit(ctx, key)
	if err != nil {
		return fmt.Errorf("version.TaskQueue: get visits error: %w", err)
	}
	ok, _, err := clientSet.DBClient.Version.GetVersion(ctx, key)
	if err != nil {
		return fmt.Errorf("version.TaskQueue: get version error: %w", err)
	}
	if !ok {
		err = clientSet.DBClient.Version.CreateVersion(ctx, &model.Visit{
			Date:   key,
			Visits: visits,
		})
		if err != nil {
			return fmt.Errorf("version.TaskQueue: create version error: %w", err)
		}
		return nil
	} else {
		err = clientSet.DBClient.Version.UpdateVersion(ctx, &model.Visit{
			Date:   key,
			Visits: visits,
		})
		if err != nil {
			return fmt.Errorf("version.TaskQueue: create version error: %w", err)
		}
		return nil
	}
}
