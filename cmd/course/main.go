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
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/internal/course"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course/courseservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

var (
	serviceName = constants.CourseServiceName
	clientSet   *base.ClientSet
	taskQueue   taskqueue.TaskQueue
)

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())
	// eshook.InitLoggerWithHook(serviceName)
	clientSet = base.NewClientSet(base.WithDBClient(), base.WithRedisClient(constants.RedisDBCourse), base.WithCommonRPCClient(), base.WithUserRPCClient())
	taskQueue = taskqueue.NewBaseTaskQueue()
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		// 如果无法解析 etcd 的地址，则无法连接到其他的微服务，说明整个服务无法运行，直接 panic
		// 因为 API 只做数据包装返回和转发请求
		logger.Fatalf("Course: etcd registry failed, error: %v", err)
	}
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("Course: get available port failed: %v", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("Course: resolve tcp addr failed, err: %v", err)
	}

	svr := courseservice.NewServer(
		course.NewCourseService(clientSet, taskQueue),
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
	server.RegisterShutdownHook(clientSet.Close)
	taskQueue.AddSchedule(constants.LocateDateTaskKey, taskqueue.ScheduleQueueTask{
		Execute: func() error {
			locateDate, err := jwch.NewStudent().GetLocateDate()
			if err = base.HandleJwchError(err); err != nil {
				return err
			}
			currentDate := time.Now().In(constants.ChinaTZ)
			formattedCurrentDate := currentDate.Format(time.DateTime)
			result := &model.LocateDate{
				Year: locateDate.Year,
				Week: locateDate.Week,
				Term: locateDate.Term,
				Date: formattedCurrentDate,
			}

			return cache.SetStructCache(clientSet.CacheClient, context.Background(),
				constants.LocateDateKey, result, constants.LocateDateExpire, "Common.SetLocateDate")
		},
		GetScheduleTime: func() time.Duration {
			return constants.LocateDateUpdateTime
		},
	})
	taskQueue.Start()
	if err = svr.Run(); err != nil {
		logger.Fatalf("Course: run server failed, err: %v", err)
	}
}
