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

	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/internal/classroom"
	"github.com/west2-online/fzuhelper-server/internal/classroom/refresh"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	baseserver "github.com/west2-online/fzuhelper-server/pkg/base/server"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/tracing"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

var (
	serviceName = constants.ClassroomServiceName
	clientSet   *base.ClientSet
	taskQueue   taskqueue.TaskQueue
)

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())
	// eshook.InitLoggerWithHook(serviceName)
	clientSet = base.NewClientSet(base.WithRedisClient(constants.RedisDBEmptyRoom))
	taskQueue = taskqueue.NewBaseTaskQueue()
}

func main() {
	// Open Telemetry provider
	shutdown := tracing.NewOtelProvider(serviceName, config.Otel.Endpoint)

	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("Classroom: etcd registry failed, error: %v", err)
	}
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("Classroom: get available port failed: %v", err)
	}
	logger.Infof("Classroom: listen addr: %v", listenAddr)
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("Classroom: listen addr failed %v", err)
	}

	svr := classroomservice.NewServer(
		classroom.NewClassroomService(clientSet),
		baseserver.AssembleCommonServerConfig(serviceName, addr, r)...,
	)
	server.RegisterShutdownHook(clientSet.Close)
	server.RegisterShutdownHook(tracing.ProviderShutdown(shutdown,
		"Classroom: otel provider shutdown failed: %v")) // otel provider

	taskQueue.AddSchedule("refresh", taskqueue.ScheduleQueueTask{
		Execute: func(ctx context.Context) error {
			stu := jwch.NewStudent().WithUser(config.DefaultUser.Account, config.DefaultUser.Password)
			if err := utils.RetryLogin(stu); err != nil {
				return fmt.Errorf("ScheduledUpdateClassroomsInfo: failed to login: %w", err)
			}
			runner := refresh.Refresher{
				Source:       stu,
				Cache:        clientSet.CacheClient.Classroom,
				RequestDelay: time.Second,
				RetryDelay:   time.Second,
			}
			return runner.RunScheduled(ctx, time.Now(), constants.CampusArray)
		},
		GetScheduleTime: func() time.Duration {
			return refresh.NextDelay(time.Now())
		},
	})

	taskQueue.Start()

	if err = svr.Run(); err != nil {
		logger.Fatalf("Classroom: server run failed: %v", err)
	}
}
