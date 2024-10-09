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
	"flag"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/netpoll"
	elastic "github.com/elastic/go-elasticsearch"
	etcd "github.com/kitex-contrib/registry-etcd"
	kopentracing "github.com/kitex-contrib/tracer-opentracing"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal"
	"github.com/west2-online/fzuhelper-server/config"
	launch_screen "github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen/launchscreenservice"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/tracer"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

var (
	path     *string
	EsClient *elastic.Client
)

func Init() {
	path = flag.String("config", "./config", "config path")
	flag.Parse()
	config.Init(*path, constants.LaunchScreenServiceName)
	dal.Init()
	tracer.InitJaeger(constants.LaunchScreenServiceName)
}

func main() {
	Init()
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("launchScreen: etcd registry failed, error: %v", err)
	}

	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("launchScreen: get available port failed: %v", err)
	}

	launchScreenServiceImpl := new(LaunchScreenServiceImpl)
	launchScreenCli, _ := NewLaunchScreenClient(listenAddr)

	serviceAddr, err := netpoll.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("launchScreen: listen addr failed %v", err)
	}

	launchScreenServiceImpl.launchScreenCli = launchScreenCli

	svr := launch_screen.NewServer(launchScreenServiceImpl, // 指定 Registry 与服务基本信息
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: constants.LaunchScreenServiceName,
			}),
		server.WithSuite(kopentracing.NewDefaultServerSuite()), //jaeger
		server.WithRegistry(r),
		server.WithServiceAddr(serviceAddr),
		server.WithLimit(
			&limit.Option{
				MaxConnections: constants.MaxConnections,
				MaxQPS:         constants.MaxQPS,
			},
		),
	)

	err = svr.Run()

	if err != nil {
		logger.Fatalf("launchScreen: server run failed: %v", err)
	}
}
