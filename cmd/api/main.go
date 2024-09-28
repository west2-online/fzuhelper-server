// Code generated by hertz generator.

package main

import (
	"flag"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/rpc"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

var (
	path       *string
	listenAddr string // listen port
)

func init() {
	//日志初始化
	// config init
	path = flag.String("config", "./config", "config path")
	flag.Parse()
	config.Init(*path, constants.ApiServiceName)
	//rpc
	utils.LoggerInit()
	rpc.Init()

}
func main() {
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

	h := server.New(
		server.WithHostPorts(listenAddr),
		server.WithHandleMethodNotAllowed(true),
		server.WithMaxRequestBodySize(1<<31),
	)

	register(h)
	h.Spin()
}
