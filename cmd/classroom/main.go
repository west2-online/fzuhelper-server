package main

import (
	"flag"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/dal"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/service"
	"github.com/west2-online/fzuhelper-server/config"
	classroom "github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"net"
)

var (
	serviceName = constants.ClassroomService
	path        *string
	listenAddr  string // listen port

)

func Init() {
	// config init
	utils.LoggerInit()
	path = flag.String("config", "./config", "config path")
	flag.Parse()
	config.Init(*path, serviceName)

	dal.Init()
	klog.SetLevel(klog.LevelDebug)
}

func main() {
	Init()
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		klog.Fatal(err)
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
	//提前缓存空教室数据
	go service.CacheEmptyRooms()
	if err = svr.Run(); err != nil {
		panic(err)
	}
}
