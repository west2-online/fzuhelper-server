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
	"net"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/internal/classroom"
	"github.com/west2-online/fzuhelper-server/internal/classroom/syncer"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

var (
	serviceName     = constants.ClassroomServiceName
	clientSet       *base.ClientSet
	classroomSyncer *syncer.EmptyRoomSyncer
)

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())
	// eshook.InitLoggerWithHook(serviceName)
	clientSet = base.NewClientSet(base.WithRedisClient(constants.RedisDBEmptyRoom))
	classroomSyncer = syncer.InitEmptyRoomSyncer(clientSet.CacheClient)
}

func main() {
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

	// update用于启动定期更新当天数据的任务
	classroomSyncer.Add("update")
	// 将scheduled放入队列，开启定时任务
	classroomSyncer.Add("schedule")
	classroomSyncer.Start()

	if err = svr.Run(); err != nil {
		logger.Fatalf("Classroom: server run failed: %v", err)
	}
}
