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
	"flag"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/remote/codec/thrift"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/netpoll"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/internal/academic"
	"github.com/west2-online/fzuhelper-server/internal/academic/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/academic/academicservice"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

var (
	serviceName = constants.AcademicServiceName
	clientSet   *base.ClientSet
	taskQueue   taskqueue.TaskQueue
	runTask     = flag.String("run-task", "", "manually run a specific task and exit")
)

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())
	clientSet = base.NewClientSet(base.WithDBClient(), base.WithRedisClient(constants.RedisDBAcademic))
	taskQueue = taskqueue.NewBaseTaskQueue()
}

func main() {
	flag.Parse()

	if *runTask != "" {
		if err := runManualTask(*runTask); err != nil {
			logger.Fatalf("Academic: manual task %s failed: %v", *runTask, err)
		}
		logger.Infof("Academic: manual task %s completed successfully", *runTask)
		return
	}

	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("Academic: etcd registry failed, error: %v", err)
	}
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("Academic: get available port failed: %v", err)
	}
	addr, err := netpoll.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("Academic: listen addr failed %v", err)
	}

	code := thrift.NewThriftCodecWithConfig(thrift.FrugalRead | thrift.FrugalWrite)
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
		server.WithPayloadCodec(code),
	)
	server.RegisterShutdownHook(clientSet.Close)

	taskQueue.AddSchedule(constants.CourseTeacherScoresTaskKey, taskqueue.ScheduleQueueTask{
		Execute: updateCourseTeacherScoresTask,
		GetScheduleTime: func() time.Duration {
			// 每天凌晨4点
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())
			if !next.After(now) {
				next = next.Add(constants.CourseTeacherScoresInterval)
			}
			return next.Sub(now)
		},
	})

	taskQueue.Start()
	if err = svr.Run(); err != nil {
		logger.Fatalf("Academic: server run failed: %v", err)
	}
}

func runManualTask(taskName string) error {
	switch taskName {
	case constants.CourseTeacherScoresTaskKey:
		return updateCourseTeacherScoresTask()
	default:
		return fmt.Errorf("unknown task: %s", taskName)
	}
}

func updateCourseTeacherScoresTask() error {
	logger.Infof("Academic: update course teacher scores task start")
	ctx := context.Background()
	svc := service.NewAcademicService(ctx, clientSet, nil)
	if err := svc.UpdateCourseTeacherScores(); err != nil {
		logger.Errorf("Academic: update course teacher scores task failed: %v", err)
		return err
	}
	logger.Infof("Academic: update course teacher scores task finished")
	return nil
}
