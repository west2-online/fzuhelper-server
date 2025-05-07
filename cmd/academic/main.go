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
	"github.com/west2-online/fzuhelper-server/internal/academic"
	"github.com/west2-online/fzuhelper-server/kitex_gen/academic/academicservice"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

const (
	dailyTriggerHour = 4
	dayHours         = 24
)

var (
	serviceName = constants.AcademicServiceName
	clientSet   *base.ClientSet
	taskQueue   taskqueue.TaskQueue
)

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())
	clientSet = base.NewClientSet(base.WithDBClient(), base.WithRedisClient(constants.RedisDBAcademic))
	taskQueue = taskqueue.NewBaseTaskQueue()
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("Academic: etcd registry failed, error: %v", err)
	}
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("Academic: get available port failed: %v", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("Academic: listen addr failed %v", err)
	}

	svr := academicservice.NewServer(
		academic.NewAcademicService(clientSet, taskQueue),
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

	taskQueue.AddSchedule(constants.CourseTeacherScoresTaskKey, taskqueue.ScheduleQueueTask{
		Execute:         UpdateCourseTeacherScoresTask,
		GetScheduleTime: durationUntilNext4AM,
	})
	taskQueue.Start()
	if err = svr.Run(); err != nil {
		logger.Fatalf("Academic: server run failed: %v", err)
	}
}

func UpdateCourseTeacherScoresTask() error {
	logger.Infof("Academic: update course teacher scores task start")
	now := time.Now()
	if now.Hour() != dailyTriggerHour {
		logger.Infof("current time is not 4 a.m. skip the execution")
		return nil
	}
	ctx := context.Background()
	// 统计总条数
	total, err := clientSet.DBClient.Academic.GetScoresCount(ctx)
	if err != nil {
		return fmt.Errorf("update course teacher scores task: academic GetScoresCount failed: %w", err)
	}
	pages := int((total + constants.SQLBatchSize - 1) / constants.SQLBatchSize)
	for i := 0; i < pages; i++ {
		offset := i * constants.SQLBatchSize
		key := fmt.Sprintf("%s-%d", constants.CourseTeacherScoresTaskKey, i)
		// 捕获局部 offset
		taskQueue.Add(key, taskqueue.QueueTask{
			Execute: func() error {
				return clientSet.DBClient.Academic.UpdateCourseTeacherScores(ctx, offset, constants.SQLBatchSize)
			},
		})
	}
	return nil
}

// durationUntilNext4AM 计算距离下一个凌晨4点的时间
func durationUntilNext4AM() time.Duration {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())
	if !next.After(now) {
		next = next.Add(dayHours * time.Hour)
	}
	return next.Sub(now)
}
