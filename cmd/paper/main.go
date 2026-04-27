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
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/netpoll"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/internal/paper"
	"github.com/west2-online/fzuhelper-server/kitex_gen/paper/paperservice"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	baseserver "github.com/west2-online/fzuhelper-server/pkg/base/server"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/tracing"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

var (
	serviceName = constants.PaperServiceName
	clientSet   *base.ClientSet
)

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())
	// eshook.InitLoggerWithHook(serviceName)
	clientSet = base.NewClientSet(base.WithRedisClient(constants.RedisDBPaper))
	upyun.NewUpYun()
}

func main() {
	// Open Telemetry provider
	shutdown := tracing.NewOtelProvider(serviceName, config.Otel.Endpoint, config.Uptrace.DSN)

	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("Paper: etcd registry failed, error: %v", err)
	}
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("Paper: get available port failed: %v", err)
	}
	addr, err := netpoll.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("Paper: listen addr failed %v", err)
	}

	svr := paperservice.NewServer(
		paper.NewPaperService(clientSet),
		baseserver.AssembleCommonServerConfig(serviceName, addr, r)...,
	)
	server.RegisterShutdownHook(clientSet.Close)
	server.RegisterShutdownHook(tracing.ProviderShutdown(shutdown,
		"Paper: otel provider shutdown failed: %v")) // otel provider

	if err = svr.Run(); err != nil {
		logger.Fatalf("Paper: server run failed: %v", err)
	}
}
