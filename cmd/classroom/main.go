package main

import (
	"flag"
	"net"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/dal"
	"github.com/west2-online/fzuhelper-server/config"
	classroom "github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

var (
	serviceName = constants.ClassroomServiceName
	path        *string
)

func Init() {
	// config init
	logger.LoggerInit()

	path = flag.String("config", "./config", "config path")
	flag.Parse()
	config.Init(*path, serviceName)

	dal.Init()
	InitWorkerQueue()
}

func main() {
	Init()
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		// 如果无法解析etcd的地址，则无法连接到其他的微服务，说明整个服务无法运行,直接panic
		// 因为api只做数据包装返回和转发请求
		logger.LoggerObj.Fatalf("Classroom: etcd registry failed, error: %v", err)
	}
	// get available port from config set
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.LoggerObj.Fatalf("Classroom: get available port failed: %v", err)
	}

	addr, err := net.ResolveTCPAddr("tcp", listenAddr)

	if err != nil {
		logger.LoggerObj.Fatalf("Classroom: listen addr failed %v", err)
	}

	svr := classroom.NewServer(
		new(ClassroomServiceImpl),
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
	// 提前缓存空教室数据
	// 将signal放入队列，开启定时任务
	WorkQueue.Add("signal")

	if err = svr.Run(); err != nil {
		logger.LoggerObj.Fatalf("Classroom: server run failed: %v", err)
	}
}
