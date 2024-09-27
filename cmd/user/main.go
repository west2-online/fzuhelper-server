package main

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/spf13/viper"
	"github.com/west2-online/fzuhelper-server/cmd/user/config"
	user "github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
	"github.com/west2-online/fzuhelper-server/pkg/utils"

	"net"
)

func Init() {
	// config init
	//初始化配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&config.Config); err != nil {
		panic(err)
	}
	klog.SetLevel(klog.LevelDebug)
	// init logger
	utils.LoggerInit()
}

func main() {
	Init()
	conf := config.Config
	r, err := etcd.NewEtcdRegistry([]string{conf.EtcdHost + ":" + conf.EtcdPort})
	if err != nil {
		klog.Fatal(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", conf.System.Host+":"+conf.System.Port)
	if err != nil {
		panic(err)
	}
	svr := user.NewServer(new(UserServiceImpl), server.WithServiceAddr(addr), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "user"}), server.WithRegistry(r))
	err = svr.Run()
	if err != nil {
		panic(err)
	}
}
