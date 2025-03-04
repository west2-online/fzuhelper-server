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

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/internal/common"
	taskmodel "github.com/west2-online/fzuhelper-server/internal/common/task_model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/common/commonservice"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

var (
	serviceName = constants.CommonServiceName
	clientSet   *base.ClientSet
	taskQueue   taskqueue.TaskQueue
)

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())
	// eshook.InitLoggerWithHook(serviceName)
	clientSet = base.NewClientSet(base.WithDBClient(), base.WithRedisClient(constants.RedisDBCommon))
	taskQueue = taskqueue.NewBaseTaskQueue()
	loadNotice(clientSet.DBClient)
}

// TODO: 失败后的重试机制
func loadNotice(db *db.Database) {
	stu := jwch.NewStudent().WithUser(config.DefaultUser.Account, config.DefaultUser.Password)
	_, totalPage, err := stu.GetNoticeInfo(&jwch.NoticeInfoReq{PageNum: 1})
	if err != nil {
		logger.Errorf("syncer init: failed to get notice info: %v", err)
	}
	// 初始化数据库
	for i := 1; i <= totalPage; i++ {
		content, _, err := stu.GetNoticeInfo(&jwch.NoticeInfoReq{PageNum: i})
		if err != nil {
			logger.Errorf("syncer init: failed to get notice info in page %d: %v", i, err)
		}
		for _, row := range content {
			ctx := context.Background()
			info := &model.Notice{
				Title:       row.Title,
				PublishedAt: row.Date,
				URL:         row.URL,
			}
			err = db.Notice.CreateNotice(ctx, info)
			if err != nil {
				logger.Errorf("syncer init: failed to create notice in page %d: %v", i, err)
			}
		}
	}
	logger.Infof("syncer init: notice syncer init success")
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("Common: etcd registry failed, error: %v", err)
	}
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("Common: get available port failed: %v", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("Common: listen addr failed %v", err)
	}

	svr := commonservice.NewServer(
		common.NewCommonService(clientSet),
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

	taskQueue.Add(taskmodel.NewNoticeSyncTask(clientSet.DBClient))
	taskQueue.Add(taskmodel.NewContributorInfoSyncTask(clientSet.CacheClient))
	taskQueue.Start()

	if err = svr.Run(); err != nil {
		logger.Fatalf("Common: server run failed: %v", err)
	}
}
