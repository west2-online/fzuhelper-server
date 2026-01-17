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
	"os"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/netpoll"
	etcd "github.com/kitex-contrib/registry-etcd"
	"golang.org/x/sync/errgroup"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/internal/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
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
	var watcherCancel context.CancelFunc
	if os.Getenv("DEPLOY_ENV") != "k8s" {
		watcherCtx, cancel := context.WithCancel(context.Background())
		watcherCancel = cancel
		go config.StartEtcdWatcher(watcherCtx, serviceName)
	}

	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("Classroom: etcd registry failed, error: %v", err)
	}
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("Classroom: get available port failed: %v", err)
	}
	logger.Infof("Classroom: listen addr: %v", listenAddr)
	addr, err := netpoll.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("Classroom: listen addr failed %v", err)
	}

	svr := classroomservice.NewServer(
		classroom.NewClassroomService(clientSet),
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
	server.RegisterShutdownHook(func() {
		if watcherCancel != nil {
			logger.Info("Shutting down etcd config watcher...")
			watcherCancel()
		}
		logger.Info("Closing client resources...")
		clientSet.Close()
	})

	taskQueue.AddSchedule("update", taskqueue.ScheduleQueueTask{
		Execute: func() error {
			return updateEmptyClassroomsInfo(time.Now())
		},
		GetScheduleTime: func() time.Duration {
			return constants.ClassroomUpdatedTime
		},
	})
	taskQueue.AddSchedule("schedule", taskqueue.ScheduleQueueTask{
		Execute: func() error {
			return scheduleUpdateEmptyClassroomsInfo(time.Now())
		},
		GetScheduleTime: func() time.Duration {
			return constants.ClassroomScheduledTime
		},
	})

	taskQueue.Start()

	if err = svr.Run(); err != nil {
		logger.Fatalf("Classroom: server run failed: %v", err)
	}
}

func scheduleUpdateEmptyClassroomsInfo(date time.Time) error {
	var dates []time.Time
	// 设定一周时间
	for i := 1; i < 7; i++ {
		d := date.AddDate(0, 0, i)
		dates = append(dates, d)
	}

	var eg errgroup.Group
	// 对每个日期启动一个 goroutine
	for _, d := range dates {
		// 将 date 作为参数传递给 goroutine
		currentDate := d
		eg.Go(func() error {
			return updateEmptyClassroomsInfo(currentDate)
		})
	}

	// 等待所有 goroutine 完成，并收集错误
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("ScheduledUpdateClassroomsInfo: failed to refresh empty room info: %w", err)
	}
	logger.Infof("scheduleUpdateEmptyClassroomsInfo: complete all tasks of dates %v", dates)
	return nil
}

func updateEmptyClassroomsInfo(date time.Time) error {
	currentDate := date.Format("2006-01-02")
	ctx := context.Background()
	// 定义 jwch 的 stu 客户端
	stu := jwch.NewStudent().WithUser(config.DefaultUser.Account, config.DefaultUser.Password)
	// 登录，id 和 cookies 会自动保存在 client 中
	// 如果登录失败，重试
	err := utils.RetryLogin(stu)
	if err != nil {
		return fmt.Errorf("updateEmptyClassroomsInfo: failed to login: %w", err)
	}
	for _, campus := range constants.CampusArray {
		for startTime := 1; startTime <= 11; startTime++ {
			for endTime := startTime; endTime <= 11; endTime++ {
				args := jwch.EmptyRoomReq{
					Campus: campus,
					Time:   currentDate, // 使用传递进来的 date 参数
					Start:  strconv.Itoa(startTime),
					End:    strconv.Itoa(endTime),
				}
				var res []string
				var err error
				// 从 jwch 获取空教室信息
				switch campus {
				case "旗山校区":
					res, err = stu.GetQiShanEmptyRoom(args)
				default:
					res, err = stu.GetEmptyRoom(args)
				}
				if err != nil {
					return fmt.Errorf("updateEmptyClassroomsInfo: failed to get empty room info: %w", err)
				}
				// 收集结果并缓存
				switch campus {
				case "厦门工艺美院":
					err = clientSet.CacheClient.Classroom.SetXiaMenEmptyRoomCache(ctx, currentDate, args.Start, args.End, res)
					if err != nil {
						return fmt.Errorf("updateEmptyClassroomsInfo: failed to set xiamen empty room cache: %w", err)
					}
				default:
					key := fmt.Sprintf("%s.%s.%s.%s", args.Time, args.Campus, args.Start, args.End)
					err = clientSet.CacheClient.Classroom.SetEmptyRoomCache(ctx, key, res)
					if err != nil {
						return fmt.Errorf("updateEmptyClassroomsInfo: failed to set empty room cache: %w", err)
					}
				}
			}
		}
	}
	return nil
}
