package main

import (
	"flag"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/netpoll"
	elastic "github.com/elastic/go-elasticsearch"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	etcd "github.com/kitex-contrib/registry-etcd"
	kopentracing "github.com/kitex-contrib/tracer-opentracing"
	"github.com/sirupsen/logrus"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal"
	"github.com/west2-online/fzuhelper-server/config"
	launch_screen "github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen/launchscreenservice"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/eslogrus"
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
	EsInit()
	klog.SetLevel(klog.LevelWarn)
	klog.SetLogger(kitexlogrus.NewLogger(kitexlogrus.WithHook(EsHookLog())))
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

func EsHookLog() *eslogrus.ElasticHook {
	hook, err := eslogrus.NewElasticHook(EsClient, config.Elasticsearch.Host, logrus.DebugLevel, constants.LaunchScreenServiceName)
	if err != nil {
		logger.Errorf("launchScreen: connect to es failed: %v", err)
	}

	return hook
}

// InitEs 初始化es
func EsInit() {
	esConn := fmt.Sprintf("http://%s", config.Elasticsearch.Addr)
	cfg := elastic.Config{
		Addresses: []string{esConn},
	}
	klog.Infof("esConn:%v", esConn)
	client, err := elastic.NewClient(cfg)
	if err != nil {
		logger.Errorf("launchScreen: connect to es failed: %v", err)
	}
	EsClient = client
}
