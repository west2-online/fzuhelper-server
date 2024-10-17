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
	"fmt"
	"net"

	"github.com/west2-online/fzuhelper-server/cmd/template/dal/mq"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/west2-online/fzuhelper-server/pkg/eszap"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/elastic/go-elasticsearch"
	kitexzap "github.com/kitex-contrib/obs-opentelemetry/logging/zap"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"

	"github.com/west2-online/fzuhelper-server/cmd/template/dal"
	"github.com/west2-online/fzuhelper-server/cmd/template/rpc"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/template/templateservice"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/tracer"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

var (
	serviceName = constants.TemplateServiceName
	path        *string
	listenAddr  string // listen port

	EsClient *elasticsearch.Client
)

func Init() {
	// config init
	path = flag.String("config", "./config", "config path")
	flag.Parse()
	config.Init(*path, serviceName)

	// dal
	dal.Init()

	// trace
	tracer.InitJaeger(serviceName)

	// rpc
	rpc.Init()

	// log
	EsInit()
	klog.SetLevel(klog.LevelDebug)
	klog.SetLogger(kitexzap.NewLogger(EsHookLog()...))

	// mq
	mq.Init()
}

func main() {
	Init() // 做一些中间件的初始化

	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	// get available port from config set
	for index, addr := range config.Service.AddrList {
		if ok := utils.AddrCheck(addr); ok {
			listenAddr = addr
			break
		}

		if index == len(config.Service.AddrList)-1 {
			klog.Fatal("not available port from config")
		}
	}

	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		panic(err)
	}

	svr := templateservice.NewServer(
		new(TemplateServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: serviceName,
		}),
		server.WithMuxTransport(),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithSuite(trace.NewDefaultServerSuite()),
		server.WithLimit(&limit.Option{
			MaxConnections: constants.MaxConnections,
			MaxQPS:         constants.MaxQPS,
		}),
	)

	if err = svr.Run(); err != nil {
		panic(err)
	}
}

func EsHookLog() []kitexzap.Option {
	hook := eszap.NewElasticHook(EsClient, config.Elasticsearch.Host, serviceName, zapcore.DebugLevel)
	var options []kitexzap.Option
	options = append(options, kitexzap.WithCoreEnc(hook.Enc()))
	options = append(options, kitexzap.WithCoreLevel(hook.Lvl()))
	options = append(options, kitexzap.WithCoreWs(hook.Ws()))
	options = append(options, kitexzap.WithZapOptions(zap.Hooks(hook.Fire)))
	return options
}

// InitEs 初始化es
func EsInit() {
	esConn := fmt.Sprintf("http://%s", config.Elasticsearch.Addr)
	cfg := elasticsearch.Config{
		Addresses: []string{esConn},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	EsClient = client
}
