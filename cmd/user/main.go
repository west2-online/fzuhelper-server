package main

import (
	"flag"

	"net"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/west2-online/fzuhelper-server/config"
	user "github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

var (
	serviceName = constants.UserServiceName
	path        *string
)

func Init() {
	// config init
	path = flag.String("config", "./config", "config path")
	flag.Parse()
	config.Init(*path, serviceName)
}

func main() {
	Init()
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("User: new etcd registry failed, err: %v", err)
	}
	// get available port from config set
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("User: get available port failed, err: %v", err)
	}

	addr, err := net.ResolveTCPAddr("tcp", listenAddr)

	if err != nil {
		logger.Fatalf("User: resolve tcp addr failed, err: %v", err)
	}

	svr := user.NewServer(
		new(UserServiceImpl),
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
	if err = svr.Run(); err != nil {
		logger.Fatalf("User: run server failed, err: %v", err)
	}
}
