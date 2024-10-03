package main

import (
	"flag"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/netpoll"
	"github.com/elastic/go-elasticsearch"
	elastic "github.com/elastic/go-elasticsearch"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	etcd "github.com/kitex-contrib/registry-etcd"
	kopentracing "github.com/kitex-contrib/tracer-opentracing"
	"github.com/sirupsen/logrus"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/pack/tracer"
	"github.com/west2-online/fzuhelper-server/cmd/user/dal"
	"github.com/west2-online/fzuhelper-server/config"
	user "github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/eslogrus"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"log"
)

var (
	path       *string
	listenAddr string
	EsClient   *elastic.Client
)

func Init() {
	path = flag.String("config", "./config", "config path")
	flag.Parse()
	logger.LoggerInit()
	config.Init(*path, constants.UserServiceName)
	dal.Init()
	tracer.InitJaegerTracer(constants.UserServiceName)
	EsInit()
	klog.SetLevel(klog.LevelWarn)
	klog.SetLogger(kitexlogrus.NewLogger(kitexlogrus.WithHook(EsHookLog())))
}

func main() {
	Init()
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	//获取addr
	for index, addr := range config.Service.AddrList {
		if ok := utils.AddrCheck(addr); ok {
			listenAddr = addr
			break
		}

		if index == len(config.Service.AddrList)-1 {
			klog.Fatal("not available addr")
		}
	}

	userHandlerImpl := new(UserServiceImpl)
	userCli, err := NewUserClient(listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	serviceAddr, err := netpoll.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	userHandlerImpl.userCli = userCli
	//然而不使用WithServiceAddr方法的话，server还是在监听8888
	//那Impl携带一个Client就没用了

	svr := user.NewServer(userHandlerImpl, // 指定 Registry 与服务基本信息
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: constants.UserServiceName,
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
		klog.Error(err.Error())
	}
}

func EsHookLog() *eslogrus.ElasticHook {
	hook, err := eslogrus.NewElasticHook(EsClient, config.Elasticsearch.Host, logrus.DebugLevel, constants.UserServiceName)
	if err != nil {
		panic(err)
	}

	return hook
}

// InitEs 初始化es
func EsInit() {
	esConn := fmt.Sprintf("http://%s", config.Elasticsearch.Addr)
	cfg := elasticsearch.Config{
		Addresses: []string{esConn},
	}
	klog.Infof("esConn:%v", esConn)
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	EsClient = client
}
