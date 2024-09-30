package client

import (
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/pkg/errors"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

//定义一系列的RPC客户端

func InitUserRPC() (*userservice.Client, error) {
	if config.Etcd.Addr == "" {
		logger.LoggerObj.Fatalf("config.Etcd.Addr is nil")
		return nil, errors.New("config.Etcd.Addr is nil")
	}
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		logger.LoggerObj.Errorf("etcd.NewEtcdResolver failed, err is %v", err)
		return nil, errors.Wrap(err, "pkg.rpc.rpc.InitUserRPC etcd.NewEtcdResolver failed")
	}
	client, err := userservice.NewClient("user", client.WithResolver(r), client.WithMuxConnection(constants.MuxConnection))
	if err != nil {
		logger.LoggerObj.Errorf("userservice.NewClient failed, err is %v", err)
		return nil, errors.Wrap(err, "pkg.rpc.rpc.InitUserRPC NewClient failed")
	}
	logger.LoggerObj.Info("InitUserRPC success")
	return &client, nil
}

func InitClassroomRPC() (*classroomservice.Client, error) {
	if config.Etcd.Addr == "" {
		logger.LoggerObj.Fatalf("config.Etcd.Addr is nil")
		return nil, errors.New("config.Etcd.Addr is nil")
	}
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		logger.LoggerObj.Errorf("etcd.NewEtcdResolver failed, err is %v", err)
		return nil, errors.Wrap(err, "pkg.rpc.rpc.InitClassroomRPC etcd.NewEtcdResolver failed")
	}
	client, err := classroomservice.NewClient("classroom", client.WithResolver(r), client.WithMuxConnection(constants.MuxConnection))
	if err != nil {
		logger.LoggerObj.Errorf("classroomservice.NewClient failed, err is %v", err)
		return nil, errors.Wrap(err, "pkg.rpc.rpc.InitClassroomRPC NewClient failed")
	}
	logger.LoggerObj.Info("InitClassroomRPC success")

	return &client, nil
}
